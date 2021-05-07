package nodeping_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	apiClient "terraform-nodeping/nodeping_api_client"
)

func TestCheckDataSource(t *testing.T) {
	t.Parallel()

	// prepare API client
	client := getClient()

	// prepare a single HTTP check
	notifications := make([]map[string]apiClient.Notification, 1)
	nmap := make(map[string]apiClient.Notification, 1)
	nmap["exmp"] = apiClient.Notification{
		Delay:    5,
		Schedule: "Days",
	}
	notifications[0] = nmap

	checkUpdate := apiClient.CheckUpdate{
		Label:         "TheCheck",
		Type:          "HTTP",
		Target:        "http://example.com",
		Enable:        "inactive",
		Public:        "false",
		RunLocations:  []string{"nam"},
		Notifications: notifications,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	check, err := client.CreateCheck(ctx, &checkUpdate)
	if err != nil {
		log.Fatal(err)
	}
	// prepare cleanup
	defer client.DeleteCheck(ctx, check.ID)

	const terraformDir = "testdata/checks_integration/data_source"
	const terraformMainFile = terraformDir + "/main.tf"

	// create main.tf
	copyFile(terraformDir+"/data_source", terraformMainFile)

	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
		MaxRetries:   1,
		Vars:         map[string]interface{}{"check_id": check.ID},
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// -----------------------------------
	// use data source
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, check.ID, terraform.Output(t, terraformOptions, "check_id"))
	assert.Equal(t, check.Type, terraform.Output(t, terraformOptions, "check_type"))
	assert.Equal(t, check.Parameters["target"].(string), terraform.Output(t, terraformOptions, "check_target"))
	assert.Equal(t, check.Enable, terraform.Output(t, terraformOptions, "check_enable"))
	assert.Equal(t, "false", terraform.Output(t, terraformOptions, "check_public"))
	assert.Equal(t, "nam", terraform.Output(t, terraformOptions, "check_runlocations"))
}
