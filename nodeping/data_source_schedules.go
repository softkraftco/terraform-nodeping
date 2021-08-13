package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func dataSourceSchedule() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceScheduleRead,
		Schema: map[string]*schema.Schema{
			"name":        &schema.Schema{Type: schema.TypeString, Required: true},
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true, Optional: true},
			"data": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day":      &schema.Schema{Type: schema.TypeString, Computed: true},
						"time1":    &schema.Schema{Type: schema.TypeString, Computed: true},
						"time2":    &schema.Schema{Type: schema.TypeString, Computed: true},
						"exclude":  &schema.Schema{Type: schema.TypeBool, Computed: true},
						"disabled": &schema.Schema{Type: schema.TypeBool, Computed: true},
						"allday":   &schema.Schema{Type: schema.TypeBool, Computed: true},
					},
				},
			},
		},
	}
}

func dataSourceScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	scheduleId := d.Get("name").(string)
	customerId := d.Get("customer_id").(string)

	schedule, err := client.GetSchedule(ctx, customerId, scheduleId)
	if err != nil {
		return diag.FromErr(err)
	}
	err = applyScheduleToSchema(schedule, d)
	if err != nil {
		return diag.FromErr(err)
	}
	return nil
}
