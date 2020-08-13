# terraform-provider-lastpass 
[![release](https://img.shields.io/github/release/nrkno/terraform-provider-lastpass.svg?style=flat-square)](https://github.com/nrkno/terraform-provider-lastpass/releases/latest) [![Build Status](https://travis-ci.com/nrkno/terraform-provider-lastpass.svg?branch=master)](https://travis-ci.com/nrkno/terraform-provider-lastpass) [![Go Report Card](https://goreportcard.com/badge/github.com/nrkno/terraform-provider-lastpass)](https://goreportcard.com/report/github.com/nrkno/terraform-provider-lastpass) [![GoDoc](https://godoc.org/github.com/github.com/nrkno/terraform-provider-lastpass/lastpass?status.svg)](https://godoc.org/github.com/nrkno/terraform-provider-lastpass/lastpass) [![Dependabot Status](https://api.dependabot.com/badges/status?host=github&repo=nrkno/terraform-provider-lastpass)](https://dependabot.com)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="400px">

The Lastpass provider is used to read, manage, or destroy secrets inside Lastpass. Goodbye secret .tfvars files 👋

```hcl
terraform {
  required_providers {
    lastpass = {
      source = "nrkno/lastpass"
      version = "0.5.0"
    }
  }
}

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

Documentation and examples can be found inside the Terraform registry:

- [Terraform Registry](https://registry.terraform.io/providers/nrkno/lastpass/latest)
- [Documentation](https://registry.terraform.io/providers/nrkno/lastpass/latest/docs)
 
## License

[Apache](LICENSE)
