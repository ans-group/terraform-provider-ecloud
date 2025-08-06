# ecloud_instance Data Source

This data source allows you to retrieve information about an eCloud Instance

## Example Usage

### Basic Usage

```hcl
data "ecloud_instance" "instance-1" {
  name = "my-instance"
}
```

### Accessing Instance Tags

```hcl
data "ecloud_instance" "web-server" {
  instance_id = "i-abcdef12"
}

# Output the tags for examination
output "instance_tags" {
  value = data.ecloud_instance.web-server.tags
}

# Use tag information in locals
locals {
  environment_tag = [
    for tag in data.ecloud_instance.web-server.tags : tag.name
    if tag.scope == "environment"
  ][0]
}
```

## Argument Reference

- `instance_id`: ID of instance
- `vpc_id`: ID of instance VPC
- `name`: Name of instance

## Attributes Reference

`id` is set to instance ID

- `vpc_id`: ID of instance VPC
- `name`: Name of instance
- `tags`: Set of tags assigned to the instance. Each tag contains:
  - `id`: ID of the tag
  - `name`: Name of the tag
  - `scope`: Scope of the tag