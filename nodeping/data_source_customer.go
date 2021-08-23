package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func dataSourceCustomerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	customer, err := client.GetCustomer(ctx, d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(customer.ID)
	d.Set("name", customer.Name)
	d.Set("status", customer.Status)
	d.Set("parent", customer.Parent)
	d.Set("creation_date", customer.CreationDate)
	d.Set("emailme", customer.Emailme)
	d.Set("timezone", customer.Timezone)
	d.Set("location", customer.Location)

	return nil
}

func dataSourceCustomer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCustomerRead,
		Schema: map[string]*schema.Schema{
			"id":            &schema.Schema{Type: schema.TypeString, Required: true},
			"parent":        &schema.Schema{Type: schema.TypeString, Computed: true},
			"name":          &schema.Schema{Type: schema.TypeString, Computed: true},
			"timezone":      &schema.Schema{Type: schema.TypeString, Computed: true},
			"creation_date": &schema.Schema{Type: schema.TypeInt, Computed: true},
			"status":        &schema.Schema{Type: schema.TypeString, Computed: true},
			"emailme":       &schema.Schema{Type: schema.TypeBool, Optional: true},
			"location":      &schema.Schema{Type: schema.TypeString, Computed: true},
		},
	}
}
