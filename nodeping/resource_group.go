package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupCreate,
		ReadContext:   resourceGroupRead,
		UpdateContext: resourceGroupUpdate,
		DeleteContext: resourceGroupDelete,
		Schema: map[string]*schema.Schema{
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true},
			"name":        &schema.Schema{Type: schema.TypeString, Optional: true},
			"members":     &schema.Schema{Type: schema.TypeList, Optional: true, Elem: &schema.Schema{Type: schema.TypeString}},
		},
	}
}

func getGroupFromSchema(d *schema.ResourceData) *nodeping_api_client.Group {
	var group nodeping_api_client.Group
	group.ID = d.Id()
	group.CustomerId = d.Get("customer_id").(string)
	group.Name = d.Get("name").(string)
	membrs := d.Get("members").([]interface{})
	members := []string{}

	for _, membr := range membrs {
		members = append(members, membr.(string))
	}
	group.Members = members

	return &group
}

func resourceGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)
	group := getGroupFromSchema(d)

	savedGroup, err := client.CreateGroup(ctx, group)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(savedGroup.ID)
	return resourceGroupRead(ctx, d, m)
}

func resourceGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	group, err := client.GetGroup(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.ID)
	d.Set("customer_id", group.CustomerId)
	d.Set("name", group.Name)
	d.Set("members", group.Members)

	return nil
}

func resourceGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)
	group := getGroupFromSchema(d)

	_, err := client.UpdateGroup(ctx, group)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGroupRead(ctx, d, m)
}

func resourceGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	err := client.DeleteGroup(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return nil
}
