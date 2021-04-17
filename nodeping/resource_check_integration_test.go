package nodeping

import (
	"context"
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

	// create main.tf
	copyFile(terraformDir+"/http_step_1", terraformMainFile)

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
	// -----------------------------------
	// change check properties
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
	assert.Equal(t, false, firstCheck.HomeLoc)
	assert.Equal(t, 30, firstCheck.Interval)
	assert.Equal(t, false, firstCheck.Public)
	assert.Equal(t, []string{"eur"}, firstCheck.Runlocations)
	assert.Equal(t, float64(4), firstCheck.Parameters["threshold"])
	assert.Equal(t, float64(5), firstCheck.Parameters["sens"])
	assert.Equal(t, "Testing 12345", firstCheck.Description)

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
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = client.GetCheck(ctx, firstCheckId)
	assert.Error(t, err)
}

func TestTerraformHTTPCheck(t *testing.T) {
	/*
		Checks if changes to HTTP specific attributes work properly.
	*/
	const terraformDir = "testdata/checks_integration"
	const terraformMainFile = terraformDir + "/main.tf"

	// create main.tf
	copyFile(terraformDir+"/http_step_1", terraformMainFile)

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
	terraform.Apply(t, terraformOptions)
	firstCheckId := terraform.Output(t, terraformOptions, "first_check_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err := client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "HTTP", firstCheck.Type)
	assert.Equal(t, true, firstCheck.Parameters["ipv6"])
	assert.Equal(t, true, firstCheck.Parameters["follow"])
	assert.Equal(t, "inactive", firstCheck.Enable)

	// -----------------------------------
	// change check ipv6 property
	copyFile(terraformDir+"/http_step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err = client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, false, firstCheck.Parameters["ipv6"])
	assert.Equal(t, false, firstCheck.Parameters["follow"])
	// -----------------------------------
	// destroy
	terraform.Destroy(t, terraformOptions)
}

func TestTerraformSSHCheck(t *testing.T) {
	/*
		Checks if changes to SSH specific attributes work properly.
	*/
	const terraformDir = "testdata/checks_integration"
	const terraformMainFile = terraformDir + "/main.tf"

	// create main.tf
	copyFile(terraformDir+"/ssh_step_1", terraformMainFile)

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
	// create a single SSH check
	terraform.Apply(t, terraformOptions)
	firstCheckId := terraform.Output(t, terraformOptions, "first_check_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err := client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "SSH", firstCheck.Type)
	assert.Equal(t, "contentstring", firstCheck.Parameters["contentstring"])
	assert.Equal(t, float64(1000), firstCheck.Parameters["port"])
	assert.Equal(t, "username", firstCheck.Parameters["username"])
	assert.Equal(t, "password", firstCheck.Parameters["password"])
	assert.Equal(t, "true", firstCheck.Parameters["verify"])
	assert.Equal(t, true, firstCheck.Parameters["invert"])
	assert.Equal(t, "inactive", firstCheck.Enable)

	// -----------------------------------
	// change check "enabled" property
	copyFile(terraformDir+"/ssh_step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err = client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "changed contentstring", firstCheck.Parameters["contentstring"])
	assert.Equal(t, float64(900), firstCheck.Parameters["port"])
	assert.Equal(t, "different username", firstCheck.Parameters["username"])
	assert.Equal(t, "another password", firstCheck.Parameters["password"])
	assert.Equal(t, "false", firstCheck.Parameters["verify"])
	assert.Equal(t, false, firstCheck.Parameters["invert"])
	// -----------------------------------
	// destroy
	terraform.Destroy(t, terraformOptions)
}

func TestTerraformSSLCheck(t *testing.T) {
	/*
		Checks if changes to SSL specific attributes work properly.
	*/
	const terraformDir = "testdata/checks_integration"
	const terraformMainFile = terraformDir + "/main.tf"

	// create main.tf
	copyFile(terraformDir+"/ssl_step_1", terraformMainFile)

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
	// create a single SSL check
	terraform.Apply(t, terraformOptions)
	firstCheckId := terraform.Output(t, terraformOptions, "first_check_id")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err := client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "SSL", firstCheck.Type)
	assert.Equal(t, float64(10), firstCheck.Parameters["warningdays"])
	assert.Equal(t, "http://example.eu/", firstCheck.Parameters["servername"])
	assert.Equal(t, "inactive", firstCheck.Enable)

	// -----------------------------------
	// change check "enabled" property
	copyFile(terraformDir+"/ssl_step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	firstCheck, err = client.GetCheck(ctx, firstCheckId)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, firstCheckId, firstCheck.ID)
	assert.Equal(t, "SSL", firstCheck.Type)
	assert.Equal(t, "http://example.com/", firstCheck.Parameters["target"])
	assert.Equal(t, float64(14), firstCheck.Parameters["warningdays"])
	assert.Equal(t, "http://example.com/", firstCheck.Parameters["servername"])
	assert.Equal(t, "inactive", firstCheck.Enable)
	// -----------------------------------
	// destroy
	terraform.Destroy(t, terraformOptions)
}
