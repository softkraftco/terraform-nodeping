package nodeping

import (
	"context"
	//"encoding/json"
	"fmt"
	//"net/http"
	//"strconv"
	//"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceContactsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	fmt.Println(d)

	var diags diag.Diagnostics

	return diags
}

func dataSourceContacts() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContactsRead,
		Schema: map[string]*schema.Schema{
			"contacts": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"_id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"type":        &schema.Schema{Type: schema.TypeString, Computed: true},
						"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true},
						"name":        &schema.Schema{Type: schema.TypeString, Computed: true},
						"custrole":    &schema.Schema{Type: schema.TypeString, Computed: true},
						"addresses": &schema.Schema{
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
				},
			},
		},
	}
}
