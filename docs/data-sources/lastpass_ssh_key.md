# lastpass_ssh_key Data Source

## Example Usage

```hcl
locals {
  mykeys = {
    one   = "320356527595751118",
    two   = "1820729297211478707",
    three = "6859419293998842259",
    four  = "2454822705972235622",
  }
}

data lastpass_ssh_key keys {
  for_each = local.mykeys
  id     = each.value
}

output keys {
  value = [ for i in data.lastpass_ssh_key.keys : i ]
}
```

## Argument Reference

* `id` - (Required) Must be unique numerical value.

## Attribute Reference

* `name`
* `fullname`
* `pass_phrase`
* `public_key`
* `private_key`
* `hostname`
* `bit_strength`
* `format`
* `date`
* `last_modified_gmt`
* `last_touch`
* `group`
* `note`
