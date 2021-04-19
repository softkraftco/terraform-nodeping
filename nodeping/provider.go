package nodeping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"terraform-nodeping/nodeping_api_client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("NODEPING_API_TOKEN", nil),
				Description: "NodePing API token - used for an authentication against NodePing API.",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"nodeping_contact":  resourceContact(),
			"nodeping_check":    resourceCheck(),
			"nodeping_schedule": resourceSchedule(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nodeping_contact": dataSourceContact(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := nodeping_api_client.NewClient(d.Get("token").(string))
	return client, nil
}
