package nodeping_test

import (
	"context"
	"fmt"
	"log"
	apiClient "terraform-nodeping/nodeping_api_client"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformCustomerLifeCycle(t *testing.T) {
	const terraformDir = "testdata/customer_integration/resource"
	const terraformMainFile = terraformDir + "/main.tf"

	// create main.tf
	copyFile(terraformDir+"/step_1", terraformMainFile)

	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
		MaxRetries:   0,
		Upgrade:      true,
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// prepare API client
	client := getClient()

	// prepare context for client
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// -----------------------------------
	// create a single group
	terraform.Apply(t, terraformOptions)

	customerID := terraform.Output(t, terraformOptions, "customer_id")

	customer, err := client.GetCustomer(ctx, customerID)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, customerID, customer.ID)
	assert.Equal(t, customer.Name, "new_subaccount")
	assert.Equal(t, customer.Status, "Active")
	assert.Equal(t, customer.Timezone, "1.0")
	assert.Equal(t, customer.Location, "nam")
	assert.Equal(t, customer.Emailme, false)

	// -----------------------------------
	// change group data
	copyFile(terraformDir+"/step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, customerID, terraform.Output(t, terraformOptions, "customer_id"))

	customer, err = client.GetCustomer(ctx, customerID)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, customerID, customer.ID)
	assert.Equal(t, customer.Name, "old_subaccount")
	assert.Equal(t, customer.Status, "Suspend")
	assert.Equal(t, customer.Timezone, "2.0")
	assert.Equal(t, customer.Location, "nam")
	assert.Equal(t, customer.Emailme, true)

	terraform.Destroy(t, terraformOptions)
	customer, err = client.GetCustomer(ctx, customerID)
	if err != nil {
		if assert.Error(t, err) {
			switch e := err.(type) {
			case *apiClient.CustomerDoesNotExist:
				// this is correct
			default:
				assert.Fail(t, fmt.Sprintf("Call to GetCustomer raised an unexpected error. %s", e))
			}
		}
	} else {
		if customer.Status != "Delete" {
			assert.Fail(t, fmt.Sprintf("Subaccount should be deleted."))
		}
	}
}
