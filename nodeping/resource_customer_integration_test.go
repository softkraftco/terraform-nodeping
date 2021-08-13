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
	// create a single customer
	terraform.Apply(t, terraformOptions)

	customerId := terraform.Output(t, terraformOptions, "customer_id")

	customer, err := client.GetCustomer(ctx, customerId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, customerId, customer.ID)
	assert.Equal(t, customer.Name, "old_subaccount")
	assert.Equal(t, customer.Status, "Active")
	assert.Equal(t, customer.Timezone, "1.0")
	assert.Equal(t, customer.Location, "nam")
	assert.Equal(t, customer.Emailme, false)

	contacts, err := client.GetContacts(ctx, customer.ID)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(contacts), 1)
	assert.Equal(t, contacts[0].Name, "John")
	assert.Equal(t, contacts[0].Custrole, "edit")
	assert.Equal(t, len(contacts[0].Addresses), 1)
	assert.Equal(t, contacts[0].Addresses["john@doe.com"].Address, "john@doe.com")

	// change customer data
	copyFile(terraformDir+"/step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, customerId, terraform.Output(t, terraformOptions, "customer_id"))

	customer, err = client.GetCustomer(ctx, customerId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, customerId, customer.ID)
	assert.Equal(t, customer.Name, "new_subaccount")
	assert.Equal(t, customer.Status, "Suspend")
	assert.Equal(t, customer.Timezone, "2.0")
	assert.Equal(t, customer.Location, "nam")
	assert.Equal(t, customer.Emailme, true)

	contacts, err = client.GetContacts(ctx, customer.ID)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, len(contacts), 1)
	// assert.Equal(t, contacts[0].Name, "Mike") you cannot update name by api
	// assert.Equal(t, contacts[0].Addresses["john@doe1.com"].Address, "john@doe1.com") you cannot update email by api
	assert.Equal(t, contacts[0].Custrole, "edit")
	assert.Equal(t, len(contacts[0].Addresses), 1)

	copyFile(terraformDir+"/step_3", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	firstContractId := terraform.Output(t, terraformOptions, "first_contact_id")
	firstContact, err := client.GetContact(ctx, customerId, firstContractId)
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

	groupID := terraform.Output(t, terraformOptions, "group_id")

	group, err := client.GetGroup(ctx, customerId, groupID)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, groupID, group.ID)
	assert.Equal(t, "test", group.Name)
	assert.Equal(t, 1, len(group.Members))

	scheduleName := terraform.Output(t, terraformOptions, "first_schedule_name")

	schedule, err := client.GetSchedule(ctx, customerId, scheduleName)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, scheduleName, schedule.Name)
	assert.Equal(t, "FirstSchedule", schedule.Name)
	assert.Equal(t, 7, len(schedule.Data))
	// check days
	assert.Equal(t, schedule.Data["monday"]["time1"], "6:00")
	assert.Equal(t, schedule.Data["monday"]["time2"], "7:00")
	assert.Equal(t, schedule.Data["monday"]["exclude"], false)
	assert.Equal(t, schedule.Data["monday"]["disabled"], false)
	assert.Equal(t, schedule.Data["monday"]["allday"], false)

	assert.Equal(t, schedule.Data["saturday"]["time1"], "6:00")
	assert.Equal(t, schedule.Data["saturday"]["time2"], "7:00")
	assert.Equal(t, schedule.Data["saturday"]["exclude"], true)
	assert.Equal(t, schedule.Data["saturday"]["disabled"], false)
	assert.Equal(t, schedule.Data["saturday"]["allday"], false)

	assert.Equal(t, schedule.Data["sunday"]["time1"], "")
	assert.Equal(t, schedule.Data["sunday"]["time2"], "")
	assert.Equal(t, schedule.Data["sunday"]["exclude"], false)
	assert.Equal(t, schedule.Data["sunday"]["disabled"], false)
	assert.Equal(t, schedule.Data["sunday"]["allday"], true)

	firstCheckId := terraform.Output(t, terraformOptions, "first_check_id")
	firstAddressId := terraform.Output(t, terraformOptions, "first_address_id")
	firstCheck, err := client.GetCheck(ctx, customerId, firstCheckId)

	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "FirstCheck", firstCheck.Label)
	assert.Equal(t, "HTTP", firstCheck.Type)
	assert.Equal(t, "inactive", firstCheck.Enable)
	assert.Equal(t, false, firstCheck.HomeLoc.(bool)) // Homeloc is available only with Provider plan
	assert.Equal(t, 25, firstCheck.Interval)
	assert.Equal(t, true, firstCheck.Public)

	assert.Contains(t, firstCheck.Runlocations, "eur")
	assert.Contains(t, firstCheck.Runlocations, "nam")
	assert.Equal(t, 2, len(firstCheck.Runlocations))
	assert.Equal(t, float64(3), firstCheck.Parameters["threshold"])
	assert.Equal(t, float64(1), firstCheck.Parameters["sens"])
	assert.Equal(t, "Testing 123", firstCheck.Description)

	assert.Equal(t, 1, len(firstCheck.Notifications))
	assert.Equal(t, 1, len(firstCheck.Notifications[0]))
	assert.Equal(t, 1, firstCheck.Notifications[0][firstAddressId].Delay)
	assert.Equal(t, "Weekdays", firstCheck.Notifications[0][firstAddressId].Schedule)

	/*
		test customer ids for resources
	*/

	terraform.Destroy(t, terraformOptions)
	customer, err = client.GetCustomer(ctx, customerId)
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
