# lastpass_secret Resource

## Example Usage

```hcl
resource "lastpass_secret" "mysecret" {
    name = "My site"
    username = "foobar"
    password = file("${path.module}/secret")
    url = "https://example.com"
    note = <<EOF
Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nam sed elit nec orci
cursus rhoncus. Morbi lacus turpis, volutpat in lobortis vel, mattis nec magna.
Cras gravida libero vitae nisl iaculis ultrices. Fusce odio ligula, pharetra ac
viverra semper, consequat quis risus.
EOF
}
```

## Argument Reference

* `name` - (Required) Must be unique, and can contain full directory path. Changing name will force recreation, and generated passwords will change.
* `username` - (Optional) 
* `password` - (Optional) 
* `generate` - (Optional) Settings for autogenerating password. Either password or generate must be defined.
  * `length` - (Required) The length of the password to generate.
  * `use_symbols` - (Optional) Whether the secret should contain symbols.
* `url` - (Optional) 
* `note` - (Optional)

## Attribute Reference

* `fullname`
* `username`
* `password`
* `last_modified_gmt`
* `last_touch`
* `group`
* `url`
* `note`

## Importer

Import a pre-existing secret in Lastpass. Example:

```
terraform import lastpass_secret.mysecret 4252909269944373577
```

The ID needs to be a unique numerical value.