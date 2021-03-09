package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func dataSourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	var diags diag.Diagnostics

	contact, err := client.GetContact(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(contact.ID)
	d.Set("type", contact.Type)
	d.Set("name", contact.Name)
	d.Set("customer_id", contact.CustomerId)
	d.Set("custrole", contact.Custrole)

	addresses := flattenAddresses(&contact.Addresses)
	if err := d.Set("addresses", &addresses); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func flattenAddresses(addresses *map[string]nodeping_api_client.Address) []interface{} {
	if addresses == nil { // return fast if nothing to do
		return make([]interface{}, 0)
	}

	flattenedAddresses := make([]interface{}, len(*addresses), len(*addresses))
	i := 0
	for addressId, address := range *addresses {
		flattenedAddress := make(map[string]interface{})
		flattenedAddress["id"] = addressId
		flattenedAddress["address"] = address.Address
		flattenedAddress["type"] = address.Type
		flattenedAddress["suppressup"] = address.Suppressup
		flattenedAddress["suppressdown"] = address.Suppressdown
		flattenedAddress["suppressfirst"] = address.Suppressfirst
		flattenedAddress["suppressall"] = address.Suppressall

		flattenedAddresses[i] = flattenedAddress
		i++
	}

	return flattenedAddresses
}

func dataSourceContact() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContactRead,
		Schema: map[string]*schema.Schema{
			"id":          &schema.Schema{Type: schema.TypeString, Required: true},
			"type":        &schema.Schema{Type: schema.TypeString, Computed: true},
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true},
			"name":        &schema.Schema{Type: schema.TypeString, Computed: true},
			"custrole":    &schema.Schema{Type: schema.TypeString, Computed: true},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":            &schema.Schema{Type: schema.TypeString, Computed: true},
						"address":       &schema.Schema{Type: schema.TypeString, Computed: true},
						"type":          &schema.Schema{Type: schema.TypeString, Computed: true},
						"suppressup":    &schema.Schema{Type: schema.TypeBool, Computed: true},
						"suppressdown":  &schema.Schema{Type: schema.TypeBool, Computed: true},
						"suppressfirst": &schema.Schema{Type: schema.TypeBool, Computed: true},
						"suppressall":   &schema.Schema{Type: schema.TypeBool, Computed: true},
					},
				},
			},
		},
	}
}
