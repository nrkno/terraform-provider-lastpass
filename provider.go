package main

import (
	"errors"

	"github.com/hashicorp/terraform/helper/schema"
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
			"lastpass_record": ResourceRecord(),
		},
		Schema: map[string]*schema.Schema{
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Lastpass login e-mail",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The Lastpass login password",
			},
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	cfg := config{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	if cfg.Username != "" && cfg.Password == "" {
		return nil, errors.New("Lastpass password is not set")
	}
	return cfg, nil
}
