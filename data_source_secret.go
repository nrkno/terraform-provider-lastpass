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
			"custom_fields": {
				Type:      schema.TypeMap,
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
	secrets, err := client.Read(id)
	if err != nil {
		return err
	}
	if len(secrets) == 0 {
		d.SetId("")
		return nil
	} else if len(secrets) > 1 {
		var err = errors.New("got duplicate IDs")
		return err
	}
	d.SetId(secrets[0].ID)
	d.Set("name", secrets[0].Name)
	d.Set("fullname", secrets[0].Fullname)
	d.Set("username", secrets[0].Username)
	d.Set("password", secrets[0].Password)
	d.Set("last_modified_gmt", secrets[0].LastModifiedGmt)
	d.Set("last_touch", secrets[0].LastTouch)
	d.Set("group", secrets[0].Group)
	d.Set("url", secrets[0].URL)
	d.Set("note", secrets[0].Note)
	d.Set("custom_fields", secrets[0].CustomFields)
	return nil
}
