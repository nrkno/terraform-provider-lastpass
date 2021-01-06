package lastpass

import (
    "fmt"
    "context"
    "errors"
    "strconv"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
    "github.com/rezroo/terraform-provider-lastpass/api"
)

// ResourceSshKey describes node-type lastpass_ssh_key resource
// Reference: https://github.com/lastpass/lastpass-cli/blob/8767b5e53192ad4e72d1352db4aa9218e928cbe1/notes.c#L102-L105
func ResourceSshKey() *schema.Resource {
    return &schema.Resource{
        CreateContext: ResourceSshKeyCreate,
        ReadContext:   ResourceSshKeyRead,
        UpdateContext: ResourceSshKeyUpdate,
        DeleteContext: ResourceSshKeyDelete,
        Importer: &schema.ResourceImporter{
            State: ResourceSshKeyImporter,
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
                Type:        schema.TypeString,
                Optional:    true,
                //Sensitive:   true,
                Computed:    true,
                Description: "The secret note content.",
            },
            "notetype": {
                Type:      schema.TypeString,
                Computed:  true,
            },
            "pass_phrase": {
                Type:      schema.TypeString,
                Optional: true,
                Sensitive: true,
            },
            "private_key": {
                Type:     schema.TypeString,
                Required: true,
                Sensitive: true,
            },
            "public_key": {
                Type:      schema.TypeString,
                Required:  true,
                Sensitive: true,
            },
            "hostname": {
                Type:      schema.TypeString,
                Optional: true,
                //Computed:  true,
                //Sensitive: true,
            },
            "bit_strength": {
                Type:      schema.TypeString,
                Optional: true,
                Computed:  true,
            },
            "format": {
                Type:      schema.TypeString,
                Optional: true,
                Computed:  true,
            },
            "date": {
                Type:      schema.TypeString,
                Optional: true,
                Computed:  true,
            },
        },
    }
}

func getSshKeyTemplate(date string, hostname string, pubkey string, prvkey string, phrase string, format string, strength string, notes string) string {
    template := fmt.Sprintf(`Date: %s
Hostname: %s
Public Key: %s
Private Key: %s
Passphrase: %s
Format: %s
Bit Strength: %s
Notes:    # Add notes below this line.
%s`, date, hostname, pubkey, prvkey, phrase, format, strength, notes)
    return template
}

// ResourceSshKeyCreate is used to create a new resource and generate ID.
func ResourceSshKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*api.Client)
    var diags diag.Diagnostics
    template := getSshKeyTemplate(
        d.Get("date").(string),
        d.Get("hostname").(string),
        d.Get("public_key").(string),
        d.Get("private_key").(string),
        d.Get("pass_phrase").(string),
        d.Get("format").(string),
        d.Get("bit_strength").(string),
        d.Get("note").(string) )
    s, err := client.CreateNodeType(d.Get("name").(string), template, "ssh-key")
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(s.ID)
    ResourceSshKeyRead(ctx, d, m)
    return diags
}

// ResourceSshKeyRead is used to sync the local state with the actual state (upstream/lastpass)
func ResourceSshKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    var diags diag.Diagnostics
    secrets, err := dataSourceSecretRead(m, d.Id())

    if err != nil {
        return diag.FromErr(err)
    }
    if len(secrets) == 0 {
        d.SetId("")
        return nil
    }
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

// ResourceSshKeyUpdate is used to update our existing resource
func ResourceSshKeyUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    template := getSshKeyTemplate(
        d.Get("date").(string),
        d.Get("hostname").(string),
        d.Get("public_key").(string),
        d.Get("private_key").(string),
        d.Get("pass_phrase").(string),
        d.Get("format").(string),
        d.Get("bit_strength").(string),
        d.Get("note").(string) )
    client := m.(*api.Client)
    err := client.UpdateNodeType(d.Id(), template)
    if err != nil {
        return diag.FromErr(err)
    }
    return ResourceSshKeyRead(ctx, d, m)
}

// ResourceSshKeyDelete is called to destroy the resource.
func ResourceSshKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*api.Client)
    var diags diag.Diagnostics
    err := client.Delete(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }
    return diags
}

// ResourceSshKeyImporter is called to import an existing resource.
func ResourceSshKeyImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
    if _, err := strconv.Atoi(d.Id()); err != nil {
        err := errors.New("Not a valid Lastpass ID")
        return nil, err
    }
    client := m.(*api.Client)
    secrets, err := client.Read(d.Id())
    if err != nil {
        return nil, err
    }
    if len(secrets) == 0 {
        var err = errors.New("ID not found")
        return nil, err
    } else if len(secrets) > 1 {
        var err = errors.New("got duplicate IDs")
        return nil, err
    }
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
    return []*schema.ResourceData{d}, nil
}
