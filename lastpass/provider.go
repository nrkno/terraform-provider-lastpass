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
	BaseURL  string
}

// Provider is the root of the lastpass provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"lastpass_secret": ResourceSecret(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"lastpass_secret": DataSourceSecret(),
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
			"enable_2fa": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				Description: "enable two-factor authentication: LastPass Authenticator, Google Authenticator, Microsoft Authenticator, YubiKey, Transakt, Duo Security, or Sesame",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_ENABLE_2FA", true),
			},
			"onetime_password": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Sensitive:   true,
				Description: "one-time password for 2fa",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_ONETIME_PASSWORD", ""),
			},
			"trust": {
				Type:        schema.TypeBool,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				Description: "Trust will cause subsequent logins to not require multifactor authentication",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_TRUST", true),
			},
			"configdir": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				Description: "Path where the provider stores its trust ID",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_CONFIGDIR", ".terraform-provider-lastpass"),
			},
			"baseurl": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Sensitive:   false,
				Description: "Base URL https://lastpass.com or https://lastpass.eu",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_BASEURL", "https://lastpass.com"),
			},
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics
	client := api.Client{
		Username:  d.Get("username").(string),
		Password:  d.Get("password").(string),
		Trust:     d.Get("trust").(bool),
		TwoFA:     d.Get("enable_2fa").(bool),
		OnetimePW: d.Get("onetime_password").(string),
		ConfigDIR: d.Get("configdir").(string),
		BaseURL:   d.Get("baseurl").(string),
	}

	err := client.Login()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create lastpass client",
			Detail:   err.Error(),
		})
		return nil, diags
	}
	err = client.Sync()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to fetch secrets from Lastpass",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return &client, diags
}
