package nodeping

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	apiClient "terraform-nodeping/nodeping_api_client"
)

func TestTerraformContactLifeCycle(t *testing.T) {
	const terraformDir = "testdata/contacts_integration"
	const terraformMainFile = terraformDir + "/main.tf"
	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// prepare API client
	token := os.Getenv("NODEPING_API_TOKEN")
	client := apiClient.NewClient(token)

	// -----------------------------------
	// create a single contact
	copyFile(terraformDir+"/step_1", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	firstContractId := terraform.Output(t, terraformOptions, "first_contact_id")
	firstContact, err := client.GetContact(firstContractId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstContractId, firstContact.ID)
	assert.Equal(t, "First", firstContact.Name)
	assert.Equal(t, 1, len(firstContact.Addresses))
	// there is only one address fow, so this will work
	for _, addr := range firstContact.Addresses {
		assert.Equal(t, "first@o1.com", addr.Address)
		assert.Equal(t, "email", addr.Type)
		assert.Equal(t, false, addr.Suppressall)
	}

	// -----------------------------------
	// change contact name
	copyFile(terraformDir+"/step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	assert.Equal(t, firstContractId, terraform.Output(t, terraformOptions, "first_contact_id"))
	firstContact, err = client.GetContact(firstContractId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstContractId, firstContact.ID)
	assert.Equal(t, "First altered", firstContact.Name)
	assert.Equal(t, 1, len(firstContact.Addresses))

	// -----------------------------------
	// alter address attribute
	copyFile(terraformDir+"/step_3", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	assert.Equal(t, firstContractId, terraform.Output(t, terraformOptions, "first_contact_id"))
	firstContact, err = client.GetContact(firstContractId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstContractId, firstContact.ID)
	assert.Equal(t, "First", firstContact.Name)
	assert.Equal(t, 1, len(firstContact.Addresses))
	for _, addr := range firstContact.Addresses {
		assert.Equal(t, "first-altered@o1.com", addr.Address)
		assert.Equal(t, "email", addr.Type)
		assert.Equal(t, true, addr.Suppressall)
	}
	// -----------------------------------
	// add address to contact
	copyFile(terraformDir+"/step_4", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	firstContact, err = client.GetContact(firstContractId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstContractId, firstContact.ID)
	assert.Equal(t, 3, len(firstContact.Addresses))
	var email, webook, pushover bool
	for _, addr := range firstContact.Addresses {
		if addr.Type == "email" {
			assert.Equal(t, "first@o1.com", addr.Address)
			email = true
		}
		if addr.Type == "webhook" {
			assert.Equal(t, "first.com", addr.Address)
			assert.Equal(t, "PUT", addr.Action)
			assert.Equal(t, map[string]string{"the": "first"}, addr.Data)
			webook = true
		}
		if addr.Type == "pushover" {
			assert.Equal(t, "first.eu", addr.Address)
			assert.Equal(t, 2, addr.Priority)
			pushover = true
		}
	}
	assert.True(t, email, "'email' type address is missing from contact")
	assert.True(t, webook, "'webook' type address is missing from contact")
	assert.True(t, pushover, "'pushover' type address is missing from contact")
	// -----------------------------------
	// remove address from contact
	copyFile(terraformDir+"/step_5", terraformMainFile)

	terraform.Apply(t, terraformOptions)
	firstContact, err = client.GetContact(firstContractId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstContractId, firstContact.ID)
	assert.Equal(t, 2, len(firstContact.Addresses))
	email, webook, pushover = false, false, false
	for _, addr := range firstContact.Addresses {
		if addr.Type == "email" {
			email = true
		}
		if addr.Type == "webhook" {
			webook = true
		}
		if addr.Type == "pushover" {
			pushover = true
		}
	}
	assert.True(t, email, "'email' type address is missing from contact")
	assert.False(t, webook, "'webook' type address should have been removed from contact")
	assert.True(t, pushover, "'pushover' type address is missing from contact")
	// -----------------------------------
	// add new contact
	copyFile(terraformDir+"/step_6", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	firstContact, err = client.GetContact(firstContractId)
	if err != nil {
		log.Fatal(err)
	}

	secondContractId := terraform.Output(t, terraformOptions, "second_contact_id")
	secondContract, err := client.GetContact(secondContractId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, secondContractId, secondContract.ID)
	assert.Equal(t, "Second", secondContract.Name)
	assert.Equal(t, 1, len(secondContract.Addresses))
	// -----------------------------------
	// destroy the first contact
	copyFile(terraformDir+"/step_7", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	firstContact, err = client.GetContact(firstContractId)
	if assert.Error(t, err) {
		switch e := err.(type) {
		case *apiClient.ContactNotExists:
			// this is correct
		default:
			log.Fatal(e)
		}
	}

	assert.Equal(t, secondContractId, terraform.Output(t, terraformOptions, "second_contact_id"))
	secondContract, err = client.GetContact(secondContractId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, secondContractId, secondContract.ID)
	assert.Equal(t, "Second", secondContract.Name)
	assert.Equal(t, 1, len(secondContract.Addresses))
	// -----------------------------------
	// destroy
	terraform.Destroy(t, terraformOptions)
	secondContract, err = client.GetContact(secondContractId)
	if assert.Error(t, err) {
		switch e := err.(type) {
		case *apiClient.ContactNotExists:
			// this is correct
		default:
			log.Fatal(e)
		}
	}
}

func copyFile(src, dst string) {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(dst, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanupTerraformDir(terraformDir string) {
	err := os.RemoveAll(terraformDir + "/.terraform")
	if err != nil {
		log.Fatal(err)
	}
	for _, fileName := range []string{".terraform.lock.hcl", "main.tf",
		"terraform.tfstate", "terraform.tfstate.backup"} {
		err = os.Remove(terraformDir + "/" + fileName)
		if err != nil {
			log.Fatal(err)
		}
	}
}
