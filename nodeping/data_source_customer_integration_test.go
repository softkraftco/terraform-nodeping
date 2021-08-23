package nodeping_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	apiClient "terraform-nodeping/nodeping_api_client"
)

func TestCustomerDataSource(t *testing.T) {
	const terraformDir = "testdata/customer_integration/data_source"
	const terraformMainFile = terraformDir + "/main.tf"

	// prepare API client
	client := getClient()

	customer := apiClient.Customer{
		Name:        "DataSourceTesting",
		ContactName: "aaaaa",
		Email:       "aa@bb.cc",
		Timezone:    "1",
		Location:    "nam",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	customerPtr, err := client.CreateCustomer(ctx, &customer)
	if err != nil {
		log.Fatal(err)
	}
	customer = *customerPtr

	// prepare contact cleanup
	defer client.DeleteCustomer(ctx, customer.ID)

	// create main.tf
	copyFile(terraformDir+"/data_source", terraformMainFile)

	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
		MaxRetries:   1,
		Upgrade:      true,
		Vars:         map[string]interface{}{"customer_id": customer.ID},
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// -----------------------------------
	// read a contact
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, customer.Name, terraform.Output(t, terraformOptions, "name"))

	// get customer from output
	output := terraform.OutputJson(t, terraformOptions, "customer")

	var outputObj map[string]interface{}
	json.Unmarshal([]byte(output), &outputObj)

	// check content
	assert.Equal(t, customer.Name, outputObj["name"].(string))
	// assert.Equal(t, customer.ContactName, outputObj["contactname"].(string)) // think what to do
	assert.Equal(t, customer.Status, "Active")
	assert.Equal(t, customer.Timezone, outputObj["timezone"].(string))
	assert.Equal(t, customer.Location, outputObj["location"].(string))
	assert.Equal(t, customer.Emailme, outputObj["emailme"].(bool))
}
