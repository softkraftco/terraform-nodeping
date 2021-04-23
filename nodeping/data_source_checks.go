package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func dataSourceCheck() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCheckRead,
		Schema: map[string]*schema.Schema{
			"id":          &schema.Schema{Type: schema.TypeString, Required: true},
			"customer_id": &schema.Schema{Type: schema.TypeString, Computed: true},
			"label":       &schema.Schema{Type: schema.TypeString, Computed: true},
			"interval":    &schema.Schema{Type: schema.TypeInt, Computed: true},
			"notifications": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"contact":  &schema.Schema{Type: schema.TypeString, Required: true},
						"delay":    &schema.Schema{Type: schema.TypeInt, Required: true},
						"schedule": &schema.Schema{Type: schema.TypeString, Required: true},
					},
				},
			},
			"type":         &schema.Schema{Type: schema.TypeString, Computed: true},
			"target":       &schema.Schema{Type: schema.TypeString, Computed: true},
			"status":       &schema.Schema{Type: schema.TypeString, Computed: true},
			"created":      &schema.Schema{Type: schema.TypeInt, Computed: true},
			"modified":     &schema.Schema{Type: schema.TypeInt, Computed: true},
			"enabled":      &schema.Schema{Type: schema.TypeString, Computed: true}, // called "enable" in API response
			"public":       &schema.Schema{Type: schema.TypeBool, Computed: true},
			"runlocations": &schema.Schema{Type: schema.TypeSet, Optional: true, Computed: true, Elem: &schema.Schema{Type: schema.TypeString}},
			"homeloc":      &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"threshold":    &schema.Schema{Type: schema.TypeInt, Optional: true, Computed: true},
			"sens":         &schema.Schema{Type: schema.TypeInt, Optional: true, Computed: true},
			"queue":        &schema.Schema{Type: schema.TypeString, Computed: true},
			"uuid":         &schema.Schema{Type: schema.TypeString, Computed: true},
			"state":        &schema.Schema{Type: schema.TypeInt, Computed: true},
			"firstdown":    &schema.Schema{Type: schema.TypeInt, Computed: true},
			"dep":          &schema.Schema{Type: schema.TypeString, Optional: true},
			"description":  &schema.Schema{Type: schema.TypeString, Optional: true},
			// the following are optional and stored in responses "parameters" dictionary:
			"checktoken":    &schema.Schema{Type: schema.TypeString, Computed: true},
			"clientcert":    &schema.Schema{Type: schema.TypeString, Computed: true, Default: nil},
			"contentstring": &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"dohdot":        &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"dnstype":       &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"dnstoresolve":  &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"dnsrd":         &schema.Schema{Type: schema.TypeBool, Optional: true, Computed: true},
			"transport":     &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"follow":        &schema.Schema{Type: schema.TypeBool, Optional: true, Computed: true},
			"email":         &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"port":          &schema.Schema{Type: schema.TypeInt, Optional: true, Computed: true},
			"username":      &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"password":      &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"secure":        &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"verify":        &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"ignore":        &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"invert":        &schema.Schema{Type: schema.TypeBool, Optional: true, Computed: true},
			"warningdays":   &schema.Schema{Type: schema.TypeInt, Optional: true, Computed: true},
			"fields": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true, Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id":   &schema.Schema{Type: schema.TypeString, Optional: true},
						"name": &schema.Schema{Type: schema.TypeString, Required: true},
						"min":  &schema.Schema{Type: schema.TypeInt, Required: true},
						"max":  &schema.Schema{Type: schema.TypeInt, Required: true},
					},
				},
			},
			"postdata":       &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"data":           &schema.Schema{Type: schema.TypeMap, Optional: true, Computed: true},
			"receiveheaders": &schema.Schema{Type: schema.TypeMap, Optional: true, Computed: true},
			"sendheaders":    &schema.Schema{Type: schema.TypeMap, Optional: true, Computed: true},
			"edns":           &schema.Schema{Type: schema.TypeMap, Optional: true, Computed: true},
			"method":         &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"statuscode":     &schema.Schema{Type: schema.TypeInt, Optional: true, Computed: true},
			"ipv6":           &schema.Schema{Type: schema.TypeBool, Optional: true, Computed: true},
			"regex":          &schema.Schema{Type: schema.TypeBool, Optional: true, Computed: true},
			"servername":     &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"snmpv":          &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"snmpcom":        &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
			"verifyvolume":   &schema.Schema{Type: schema.TypeBool, Optional: true, Computed: true},
			"volumemin":      &schema.Schema{Type: schema.TypeInt, Optional: true, Computed: true},
			"whoisserver":    &schema.Schema{Type: schema.TypeString, Optional: true, Computed: true},
		},
	}
}

func dataSourceCheckRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*nodeping_api_client.Client)

	check, err := client.GetCheck(ctx, d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	applyCheckToSchema(check, d)

	return nil
}
