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

// ResourceServer describes node-type lastpass_server resource
// Reference: https://github.com/lastpass/lastpass-cli/blob/8767b5e53192ad4e72d1352db4aa9218e928cbe1/notes.c#L90-L93
func ResourceServer() *schema.Resource {
    return &schema.Resource{
        CreateContext: ResourceServerCreate,
        ReadContext:   ResourceServerRead,
        UpdateContext: ResourceServerUpdate,
        DeleteContext: ResourceServerDelete,
        Importer: &schema.ResourceImporter{
            State: ResourceServerImporter,
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
                Required: true,
                Sensitive: true,
                //Computed: true,
            },
            "password": {
                Type:      schema.TypeString,
                Required:  true,
                Sensitive: true,
                //Computed:  true,
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
            "hostname": {
                Type:      schema.TypeString,
                Optional: true,
                //Computed:  true,
                //Sensitive: true,
            },
        },
    }
}

func getServerTemplate(hostname string, username string, password string, notes string) string {
    template := fmt.Sprintf(`Hostname: %s
Username: %s
Password: %s
Notes:    # Add notes below this line.
%s`, hostname, username, password, notes)
    return template
}

// ResourceServerCreate is used to create a new resource and generate ID.
func ResourceServerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*api.Client)
    var diags diag.Diagnostics
    template := getServerTemplate(
        d.Get("hostname").(string),
        d.Get("username").(string),
        d.Get("password").(string),
        d.Get("note").(string) )
    s, err := client.CreateNodeType(d.Get("name").(string), template, "server")
    if err != nil {
        return diag.FromErr(err)
    }
    d.SetId(s.ID)
    ResourceServerRead(ctx, d, m)
    return diags
}

// ResourceServerRead is used to sync the local state with the actual state (upstream/lastpass)
func ResourceServerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

// ResourceServerUpdate is used to update our existing resource
func ResourceServerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    template := getServerTemplate(
        d.Get("hostname").(string),
        d.Get("username").(string),
        d.Get("password").(string),
        d.Get("note").(string) )
    client := m.(*api.Client)
    err := client.UpdateNodeType(d.Id(), template)
    if err != nil {
        return diag.FromErr(err)
    }
    return ResourceServerRead(ctx, d, m)
}

// ResourceServerDelete is called to destroy the resource.
func ResourceServerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    client := m.(*api.Client)
    var diags diag.Diagnostics
    err := client.Delete(d.Id())
    if err != nil {
        return diag.FromErr(err)
    }
    return diags
}

// ResourceServerImporter is called to import an existing resource.
func ResourceServerImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
    d.Set("username", secrets[0].CustomFields["Username"])
    d.Set("password", secrets[0].CustomFields["Password"])
    d.Set("last_modified_gmt", secrets[0].LastModifiedGmt)
    d.Set("last_touch", secrets[0].LastTouch)
    d.Set("group", secrets[0].Group)
    d.Set("note", secrets[0].CustomFields["Notes"])
    d.Set("notetype", secrets[0].CustomFields["NoteType"])
    d.Set("hostname", secrets[0].CustomFields["Hostname"])
    return []*schema.ResourceData{d}, nil
}
