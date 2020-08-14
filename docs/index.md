# Lastpass Provider

The Lastpass provider is used to read, manage, or destroy secrets inside Lastpass. Goodbye secret .tfvars files 👋

Make sure to have [lastpass-cli](https://github.com/lastpass/lastpass-cli) in your current `$PATH`. 

-> Set `LPASS_AGENT_TIMEOUT=86400` inside your `~/.lpass/env` to stay logged in for 24h. Set to `0` to never logout (less secure).

-> Set `LASTPASS_USER` and `LASTPASS_PASSWORD` env variables to avoid writing login to your .tf-files.

## Example Usage

```hcl

resource "random_password" "pw" {
  length = 32
  special = false
}

resource "lastpass_secret" "mylogin" {
    name = "My service"
    username = "foobar"
    password = random_password.pw.result
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

## Argument Reference

* `username` - (Required) 
  * Can be set via `LASTPASS_USER` env variable.
  * Can be set to empty string for manual lpass login.
  * With 2FA enabled you will need to login manually with `--trust` at least once.
* `password` - (Required)
  * Can be set via `LASTPASS_PASSWORD` env variable.
  * Can be set to empty string for manual lpass login.
