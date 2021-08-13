package nodeping_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	apiClient "terraform-nodeping/nodeping_api_client"
)

func TestTerraformScheduleLifeCycle(t *testing.T) {
	const terraformDir = "testdata/schedules_integration/resource"
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
	// create a single schedule
	terraform.Apply(t, terraformOptions)

	scheduleName := terraform.Output(t, terraformOptions, "first_schedule_name")
	scheduleCustomerId := terraform.Output(t, terraformOptions, "first_schedule_customer_id")

	schedule, err := client.GetSchedule(ctx, scheduleCustomerId, scheduleName)
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

	// -----------------------------------
	// change schedule data
	copyFile(terraformDir+"/step_2", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, scheduleName, terraform.Output(t, terraformOptions, "first_schedule_name"))

	schedule, err = client.GetSchedule(ctx, scheduleCustomerId, scheduleName)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, scheduleName, schedule.Name)
	assert.Equal(t, 7, len(schedule.Data))
	// check days
	assert.Equal(t, schedule.Data["monday"]["time1"], "16:00")
	assert.Equal(t, schedule.Data["monday"]["time2"], "17:00")
	assert.Equal(t, schedule.Data["monday"]["exclude"], false)
	assert.Equal(t, schedule.Data["monday"]["disabled"], false)
	assert.Equal(t, schedule.Data["monday"]["allday"], false)

	assert.Equal(t, schedule.Data["saturday"]["time1"], "")
	assert.Equal(t, schedule.Data["saturday"]["time2"], "")
	assert.Equal(t, schedule.Data["saturday"]["exclude"], false)
	assert.Equal(t, schedule.Data["saturday"]["disabled"], false)
	assert.Equal(t, schedule.Data["saturday"]["allday"], true)

	assert.Equal(t, schedule.Data["sunday"]["time1"], "")
	assert.Equal(t, schedule.Data["sunday"]["time2"], "")
	assert.Equal(t, schedule.Data["sunday"]["exclude"], false)
	assert.Equal(t, schedule.Data["sunday"]["disabled"], false)
	assert.Equal(t, schedule.Data["sunday"]["allday"], true)
	// -----------------------------------
	// change schedule data once more
	copyFile(terraformDir+"/step_3", terraformMainFile)
	terraform.Apply(t, terraformOptions)

	assert.Equal(t, scheduleName, terraform.Output(t, terraformOptions, "first_schedule_name"))

	schedule, err = client.GetSchedule(ctx, scheduleCustomerId, scheduleName)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, scheduleName, schedule.Name)
	assert.Equal(t, 1, len(schedule.Data))
	// check days
	assert.Equal(t, nil, schedule.Data["monday"]["time1"])
	assert.Equal(t, nil, schedule.Data["monday"]["time2"])
	assert.Equal(t, nil, schedule.Data["monday"]["exclude"])
	assert.Equal(t, nil, schedule.Data["monday"]["disabled"])
	assert.Equal(t, nil, schedule.Data["monday"]["allday"])

	assert.Equal(t, "6:00", schedule.Data["sunday"]["time1"])
	assert.Equal(t, "7:00", schedule.Data["sunday"]["time2"])
	assert.Equal(t, false, schedule.Data["sunday"]["exclude"])
	assert.Equal(t, false, schedule.Data["sunday"]["disabled"])
	assert.Equal(t, false, schedule.Data["sunday"]["allday"])
	// -----------------------------------
	// destroy
	terraform.Destroy(t, terraformOptions)
	schedule, err = client.GetSchedule(ctx, scheduleCustomerId, scheduleName)
	if assert.Error(t, err) {
		switch e := err.(type) {
		case *apiClient.ScheduleDoesNotExist:
			// this is correct
		default:
			assert.Fail(t, fmt.Sprintf("Call to GetSchedule raised an unexpected error. %s", e))
		}
	}
}
