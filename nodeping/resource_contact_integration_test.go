package nodeping_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	apiClient "terraform-nodeping/nodeping_api_client"
)

func TestTerraformContactLifeCycle(t *testing.T) {
	const terraformDir = "testdata/contacts_integration/resource"
	const terraformMainFile = terraformDir + "/main.tf"

	// create main.tf
	copyFile(terraformDir+"/step_1", terraformMainFile)

	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// prepare API client
	client := getClient()

	// -----------------------------------
	// create a single contact
	terraform.Apply(t, terraformOptions)
	firstContractId := terraform.Output(t, terraformOptions, "first_contact_id")
	firstContractCustomerId := terraform.Output(t, terraformOptions, "first_contact_customer_id")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstContact, err := client.GetContact(ctx, firstContractCustomerId, firstContractId)
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstContact, err = client.GetContact(ctx, firstContractCustomerId, firstContractId)
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstContact, err = client.GetContact(ctx, firstContractCustomerId, firstContractId)
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstContact, err = client.GetContact(ctx, firstContractCustomerId, firstContractId)
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstContact, err = client.GetContact(ctx, firstContractCustomerId, firstContractId)
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstContact, err = client.GetContact(ctx, firstContractCustomerId, firstContractId)
	if err != nil {
		log.Fatal(err)
	}

	secondContractId := terraform.Output(t, terraformOptions, "second_contact_id")
	secondContractCustomerId := terraform.Output(t, terraformOptions, "second_contact_customer_id")
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	secondContract, err := client.GetContact(ctx, secondContractCustomerId, secondContractId)
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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstContact, err = client.GetContact(ctx, secondContractCustomerId, firstContractId)
	if assert.Error(t, err) {
		switch e := err.(type) {
		case *apiClient.ContactDoesNotExist:
			// this is correct
		default:
			assert.Fail(t, fmt.Sprintf("Call to GetContact raised an unexpected error: %s", e))
		}
	}

	assert.Equal(t, secondContractId, terraform.Output(t, terraformOptions, "second_contact_id"))
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	secondContract, err = client.GetContact(ctx, secondContractCustomerId, secondContractId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, secondContractId, secondContract.ID)
	assert.Equal(t, "Second", secondContract.Name)
	assert.Equal(t, 1, len(secondContract.Addresses))
	// -----------------------------------
	// destroy
	terraform.Destroy(t, terraformOptions)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	secondContract, err = client.GetContact(ctx, secondContractCustomerId, secondContractId)
	if assert.Error(t, err) {
		switch e := err.(type) {
		case *apiClient.ContactDoesNotExist:
			// this is correct
		default:
			assert.Fail(t, fmt.Sprintf("Call to GetContact raised an unexpected error: %s", e))
		}
	}
}
