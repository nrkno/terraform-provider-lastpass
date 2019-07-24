package main

import (
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
	d.Set("url", r.URL)
	d.Set("username", r.Username)
	d.Set("password", r.Password)
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
