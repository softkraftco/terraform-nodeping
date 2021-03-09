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
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"nodeping_contact": resourceContact(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"nodeping_contact": dataSourceContact(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	client := nodeping_api_client.NewClient(d.Get("token").(string))
	var diags diag.Diagnostics // TODO: is diags really needed?
	return client, diags
}
