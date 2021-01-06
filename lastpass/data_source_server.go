package lastpass

import (
    "context"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceServer describes lastpass_server data source
// Reference: https://github.com/lastpass/lastpass-cli/blob/8767b5e53192ad4e72d1352db4aa9218e928cbe1/notes.c#L90-L93
func DataSourceServer() *schema.Resource {
    return &schema.Resource{
        ReadContext: DataSourceServerRead,
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
                Type:      schema.TypeString,
                Computed:  true,
                Sensitive: true,
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
            "note": {
                Type:      schema.TypeString,
                Computed:  true,
                Sensitive: true,
            },
            "notetype": {
                Type:      schema.TypeString,
                Computed:  true,
            },
            "hostname": {
                Type:      schema.TypeString,
                Computed:  true,
                Sensitive: true,
            },
        },
    }
}

// DataSourceServerRead reads resource from upstream/lastpass
func DataSourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
    d.Set("username", secrets[0].CustomFields["Username"])
    d.Set("password", secrets[0].CustomFields["Password"])
    d.Set("last_modified_gmt", secrets[0].LastModifiedGmt)
    d.Set("last_touch", secrets[0].LastTouch)
    d.Set("group", secrets[0].Group)
    d.Set("note", secrets[0].CustomFields["Notes"])
    d.Set("notetype", secrets[0].CustomFields["NoteType"])
    d.Set("hostname", secrets[0].CustomFields["Hostname"])
    return diags
}
