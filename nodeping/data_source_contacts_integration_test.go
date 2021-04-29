package nodeping

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	apiClient "terraform-nodeping/nodeping_api_client"
)

func TestContactDataSource(t *testing.T) {
	const terraformDir = "testdata/contacts_integration"
	const terraformMainFile = terraformDir + "/main.tf"

	// prepare API client
	token := os.Getenv("NODEPING_API_TOKEN")
	client := apiClient.NewClient(token)

	// create a contact to read
	address := apiClient.Address{
		Address: "contact-test@example.com",
		Type:    "email",
	}

	contact := apiClient.Contact{
		Name:         "DataSourceTesting",
		Custrole:     "owner",
		Addresses:    make(map[string]apiClient.Address, 0),
		NewAddresses: []apiClient.Address{address},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	contactPtr, err := client.CreateContact(ctx, &contact)
	if err != nil {
		log.Fatal(err)
	}
	contact = *contactPtr

	// prepare contact cleanup
	defer client.DeleteContact(ctx, contact.ID)

	// create main.tf
	copyFile(terraformDir+"/data_source", terraformMainFile)

	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
		MaxRetries:   0,
		Upgrade:      true,
		Vars:         map[string]interface{}{"contact_id": contact.ID},
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// -----------------------------------
	// read a contact
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, contact.Name, terraform.Output(t, terraformOptions, "contact_name"))

	// get contact from output
	output := terraform.OutputJson(t, terraformOptions, "contact")

	var outputObj map[string]interface{}
	json.Unmarshal([]byte(output), &outputObj)

	// check content
	assert.Equal(t, contact.Name, outputObj["name"].(string))
	assert.Equal(t, contact.ID, outputObj["id"].(string))
	assert.Equal(t, contact.Custrole, outputObj["custrole"].(string))
	assert.Equal(t, contact.Type, outputObj["type"].(string))

	// unwrap "addresses"
	addressData := outputObj["addresses"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, address.Address, addressData["address"].(string))
	assert.Equal(t, address.Type, addressData["type"].(string))
	assert.Greater(t, len(addressData["id"].(string)), 0)
	assert.Equal(t, false, addressData["suppressall"].(bool))
}
