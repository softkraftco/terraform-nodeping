package nodeping

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	apiClient "terraform-nodeping/nodeping_api_client"
)

func TestTerraformCheckLifeCycle(t *testing.T) {
	const terraformDir = "testdata/checks_integration"
	const terraformMainFile = terraformDir + "/main.tf"
	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
		MaxRetries:   1,
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// prepare API client
	token := os.Getenv("NODEPING_API_TOKEN")
	client := apiClient.NewClient(token)

	// -----------------------------------
	// create a single HTTP check
	copyFile(terraformDir+"/step_1", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	firstCheckId := terraform.Output(t, terraformOptions, "first_check_id")
	firstAddressId := terraform.Output(t, terraformOptions, "first_address_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err := client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "FirstCheck", firstCheck.Label)
	assert.Equal(t, "HTTP", firstCheck.Type)
	assert.Equal(t, "inactive", firstCheck.Enable)

	assert.Equal(t, 1, len(firstCheck.Notifications))
	assert.Equal(t, 1, len(firstCheck.Notifications[0]))
	assert.Equal(t, 1, firstCheck.Notifications[0][firstAddressId].Delay)
	assert.Equal(t, "Weekdays", firstCheck.Notifications[0][firstAddressId].Schedule)
	// -----------------------------------
	// change check "enabled" property
	copyFile(terraformDir+"/http_step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	firstAddressId = terraform.Output(t, terraformOptions, "first_address_id")

	assert.Equal(t, firstCheckId, terraform.Output(t, terraformOptions, "first_check_id"))

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err = client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "FirstCheck", firstCheck.Label)
	assert.Equal(t, "HTTP", firstCheck.Type)
	assert.Equal(t, "active", firstCheck.Enable)

	assert.Equal(t, 1, len(firstCheck.Notifications))
	assert.Equal(t, 1, len(firstCheck.Notifications[0]))
	assert.Equal(t, 1, firstCheck.Notifications[0][firstAddressId].Delay)
	assert.Equal(t, "Weekdays", firstCheck.Notifications[0][firstAddressId].Schedule)
	// -----------------------------------
	// add notification to check
	copyFile(terraformDir+"/http_step_3", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	firstAddressId = terraform.Output(t, terraformOptions, "first_address_id")
	secondAddressId := terraform.Output(t, terraformOptions, "second_address_id")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err = client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "inactive", firstCheck.Enable)

	assert.Equal(t, 2, len(firstCheck.Notifications))

	for _, notificationMap := range firstCheck.Notifications {
		assert.Equal(t, 1, len(notificationMap))
		if notification, exists := notificationMap[firstAddressId]; exists {
			assert.Equal(t, 1, notification.Delay)
			assert.Equal(t, "Weekdays", notification.Schedule)
		} else if notification, exists := notificationMap[secondAddressId]; exists {
			assert.Equal(t, 20, notification.Delay)
			assert.Equal(t, "Nights", notification.Schedule)
		} else {
			t.Fail()
		}
	}
	// -----------------------------------
	// remove notification from check
	copyFile(terraformDir+"/http_step_4", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	secondAddressId = terraform.Output(t, terraformOptions, "second_address_id")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err = client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 1, len(firstCheck.Notifications))
	assert.Equal(t, 1, len(firstCheck.Notifications[0]))
	assert.Equal(t, 20, firstCheck.Notifications[0][secondAddressId].Delay)
	assert.Equal(t, "Nights", firstCheck.Notifications[0][secondAddressId].Schedule)
	// -----------------------------------
	// change notification properties
	copyFile(terraformDir+"/http_step_5", terraformMainFile)
	terraform.Apply(t, terraformOptions)
	secondAddressId = terraform.Output(t, terraformOptions, "second_address_id")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err = client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, 1, len(firstCheck.Notifications))
	assert.Equal(t, 1, len(firstCheck.Notifications[0]))
	assert.Equal(t, 10, firstCheck.Notifications[0][secondAddressId].Delay)
	assert.Equal(t, "All", firstCheck.Notifications[0][secondAddressId].Schedule)
	// -----------------------------------
	// destroy
	terraform.Destroy(t, terraformOptions)
}

func copyFile(src, dst string) { // TODO: don't copy this
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
