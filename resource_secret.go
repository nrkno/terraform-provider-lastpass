package main

import (
	"errors"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/lastpass"
	"github.com/sethvargo/go-password/password"
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
				ForceNew: true,
			},
			"fullname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"generate"},
				Optional:      true,
				Sensitive:     true,
				Computed:      true,
				Description:   "The password content. Either `password` or `generate` must be defined.",
			},
			"generate": {
				Type:          schema.TypeList,
				MaxItems:      1,
				ConflictsWith: []string{"password"},
				Optional:      true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the password.",
						},
						"use_symbols": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the secret should contain symbols.",
						},
					},
				},
				Description: "Settings for autogenerating a password. Either `password` or `generate` must be defined.",
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
				Computed: true,
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Computed:    true,
				Description: "The secret note content.",
			},
		},
	}
}

// ResourceSecretCreate is used to create a new resource and generate ID.
func ResourceSecretCreate(d *schema.ResourceData, m interface{}) error {
	generate := d.Get("generate").([]interface{})
	if d.Get("password") == "" && len(generate) == 0 {
		return errors.New("either 'password' or 'generate' must be specified")
	}
	client := m.(*lastpass.Client)
	s := lastpass.Secret{
		Name:     d.Get("name").(string),
		URL:      d.Get("url").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		Note:     d.Get("note").(string),
	}
	if len(generate) == 1 {
		settings := generate[0].(map[string]interface{})
		symbols := settings["use_symbols"].(bool)
		length := settings["length"].(int)
		nrSymbols := length + 1 // no symbols by default
		if symbols {
			nrSymbols = 4
		}
		pw, err := password.Generate(length, (length / 4), (length / nrSymbols), false, false)
		if err != nil {
			return err
		}
		s.Password = pw
	}
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
	secrets, err := client.Read(d.Id())
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
	d.Set("name", secrets[0].Name)
	d.Set("fullname", secrets[0].Fullname)
	d.Set("username", secrets[0].Username)
	d.Set("password", secrets[0].Password)
	d.Set("last_modified_gmt", secrets[0].LastModifiedGmt)
	d.Set("last_touch", secrets[0].LastTouch)
	d.Set("group", secrets[0].Group)
	d.Set("url", secrets[0].URL)
	d.Set("note", secrets[0].Note)
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
	secrets, err := client.Read(d.Id())
	if err != nil {
		return nil, err
	}
	if len(secrets) == 0 {
		var err = errors.New("ID not found.")
		return nil, err
	} else if len(secrets) > 1 {
		var err = errors.New("got duplicate IDs")
		return nil, err
	}
	d.Set("name", secrets[0].Name)
	d.Set("fullname", secrets[0].Fullname)
	d.Set("username", secrets[0].Username)
	d.Set("password", secrets[0].Password)
	d.Set("last_modified_gmt", secrets[0].LastModifiedGmt)
	d.Set("last_touch", secrets[0].LastTouch)
	d.Set("group", secrets[0].Group)
	d.Set("url", secrets[0].URL)
	d.Set("note", secrets[0].Note)
	return []*schema.ResourceData{d}, nil
}
