package lastpass

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/api"
)

// Provider config
type config struct {
	Username string
	Password string
}

// Provider is the root of the lastpass provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"lastpass_secret": ResourceSecret(),
			"lastpass_server": ResourceServer(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"lastpass_secret": DataSourceSecret(),
			"lastpass_server": DataSourceServer(),
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Lastpass login e-mail",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_USER", nil),
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Lastpass login password",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_PASSWORD", nil),
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := api.Client{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	return &client, diags
}
