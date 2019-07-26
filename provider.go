package main

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/lastpass"
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
			"lastpass_secret": ResourceRecord(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"lastpass_secret": DataSourceRecord(),
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Lastpass login e-mail",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_USER", ""),
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Lastpass login password",
				DefaultFunc: schema.EnvDefaultFunc("LASTPASS_PASSWORD", ""),
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	if d.Get("username").(string) != "" && d.Get("password").(string) == "" {
		return nil, errors.New("lastpass password is not set")
	}
	client := lastpass.Client{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	return &client, nil
}
