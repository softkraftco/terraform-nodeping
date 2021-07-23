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

func TestTerraformGroupLifeCycle(t *testing.T) {
	const terraformDir = "testdata/group_integration/resource"
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

	groupID := terraform.Output(t, terraformOptions, "group_id")

	group, err := client.GetGroup(ctx, groupID)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, groupID, group.ID)
	assert.Equal(t, "test", group.Name)
	assert.Equal(t, 1, len(group.Members))

	// // -----------------------------------
	// change group data
	copyFile(terraformDir+"/step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, groupID, terraform.Output(t, terraformOptions, "group_id"))

	group, err = client.GetGroup(ctx, groupID)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, groupID, group.ID)
	assert.Equal(t, "test2", group.Name)
	assert.Equal(t, 2, len(group.Members))

	// // -----------------------------------
	// change group data once more
	copyFile(terraformDir+"/step_3", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, groupID, terraform.Output(t, terraformOptions, "group_id"))

	group, err = client.GetGroup(ctx, groupID)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, groupID, group.ID)
	assert.Equal(t, 0, len(group.Members))
	// // -----------------------------------
	// // destroy
	terraform.Destroy(t, terraformOptions)
	group, err = client.GetGroup(ctx, groupID)
	if assert.Error(t, err) {
		switch e := err.(type) {
		case *apiClient.GroupDoesNotExist:
			// this is correct
		default:
			assert.Fail(t, fmt.Sprintf("Call to GetGroup raised an unexpected error. %s", e))
		}
	}
}
