# lastpass_server Data Source

## Example Usage

```hcl
data "lastpass_server" "myserver" {
    id = "3863267983730403838"
}

resource "aws_db_instance" "myserver" {
  allocated_storage    = 10
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = data.lastpass_server.myserver.hostname
  username             = data.lastpass_server.myserver.username
  password             = data.lastpass_server.myserver.password
}

# data source with custom note template
output "hostname" {
    value = data.lastpass_server.myserver.hostname
}
```

## Argument Reference

* `id` - (Required) Must be unique numerical value.

## Attribute Reference

* `name`
* `fullname`
* `username`
* `password`
* `last_modified_gmt`
* `last_touch`
* `group`
* `note`
* `hostname`
