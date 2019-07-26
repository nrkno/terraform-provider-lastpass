package main

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/lastpass"
)

// ResourceRecord describes our lastpass record resource
func ResourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: ResourceRecordCreate,
		Read:   ResourceRecordRead,
		Update: ResourceRecordUpdate,
		Delete: ResourceRecordDelete,
		Importer: &schema.ResourceImporter{
			State: ResourceRecordImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"fullname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
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
				Optional: true,
			},
			"note": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// ResourceRecordCreate is used to create a new resource and generate ID.
func ResourceRecordCreate(d *schema.ResourceData, m interface{}) error {
	r := lastpass.Record{
		Name:     d.Get("name").(string),
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Note:     d.Get("note").(string),
	}
	client := m.(*lastpass.Client)
	r, err := client.Create(r)
	if err != nil {
		return err
	}
	d.SetId(r.ID)
	return ResourceRecordRead(d, m)
}

// ResourceRecordRead is used to sync the local state with the actual state (upstream/lastpass)
func ResourceRecordRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*lastpass.Client)
	r, err := client.Read(d.Id())
	if err != nil {
		if r.ID == "0" {
			d.SetId("")
			return nil
		}
		return err
	}
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

// ResourceRecordUpdate is used to update our existing resource
func ResourceRecordUpdate(d *schema.ResourceData, m interface{}) error {
	r := lastpass.Record{
		Name:     d.Get("name").(string),
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Note:     d.Get("note").(string),
		ID:       d.Id(),
	}
	client := m.(*lastpass.Client)
	err := client.Update(r)
	if err != nil {
		return err
	}
	return ResourceRecordRead(d, m)
}

// ResourceRecordDelete is called to destroy the resource.
func ResourceRecordDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*lastpass.Client)
	err := client.Delete(d.Id())
	if err != nil {
		return err
	}
	return nil
}

// ResourceRecordImporter is called to import an existing resource.
func ResourceRecordImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if _, err := strconv.Atoi(d.Id()); err != nil {
		err := errors.New("Not a valid Lastpass ID")
		return nil, err
	}
	client := m.(*lastpass.Client)
	r, err := client.Read(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("name", r.Name)
	d.Set("fullname", r.Fullname)
	d.Set("username", r.Username)
	d.Set("password", r.Password)
	d.Set("last_modified_gmt", r.LastModifiedGmt)
	d.Set("last_touch", r.LastTouch)
	d.Set("group", r.Group)
	d.Set("url", r.URL)
	d.Set("note", r.Note)
	return []*schema.ResourceData{d}, nil
}
