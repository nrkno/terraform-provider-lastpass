# lastpass_server Resource

## Example Usage

```hcl
resource "lastpass_server" "mysecret" {
    name = "My site"
    username = "foobar"
    password = chomp(file("${path.module}/secret"))
    hostname = "myserver"
    note = <<EOF
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam sed elit nec orci
cursus rhoncus. Morbi lacus turpis, volutpat in lobortis vel, mattis nec magna.
Cras gravida libero vitae nisl iaculis ultrices. Fusce odio ligula, pharetra ac
viverra semper, consequat quis risus.
EOF
}
```

## Argument Reference

* `name` - (Required) Must be unique, and can contain full directory path. Changing name will force recreation.
* `username` - (Required) 
* `password` - (Required) 
* `hostname` - (Optional) 
* `note` - (Optional)

## Attribute Reference

* `fullname`
* `username`
* `password`
* `last_modified_gmt`
* `last_touch`
* `group`
* `note`
* `hostname`

## Importer

Import a pre-existing server secret in Lastpass. Example:

```
terraform import lastpass_server.mysecret 4252909269944373577
```

The ID needs to be a unique numerical value.
