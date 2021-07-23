package nodeping_test

import (
	"context"
	"encoding/json"
	"log"
	"testing"
	"time"

	apiClient "terraform-nodeping/nodeping_api_client"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestGroupDataSource(t *testing.T) {
	const terraformDir = "testdata/group_integration/data_source"
	const terraformMainFile = terraformDir + "/main.tf"

	// prepare API client
	client := getClient()

	// create a contact to read
	address := apiClient.Address{
		Address: "contact-test@example.com",
		Type:    "email",
	}

	address2 := apiClient.Address{
		Address: "contact-test2@example.com",
		Type:    "email",
	}

	contact := apiClient.Contact{
		Name:         "DataSourceTesting",
		Custrole:     "owner",
		Addresses:    make(map[string]apiClient.Address, 0),
		NewAddresses: []apiClient.Address{address, address2},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	contactPtr, err := client.CreateContact(ctx, &contact)
	if err != nil {
		log.Fatal(err)
	}
	contact = *contactPtr

	members := []string{}

	for key := range contact.Addresses {
		members = append(members, key)
	}

	group := apiClient.Group{
		Name:    "DataSourceTesting",
		Members: members,
	}

	groupPtr, err := client.CreateGroup(ctx, &group)
	if err != nil {
		log.Fatal(err)
	}

	group = *groupPtr

	// prepare contact cleanup
	defer client.DeleteContact(ctx, contact.ID)

	// prepare group cleanup
	defer client.DeleteGroup(ctx, group.ID)

	// create main.tf
	copyFile(terraformDir+"/data_source", terraformMainFile)

	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
		MaxRetries:   0,
		Upgrade:      true,
		Vars:         map[string]interface{}{"group_id": group.ID},
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// -----------------------------------
	// read a group
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, group.Name, terraform.Output(t, terraformOptions, "group_name"))

	// get group from output
	output := terraform.OutputJson(t, terraformOptions, "group")

	var outputObj map[string]interface{}
	json.Unmarshal([]byte(output), &outputObj)

	// check content
	assert.Equal(t, group.Name, outputObj["name"].(string))
	assert.Equal(t, group.ID, outputObj["id"].(string))
	assert.Equal(t, group.CustomerId, outputObj["customer_id"].(string))

	// unwrap "members" and check
	member1 := outputObj["members"].([]interface{})[0].(string)
	member2 := outputObj["members"].([]interface{})[1].(string)
	assert.Equal(t, group.Members[0], member1)
	assert.Equal(t, group.Members[1], member2)
}
