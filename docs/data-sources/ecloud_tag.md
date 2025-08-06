# ecloud_tag Data Source

This data source allows you to retrieve information about an eCloud Tag

## Example Usage

### Find Tag by Name

```hcl
data "ecloud_tag" "production" {
  name = "production"
}
```

### Find Tag by ID

```hcl
data "ecloud_tag" "web_server" {
  tag_id = "tag-abcdef12"
}
```

### Find Tag by Name and Scope

```hcl
data "ecloud_tag" "instance_tag" {
  name  = "web-server"
  scope = "instance"
}
```

### Using Tag Data in Resources

```hcl
# Find existing tags
data "ecloud_tag" "environment" {
  name = "production"
}

data "ecloud_tag" "team" {
  name = "platform"
}

# Use existing tags with new instance
resource "ecloud_instance" "new_server" {
  # ... instance configuration ...
  
  tag_ids = [
    data.ecloud_tag.environment.id,
    data.ecloud_tag.team.id
  ]
}
```

### Conditional Tag Assignment

```hcl
# Find tag if it exists
data "ecloud_tag" "monitoring" {
  name = "monitoring-enabled"
}

locals {
  # Create tag IDs list conditionally
  instance_tags = compact([
    data.ecloud_tag.environment.id,
    var.enable_monitoring ? data.ecloud_tag.monitoring.id : null
  ])
}

resource "ecloud_instance" "server" {
  # ... instance configuration ...
  
  tag_ids = local.instance_tags
}
```

## Argument Reference

The following arguments are supported for filtering tags. At least one argument must be specified.

- `tag_id`: ID of the tag
- `name`: Name of the tag
- `scope`: Scope of the tag

## Attributes Reference

- `id`: ID of the tag
- `name`: Name of the tag
- `scope`: Scope of the tag
- `created_at`: Timestamp when the tag was created
- `updated_at`: Timestamp when the tag was last updated
