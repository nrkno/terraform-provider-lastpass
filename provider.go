package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider is the root of the lastpass provider
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"lastpass_record": ResourceRecord(),
		},
	}
}
