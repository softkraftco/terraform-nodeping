package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"terraform-nodeping/nodeping_api_client"
)

func resourceCustomer() *schema.Resource {
	locations := []string{"nam", "eur", "eao", "wlw"}
	return &schema.Resource{
		CreateContext: resourceCustomerCreate,
		ReadContext:   resourceCustomerRead,
		UpdateContext: resourceCustomerUpdate,
		DeleteContext: resourceCustomerDelete,
		Schema: map[string]*schema.Schema{
			"id":          &schema.Schema{Type: schema.TypeString, Computed: true},
			"name":        &schema.Schema{Type: schema.TypeString, Required: true},
			"contactname": &schema.Schema{Type: schema.TypeString, Required: true},
			"email":       &schema.Schema{Type: schema.TypeString, Required: true},
			"timezone":    &schema.Schema{Type: schema.TypeString, Required: true},
			"location":    &schema.Schema{Type: schema.TypeString, Required: true, ValidateFunc: validation.StringInSlice(locations, false)},
			"emailme":     &schema.Schema{Type: schema.TypeBool, Optional: true},
			"status":      &schema.Schema{Type: schema.TypeString, Optional: true},
		},
	}
}

func getCustomerFromSchema(d *schema.ResourceData) *nodeping_api_client.Customer {
	var customer nodeping_api_client.Customer
	customer.ID = d.Id()
	customer.Name = d.Get("name").(string)
	customer.ContactName = d.Get("contactname").(string)
	customer.Email = d.Get("email").(string)
	customer.Timezone = d.Get("timezone").(string)
	customer.Location = d.Get("location").(string)
	customer.Emailme = d.Get("emailme").(bool)
	customer.Status = d.Get("status").(string)

	//customer.Parent = d.Get("parent").(string)
	// customer.CustomerName = d.Get("customer_name").(string)
	// customer.CreationDate = d.Get("creation_date").(int)
	// customer.Status = d.Get("status").(string)

	return &customer
}

func resourceCustomerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)
	customer := getCustomerFromSchema(d)

	savedCustomer, err := client.CreateCustomer(ctx, customer)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(savedCustomer.ID)
	return resourceCustomerRead(ctx, d, m)
}

func resourceCustomerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	customer, err := client.GetCustomer(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(customer.ID)
	d.Set("timezone", customer.Timezone)
	d.Set("customer_name", customer.ContactName)

	return nil
}

func resourceCustomerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)
	customer := getCustomerFromSchema(d)

	_, err := client.UpdateCustomer(ctx, customer)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCustomerRead(ctx, d, m)
}

func resourceCustomerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	err := client.DeleteCustomer(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return nil
}
