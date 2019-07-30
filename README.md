# terraform-provider-lastpass [![release](https://img.shields.io/github/release/nrkno/terraform-provider-lastpass.svg?style=flat-square)](https://github.com/nrkno/terraform-provider-lastpass/releases/latest) [![Build Status](https://travis-ci.com/nrkno/terraform-provider-lastpass.svg?branch=master)](https://travis-ci.com/nrkno/terraform-provider-lastpass) [![Go Report Card](https://goreportcard.com/badge/github.com/nrkno/terraform-provider-lastpass)](https://goreportcard.com/report/github.com/nrkno/terraform-provider-lastpass)

The Lastpass provider is used to read, manage, or destroy secrets inside Lastpass. Goodbye secret .tfvars files 👋

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="400px">

## Getting started:

1. Install [Terraform](https://www.terraform.io/downloads.html) v0.12 or later.
1. Install the latest pre-compiled binary (Linux/MacOS/Windows) inside `~/.terraform.d/plugins`. Check [releases](https://github.com/nrkno/terraform-provider-lastpass/releases) page.
2. Make sure to have [lastpass-cli](https://github.com/lastpass/lastpass-cli) in your current `$PATH`. 
3. Once a plugin is installed, run `terraform init` to initialize it.

Bonus: 

- Set `LPASS_AGENT_TIMEOUT=86400` inside your `~/.lpass/env` to stay logged in for 24h. Set to `0` to never logout (less secure).
- Set `LASTPASS_USER` and `LASTPASS_PASSWORD` env variables to avoid writing login to your .tf-file.


## Example Usage:

```hcl
provider "lastpass" {
    version = "0.3.0"
    username = "user@example.com"
    password = file("${path.module}/.lpass")
} 

# secret with random generated password
resource "lastpass_secret" "mylogin" {
    name = "My service"
    username = "foobar"
    generate {
        length = 24
        use_symbols = false
    }
}

# secret with password set from string, file, variable, etc.
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

# data source with computed values
data "lastpass_secret" "mydb" {
    id = "3863267983730403838"
}

resource "aws_db_instance" "mydb" {
  allocated_storage    = 10
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = data.lastpass_secret.mydb.user
  password             = data.lastpass_secret.mydb.password
}

```

## Importer

Import a pre-existing secret in Lastpass. Example:

```
terraform import lastpass_secret.mysecret 4252909269944373577
```

The ID needs to be a unique numerical value.

## provider lastpass

* `username` - (Optional) 
  * Can be set via `LASTPASS_USER` env variable.
  * Leave empty if you prefer to login manually.
  * With 2FA enabled you will need to login manually with `--trust` at least once.
* `password` - (Optional) Only required if username is set.
  * Can be set via `LASTPASS_PASSWORD` env variable.


### resource lastpass_secret

**Argument Reference**

* `name` - (Required) Must be unique.
* `username` - (Optional) 
* `password` - (Optional) 
* `generate` - (Optional) Settings for autogenerating password. Either password or generate must be defined.
  * `length` - (Required) The length of the password to generate.
  * `use_symbols` - (Optional) Whether the secret should contain symbols.
* `url` - (Optional) 
* `note` - (Optional)

**Attributes Reference**

* `fullname` - (Computed) 
* `username` - (Computed)
* `password` - (Computed)
* `last_modified_gmt` - (Computed) 
* `last_touch` - (Computed) 
* `group` - (Computed) 
* `url` - (Computed) 
* `note` - (Computed)

### data source lastpass_secret

**Argument Reference**

* `id` - (Required) Must be unique numerical value.

**Attributes Reference**

* `name` - (Computed) 
* `fullname` - (Computed) 
* `username` - (Computed) 
* `password` - (Computed) 
* `last_modified_gmt` - (Computed) 
* `last_touch` - (Computed) 
* `group` - (Computed) 
* `url` - (Computed) 
* `note` - (Computed)