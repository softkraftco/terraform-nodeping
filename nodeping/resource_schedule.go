package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"terraform-nodeping/nodeping_api_client"
)

func resourceSchedule() *schema.Resource {

	weekdayNames := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}

	return &schema.Resource{
		CreateContext: resourceScheduleCreate,
		ReadContext:   resourceScheduleRead,
		UpdateContext: resourceScheduleUpdate,
		DeleteContext: resourceScheduleDelete,
		Schema: map[string]*schema.Schema{
			"name":        &schema.Schema{Type: schema.TypeString, Required: true},
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true, Optional: true},
			"data": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day":      &schema.Schema{Type: schema.TypeString, Optional: true, ValidateFunc: validation.StringInSlice(weekdayNames, false)},
						"time1":    &schema.Schema{Type: schema.TypeString, Optional: true},
						"time2":    &schema.Schema{Type: schema.TypeString, Optional: true},
						"exclude":  &schema.Schema{Type: schema.TypeBool, Optional: true},
						"disabled": &schema.Schema{Type: schema.TypeBool, Optional: true},
						"allday":   &schema.Schema{Type: schema.TypeBool, Optional: true},
					},
				},
			},
		},
	}
}

func applyScheduleToSchema(schedule *nodeping_api_client.Schedule, d *schema.ResourceData) error {
	d.SetId(schedule.Name)
	err := d.Set("name", schedule.Name)
	if err != nil {
		return err
	}

	err = d.Set("customer_id", schedule.CustomerId)
	if err != nil {
		return err
	}

	daysSchemasList := make([]map[string]interface{}, len(schedule.Data))
	counter := 0
	for day, dayConfig := range schedule.Data {
		daySchema := make(map[string]interface{})
		daySchema["day"] = day
		for configName, configValue := range dayConfig {
			// exclude can be a bool, but also an 0 or 1
			if configName == "exclude" {
				if configValue == float64(0) {
					configValue = false
				} else if configValue == float64(1) {
					configValue = true
				}
			}
			daySchema[configName] = configValue
		}
		daysSchemasList[counter] = daySchema
		counter++
	}

	err = d.Set("data", daysSchemasList)
	if err != nil {
		return err
	}

	return nil
}

func getScheduleFromSchema(d *schema.ResourceData) *nodeping_api_client.Schedule {
	var schedule nodeping_api_client.Schedule

	schedule.Name = d.Get("name").(string)
	schedule.CustomerId = d.Get("customer_id").(string)

	daySchemaList := d.Get("data").(*schema.Set).List()
	schedule.Data = make(map[string]map[string]interface{}, len(daySchemaList))
	for _, ds := range daySchemaList {
		daySchema := ds.(map[string]interface{})
		day := daySchema["day"].(string)
		schedule.Data[day] = make(map[string]interface{}, len(daySchema))
		schedule.Data[day]["time1"] = daySchema["time1"].(string)
		schedule.Data[day]["time2"] = daySchema["time2"].(string)
		schedule.Data[day]["exclude"] = daySchema["exclude"].(bool)
		schedule.Data[day]["disabled"] = daySchema["disabled"].(bool)
		schedule.Data[day]["allday"] = daySchema["allday"].(bool)
	}

	return &schedule
}

func resourceScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	schedule := getScheduleFromSchema(d)

	_, err := client.CreateSchedule(ctx, schedule)
	if err != nil {
		return diag.FromErr(err)
	}

	// this might seem silly, but API returns some kind of Id, that isn't used
	// anywhere. Instead schedule name (in API docs under the name "id") is used
	// for schedule identification.
	d.SetId(schedule.Name)
	d.Set("customer_id", schedule.CustomerId)
	return resourceScheduleRead(ctx, d, m)
}

func resourceScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	schedule, err := client.GetSchedule(ctx, d.Get("customer_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = applyScheduleToSchema(schedule, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceScheduleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	schema := getScheduleFromSchema(d)

	_, err := client.UpdateSchedule(ctx, schema)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceScheduleRead(ctx, d, m)
}

func resourceScheduleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	scheduleId := d.Id()
	customerId := d.Get("customer_id").(string)

	err := client.DeleteSchedule(ctx, customerId, scheduleId)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("") // added here for explicitness

	return nil

}
