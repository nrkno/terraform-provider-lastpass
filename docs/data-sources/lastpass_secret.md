# lastpass_secret Data Source

## Example Usage

```hcl
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
  username             = data.lastpass_secret.mydb.username
  password             = data.lastpass_secret.mydb.password
}

# data source with custom note template 
output "custom_field" {
    value = data.lastpass_secret.mydb.custom_fields.host
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
* `url`
* `note`
* `custom_fields`