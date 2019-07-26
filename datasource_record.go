package main

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/lastpass"
)

// DataSourceRecord describes our lastpass record data source
func DataSourceRecord() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceRecordRead,
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

// DataSourceRecordRead reads resource from upstream/lastpass
func DataSourceRecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*lastpass.Client)
	id := d.Get("id").(string)
	if _, err := strconv.Atoi(id); err != nil {
		err := errors.New("Not a valid Lastpass ID")
		return err
	}
	r, err := client.Read(id)
	if err != nil {
		if r.ID == "0" {
			d.SetId("")
			return nil
		}
		return err
	}
	d.SetId(r.ID)
	d.Set("name", r.Name)
	d.Set("fullname", r.Fullname)
	d.Set("username", r.Username)
	d.Set("password", r.Password)
	d.Set("last_modified_gmt", r.LastModifiedGmt)
	d.Set("last_touch", r.LastTouch)
	d.Set("group", r.Group)
	d.Set("url", r.URL)
	d.Set("note", r.Note)
	return nil
}
