package main

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/lastpass"
)

// DataSourceSecret describes our lastpass secret data source
func DataSourceSecret() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceSecretRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fullname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"last_modified_gmt": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_touch": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"note": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

// DataSourceSecretRead reads resource from upstream/lastpass
func DataSourceSecretRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*lastpass.Client)
	id := d.Get("id").(string)
	if _, err := strconv.Atoi(id); err != nil {
		err := errors.New("Not a valid Lastpass ID")
		return err
	}
	s, err := client.Read(id)
	if err != nil {
		if s.ID == "0" {
			d.SetId("")
			return nil
		}
		return err
	}
	d.SetId(s.ID)
	d.Set("name", s.Name)
	d.Set("fullname", s.Fullname)
	d.Set("username", s.Username)
	d.Set("password", s.Password)
	d.Set("last_modified_gmt", s.LastModifiedGmt)
	d.Set("last_touch", s.LastTouch)
	d.Set("group", s.Group)
	d.Set("url", s.URL)
	d.Set("note", s.Note)
	return nil
}
