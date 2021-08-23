package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"terraform-nodeping/nodeping_api_client"
)

func resourceContact() *schema.Resource {
	// prepare accepted values for validation
	addressTypes := []string{"email", "sms", "webhook", "slack", "hipchat", "pushover", "pagerduty", "voice"}
	custroles := []string{"owner", "edit", "view", "notify"}
	webhookActions := []string{"GET", "POST", "PUT", "DELETE", "HEAD", "PATCH"}

	return &schema.Resource{
		CreateContext: resourceContactCreate,
		ReadContext:   resourceContactRead,
		UpdateContext: resourceContactUpdate,
		DeleteContext: resourceContactDelete,
		Schema: map[string]*schema.Schema{
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true, Optional: true},
			"name":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"custrole":    &schema.Schema{Type: schema.TypeString, Optional: true, ValidateFunc: validation.StringInSlice(custroles, false)},
			"addresses": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":            &schema.Schema{Type: schema.TypeString, Computed: true},
						"address":       &schema.Schema{Type: schema.TypeString, Required: true},
						"type":          &schema.Schema{Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice(addressTypes, false)},
						"suppressup":    &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressdown":  &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressfirst": &schema.Schema{Type: schema.TypeBool, Optional: true},
						"suppressall":   &schema.Schema{Type: schema.TypeBool, Optional: true},
						// webhooks related attributes
						"action":       &schema.Schema{Type: schema.TypeString, Optional: true, ValidateFunc: validation.StringInSlice(webhookActions, false)},
						"data":         &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString},
						"headers":      &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString},
						"querystrings": &schema.Schema{Type: schema.TypeMap, Optional: true, Elem: schema.TypeString},
						// pushover attributes
						"priority": &schema.Schema{Type: schema.TypeInt, Optional: true, ValidateFunc: validation.IntBetween(-2, 2)},
					},
				},
			},
		},
	}
}

func getContactFromSchema(d *schema.ResourceData) *nodeping_api_client.Contact {
	var contact nodeping_api_client.Contact
	contact.ID = d.Id()
	contact.CustomerId = d.Get("customer_id").(string)
	contact.Name = d.Get("name").(string)
	contact.Custrole = d.Get("custrole").(string)

	addrs := d.Get("addresses").([]interface{})
	addresses := make(map[string]nodeping_api_client.Address)
	newAddresses := make([]nodeping_api_client.Address, 0)
	for _, addr := range addrs {
		a := addr.(map[string]interface{})

		// get address Id (if present)
		addressId := a["id"].(string)

		// convert "data", "headers" and "querystrings" from interface{}
		// to map[string]string
		data := make(map[string]string)
		for key, val := range a["data"].(map[string]interface{}) {
			data[key] = val.(string)
		}
		headers := make(map[string]string)
		for key, val := range a["headers"].(map[string]interface{}) {
			headers[key] = val.(string)
		}
		querystrings := make(map[string]string)
		for key, val := range a["querystrings"].(map[string]interface{}) {
			querystrings[key] = val.(string)
		}

		address := nodeping_api_client.Address{
			ID:            a["id"].(string),
			Address:       a["address"].(string),
			Type:          a["type"].(string),
			Suppressup:    a["suppressup"].(bool),
			Suppressdown:  a["suppressdown"].(bool),
			Suppressfirst: a["suppressfirst"].(bool),
			Suppressall:   a["suppressall"].(bool),
			Action:        a["action"].(string),
			Data:          data,
			Headers:       headers,
			Querystrings:  querystrings,
			Priority:      a["priority"].(int),
		}

		// for lack of better documentation, I assume that addresses without a
		// type are the ones that got deleted
		if len(address.Type) != 0 {
			// addresses that have an id go to addresses, the ones that don't,
			// go to new addresses
			if len(addressId) > 0 {
				addresses[addressId] = address
			} else {
				newAddresses = append(newAddresses, address)
			}
		}
	}
	contact.Addresses = addresses
	contact.NewAddresses = newAddresses

	return &contact
}

func resourceContactCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	contact := getContactFromSchema(d)

	savedContact, err := client.CreateContact(ctx, contact)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(savedContact.ID)
	d.Set("customer_id", savedContact.CustomerId)
	return resourceContactRead(ctx, d, m)
}

func resourceContactRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	contact, err := client.GetContact(ctx, d.Get("customer_id").(string), d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(contact.ID)
	d.Set("customer_id", contact.CustomerId)
	d.Set("name", contact.Name)
	d.Set("custrole", contact.Custrole)

	addresses := flattenAddresses(&contact.Addresses)

	// sort addresses to match previous ordering
	orderedAddresses := make([]interface{}, len(addresses))
	for idx, a := range d.Get("addresses").([]interface{}) {
		addrSchema := a.(map[string]interface{})
		addressId := addrSchema["id"].(string)

		for _, ad := range addresses {
			address := ad.(map[string]interface{})
			if addressId == address["id"].(string) ||
				address["type"].(string) == addrSchema["type"].(string) &&
					address["address"].(string) == addrSchema["address"].(string) {
				orderedAddresses[idx] = address
				break
			}
		}
	}

	if err := d.Set("addresses", orderedAddresses); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceContactUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	contact := getContactFromSchema(d)

	_, err := client.UpdateContact(ctx, contact)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceContactRead(ctx, d, m)
}

func resourceContactDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	contactId := d.Id()
	customerId := d.Get("customer_id").(string)

	err := client.DeleteContact(ctx, customerId, contactId)
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return nil
}
