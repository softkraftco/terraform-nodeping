package nodeping_test

import (
	"encoding/json"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestScheduleDataSource(t *testing.T) {
	const terraformDir = "testdata/schedules_integration/data_source"
	const terraformMainFile = terraformDir + "/main.tf"

	// create main.tf
	copyFile(terraformDir+"/data_source", terraformMainFile)

	// initialize terraform
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: terraformDir,
		MaxRetries:   0,
		Upgrade:      true,
	})
	terraform.Init(t, terraformOptions)

	// prepare cleanup
	defer cleanupTerraformDir(terraformDir)

	// -----------------------------------
	// read a schedule - "Weekdays" is available by default
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, "Weekdays", terraform.Output(t, terraformOptions, "schedule_name"))

	// get schedule from output
	output := terraform.OutputJson(t, terraformOptions, "schedule")

	var outputObj map[string]interface{}
	json.Unmarshal([]byte(output), &outputObj)

	// check content
	assert.Equal(t, "Weekdays", outputObj["name"].(string))
	assert.Equal(t, "Weekdays", outputObj["id"].(string))
	assert.Equal(t, "", outputObj["customer_id"].(string))

	// unwrap "data"
	weekdays := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	for _, dayMapInterface := range outputObj["data"].([]interface{}) {
		dayMap := dayMapInterface.(map[string]interface{})
		day := dayMap["day"].(string)
		assert.Contains(t, weekdays, day)
		if day == "saturday" || day == "sunday" {
			assert.False(t, dayMap["allday"].(bool))
			assert.True(t, dayMap["disabled"].(bool))
			assert.False(t, dayMap["exclude"].(bool))
			assert.Equal(t, "12:00", dayMap["time1"].(string))
			assert.Equal(t, "12:00", dayMap["time2"].(string))
		} else {
			assert.True(t, dayMap["allday"].(bool))
			assert.False(t, dayMap["disabled"].(bool))
			assert.False(t, dayMap["exclude"].(bool))
			assert.Equal(t, "0:00", dayMap["time1"].(string))
			assert.Equal(t, "23:59", dayMap["time2"].(string))
		}
	}
}
