package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func resourceContact() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContactCreate,
		ReadContext:   resourceContactRead,
		UpdateContext: resourceContactUpdate,
		DeleteContext: resourceContactDelete,
		Schema: map[string]*schema.Schema{
			"customer_id": &schema.Schema{Type: schema.TypeString, Optional: true},
			"name":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"custrole":    &schema.Schema{Type: schema.TypeString, Optional: true},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":            &schema.Schema{Type: schema.TypeString, Computed: true},
						"address":       &schema.Schema{Type: schema.TypeString, Required: true},
						"type":          &schema.Schema{Type: schema.TypeString, Required: true},
						"suppressup":    &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressdown":  &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressfirst": &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressall":   &schema.Schema{Type: schema.TypeBool, Optional: true},
					},
				},
			},
			// webhooks related attributes
			"action":       &schema.Schema{Type: schema.TypeString, Optional: true},
			"data":         &schema.Schema{Type: schema.TypeString, Optional: true},
			"headers":      &schema.Schema{Type: schema.TypeString, Optional: true},
			"querystrings": &schema.Schema{Type: schema.TypeString, Optional: true},
		},
	}
}

func resourceContactCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// map resource into Contact
	var contact nodeping_api_client.NewContact
	contact.CustomerId = d.Get("customer_id").(string)
	contact.Name = d.Get("name").(string)
	contact.Custrole = d.Get("custrole").(string)

	addrs := d.Get("addresses").([]interface{})
	addresses := make([]nodeping_api_client.Address, len(addrs))
	for i, addr := range addrs {
		a := addr.(map[string]interface{})
		address := nodeping_api_client.Address{
			Address:       a["address"].(string),
			Type:          a["type"].(string),
			Suppressup:    a["suppressup"].(bool),
			Suppressdown:  a["suppressdown"].(bool),
			Suppressfirst: a["suppressfirst"].(bool),
			Suppressall:   a["suppressall"].(bool),
		}
		addresses[i] = address
	}
	contact.Addresses = addresses

	savedContact, err := client.CreateContact(&contact)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(savedContact.ID)
	resourceContactRead(ctx, d, m)

	return diags
}

func resourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	contact, err := client.GetContact(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	addresses := flattenAddresses(&contact.Addresses)
	if err := d.Set("addresses", addresses); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceContactUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}

func resourceContactDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	contactID := d.Id()

	err := client.DeleteContact(contactID)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}
