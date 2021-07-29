package lastpass

import (
	"context"
	"errors"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nrkno/terraform-provider-lastpass/api"
)

// DataSourceSecret describes our lastpass secret data source
func DataSourceSecret() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceSecretRead,
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
func DataSourceSecretRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*api.Client)
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	if _, err := strconv.Atoi(id); err != nil {
		err := errors.New("Not a valid Lastpass ID")
		return diag.FromErr(err)
	}
	secret, err := client.Read(id)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(secret.ID)
	d.Set("name", secret.Name)
	d.Set("username", secret.Username)
	d.Set("password", secret.Password)
	d.Set("url", secret.URL)
	d.Set("group", secret.Group)
	d.Set("note", secret.Note)
	d.Set("custom_fields", secret.CustomFields)
	return diags
}
