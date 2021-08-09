package lastpass

import (
	"context"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/api"
)

// ResourceSecret describes our lastpass secret resource
func ResourceSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceSecretCreate,
		ReadContext:   ResourceSecretRead,
		UpdateContext: ResourceSecretUpdate,
		DeleteContext: ResourceSecretDelete,
		Importer: &schema.ResourceImporter{
			State: ResourceSecretImporter,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Computed:  true,
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
				Optional: true,
				Computed: true,
				ForceNew: true,
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
func ResourceSecretCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.Client)
	var diags diag.Diagnostics
	s := api.Secret{
		Name:     d.Get("name").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		URL:      d.Get("url").(string),
		Group:    d.Get("group").(string),
		Note:     d.Get("note").(string),
	}
	err := client.Create(&s)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(s.ID)
	ResourceSecretRead(ctx, d, m)

	return diags
}

// ResourceSecretRead is used to sync the local state with the actual state (upstream/lastpass)
func ResourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.Client)
	var diags diag.Diagnostics
	secret, err := client.Read(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}
	d.Set("name", secret.Name)
	d.Set("username", secret.Username)
	d.Set("password", secret.Password)
	d.Set("url", secret.URL)
	d.Set("group", secret.Group)
	d.Set("note", secret.Note)
	d.Set("last_modified_gmt", secret.LastModifiedGmt)
	d.Set("last_touch", secret.LastTouch)
	return diags
}

// ResourceSecretUpdate is used to update our existing resource
func ResourceSecretUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	s := api.Secret{
		ID:       d.Id(),
		Name:     d.Get("name").(string),
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
		URL:      d.Get("url").(string),
		Group:    d.Get("group").(string),
		Note:     d.Get("note").(string),
	}
	client := m.(*api.Client)
	err := client.Update(&s)
	if err != nil {
		return diag.FromErr(err)
	}
	return ResourceSecretRead(ctx, d, m)
}

// ResourceSecretDelete is called to destroy the resource.
func ResourceSecretDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.Client)
	var diags diag.Diagnostics
	err := client.Delete(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

// ResourceSecretImporter is called to import an existing resource.
func ResourceSecretImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if _, err := strconv.Atoi(d.Id()); err != nil {
		err := errors.New("Not a valid Lastpass ID")
		return nil, err
	}
	client := m.(*api.Client)
	secret, err := client.Read(d.Id())
	if err != nil {
		return nil, err
	}
	d.Set("name", secret.Name)
	d.Set("username", secret.Username)
	d.Set("password", secret.Password)
	d.Set("url", secret.URL)
	d.Set("group", secret.Group)
	d.Set("note", secret.Note)
	d.Set("last_modified_gmt", secret.LastModifiedGmt)
	d.Set("last_touch", secret.LastTouch)

	return []*schema.ResourceData{d}, nil
}
