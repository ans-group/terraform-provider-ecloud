# ecloud_tag Resource

This resource is for managing eCloud Tags. Tags are metadata labels that can be assigned to various eCloud resources for organization, categorization, and management purposes.

## Example Usage

### Basic Tag

```hcl
resource "ecloud_tag" "environment" {
  name = "production"
}
```

### Tag with Scope

```hcl
resource "ecloud_tag" "web_server" {
  name  = "web-server"
  scope = "instance"
}
```

### Multiple Tags for Resource Organization

```hcl
resource "ecloud_tag" "environment" {
  name  = "production"
  scope = "instance"
}

resource "ecloud_tag" "team" {
  name  = "platform"
  scope = "instance"
}

resource "ecloud_tag" "application" {
  name  = "api-gateway"
  scope = "instance"
}

# Use tags with an instance
resource "ecloud_instance" "api_server" {
  # ... instance configuration ...
  
  tag_ids = [
    ecloud_tag.environment.id,
    ecloud_tag.team.id,
    ecloud_tag.application.id
  ]
}
```

## Argument Reference

- `name`: (Required) Name of the tag. This is the primary identifier and label for the tag
- `scope`: Scope of the tag, indicating what type of resources this tag is intended for (e.g., "instance", "vpc", "network")

## Attributes Reference

- `id`: ID of the tag
- `name`: Name of the tag
- `scope`: Scope of the tag
- `created_at`: Timestamp when the tag was created
- `updated_at`: Timestamp when the tag was last updated

## Import

Tags can be imported using the tag ID:

```bash
terraform import ecloud_tag.example tag-abcdef12
```

## Notes on Tag Management

### Tag Assignment
Tags are assigned to resources using the resource's `tag_ids` attribute. When updating tag assignments on resources:
- The complete list of tag IDs must be provided
- Any tags not included in the list will be removed from the resource
- Tags themselves are not deleted when removed from resources

### Tag Naming
- Tag names should be descriptive and follow your organization's naming conventions
- Consider using consistent naming patterns across your infrastructure
- Tag names are case-sensitive

### Tag Scopes
The `scope` field is optional but recommended for organizing tags by resource type:
- Use `"instance"` for tags intended for instances
- Use `"vpc"` for VPC-related tags
- Use `"network"` for network-related tags
- Custom scopes can be defined based on your organizational needs