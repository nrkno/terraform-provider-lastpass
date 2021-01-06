package lastpass

import (
    "context"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceSshKey describes lastpass_ssh_key data source
// Reference: https://github.com/lastpass/lastpass-cli/blob/8767b5e53192ad4e72d1352db4aa9218e928cbe1/notes.c#L102-L105
func DataSourceSshKey() *schema.Resource {
    return &schema.Resource{
        ReadContext: DataSourceSshKeyRead,
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
            "pass_phrase": {
                Type:      schema.TypeString,
                Computed:  true,
                Sensitive: true,
            },
            "private_key": {
                Type:      schema.TypeString,
                Computed:  true,
                Sensitive: true,
            },
            "public_key": {
                Type:      schema.TypeString,
                Computed:  true,
                Sensitive: true,
            },
            "hostname": {
                Type:      schema.TypeString,
                Computed:  true,
                Sensitive: true,
            },
            "bit_strength": {
                Type:      schema.TypeString,
                Computed:  true,
            },
            "format": {
                Type:      schema.TypeString,
                Computed:  true,
            },
            "date": {
                Type:      schema.TypeString,
                Computed:  true,
            },
        },
    }
}

// DataSourceSshKeyRead reads resource from upstream/lastpass
func DataSourceSshKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
    d.Set("last_modified_gmt", secrets[0].LastModifiedGmt)
    d.Set("last_touch", secrets[0].LastTouch)
    d.Set("group", secrets[0].Group)
    d.Set("notetype", secrets[0].CustomFields["NoteType"])
    d.Set("pass_phrase", secrets[0].CustomFields["Passphrase"])
    d.Set("private_key", secrets[0].CustomFields["Private Key"])
    d.Set("public_key", secrets[0].CustomFields["Public Key"])
    d.Set("hostname", secrets[0].CustomFields["Hostname"])
    d.Set("bit_strength", secrets[0].CustomFields["Bit Strength"])
    d.Set("format", secrets[0].CustomFields["Format"])
    d.Set("date", secrets[0].CustomFields["Date"])
    d.Set("note", secrets[0].CustomFields["Notes"])
    return diags
}
