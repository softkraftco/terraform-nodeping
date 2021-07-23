package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func dataSourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	group, err := client.GetGroup(ctx, d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.ID)
	d.Set("customer_id", group.CustomerId)
	d.Set("name", group.Name)
	d.Set("members", group.Members)

	return nil
}

func dataSourceGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGroupRead,
		Schema: map[string]*schema.Schema{
			"id":          &schema.Schema{Type: schema.TypeString, Required: true},
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true},
			"name":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"members":     &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
}
