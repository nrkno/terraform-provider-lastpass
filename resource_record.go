package main

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/lastpass"
)

// ResourceSecret describes our lastpass secret resource
func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		Create: ResourceSecretCreate,
		Read:   ResourceSecretRead,
		Update: ResourceSecretUpdate,
		Delete: ResourceSecretDelete,
		Importer: &schema.ResourceImporter{
			State: ResourceSecretImporter,
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

// ResourceSecretCreate is used to create a new resource and generate ID.
func ResourceSecretCreate(d *schema.ResourceData, m interface{}) error {
	s := lastpass.Secret{
		Name:     d.Get("name").(string),
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Note:     d.Get("note").(string),
	}
	client := m.(*lastpass.Client)
	s, err := client.Create(s)
	if err != nil {
		return err
	}
	d.SetId(s.ID)
	return ResourceSecretRead(d, m)
}

// ResourceSecretRead is used to sync the local state with the actual state (upstream/lastpass)
func ResourceSecretRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*lastpass.Client)
	s, err := client.Read(d.Id())
	if err != nil {
		if s.ID == "0" {
			d.SetId("")
			return nil
		}
		return err
	}
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

// ResourceSecretUpdate is used to update our existing resource
func ResourceSecretUpdate(d *schema.ResourceData, m interface{}) error {
	s := lastpass.Secret{
		Name:     d.Get("name").(string),
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Note:     d.Get("note").(string),
		ID:       d.Id(),
	}
	client := m.(*lastpass.Client)
	err := client.Update(s)
	if err != nil {
		return err
	}
	return ResourceSecretRead(d, m)
}

// ResourceSecretDelete is called to destroy the resource.
func ResourceSecretDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*lastpass.Client)
	err := client.Delete(d.Id())
	if err != nil {
		return err
	}
	return nil
}

// ResourceSecretImporter is called to import an existing resource.
func ResourceSecretImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if _, err := strconv.Atoi(d.Id()); err != nil {
		err := errors.New("Not a valid Lastpass ID")
		return nil, err
	}
	client := m.(*lastpass.Client)
	s, err := client.Read(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("name", s.Name)
	d.Set("fullname", s.Fullname)
	d.Set("username", s.Username)
	d.Set("password", s.Password)
	d.Set("last_modified_gmt", s.LastModifiedGmt)
	d.Set("last_touch", s.LastTouch)
	d.Set("group", s.Group)
	d.Set("url", s.URL)
	d.Set("note", s.Note)
	return []*schema.ResourceData{d}, nil
}
