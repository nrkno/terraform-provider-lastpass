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
	var diags diag.Diagnostics
	secrets, err := dataSourceSecretRead(m, d.Get("id").(string) )

	if err != nil {
		return diag.FromErr(err)
	} else if len(secrets) == 0 {
		d.SetId("")
		return diags
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
	return diags
}

func dataSourceSecretRead(m interface{}, id string) ([]api.Secret, error) {
    if _, err := strconv.Atoi(id); err != nil {
        err := errors.New("Not a valid Lastpass ID")
        return []api.Secret{}, err
    }
    client := m.(*api.Client)
    secrets, err := client.Read(id)
    if err != nil {
        return []api.Secret{}, err
    } else if len(secrets) > 1 {
        var err = errors.New("got duplicate IDs")
        return secrets, err
    }
    return secrets, nil
}
