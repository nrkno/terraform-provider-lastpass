# lastpass_server Resource

## Example Usage

```hcl
resource "lastpass_ssh_key" "mysshkey" {
    name = "My key"
    pass_phrase = each.value[1]
    public_key  = chomp( file("key.pub") )
    private_key = chomp( file("key.pem") )
    hostname    = "myserver"
    note        = <<EOF
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam sed elit nec orci
cursus rhoncus. Morbi lacus turpis, volutpat in lobortis vel, mattis nec magna.
Cras gravida libero vitae nisl iaculis ultrices. Fusce odio ligula, pharetra ac
viverra semper, consequat quis risus.
EOF
}
```

## Argument Reference

* `name` - (Required) Must be unique, and can contain full directory path. Changing name will force recreation.
* `public_key` - (Required) 
* `private_key` - (Required) 
* `pass_phrase` - (Optional)
* `hostname` - (Optional) 
* `bit_strength` - (Optional)
* `format` - (Optional) 
* `date` - (Optional) 
* `note` - (Optional)

## Additional ttribute Reference

* `fullname`
* `last_modified_gmt`
* `last_touch`
* `group`

## Importer

Import a pre-existing server secret in Lastpass. Example:

```
terraform import lastpass_server.mysshkey 4252909269944373577
```

The ID needs to be a unique numerical value.
