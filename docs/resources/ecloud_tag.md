# ecloud_tag Resource

This resource is for managing eCloud Tags. Tags are metadata labels that can be assigned to various eCloud resources for organization, categorization, and management purposes.

## Example Usage

### Basic Tag

```hcl
resource "ecloud_tag" "my_first_tag" {
  name = "production"
}
```

### Tag with Scope

```hcl
resource "ecloud_tag" "env_prod" {
  scope = "environment"
  name  = "production"
}
```

### Multiple Tags for Resource Organization

```hcl
resource "ecloud_tag" "env_prod" {
  scope = "environment"
  name  = "production"
}

resource "ecloud_tag" "env_dev" {
  scope = "environment"
  name  = "development"
}

resource "ecloud_tag" "team_devops" {
  scope = "team"
  name  = "devops"
}

resource "ecloud_tag" "monitoring_enabled" {
  scope = "monitoring"
  name  = "enabled"
}

# Use tags with an instance
resource "ecloud_instance" "api_server" {
  # ... instance configuration ...
  
  tag_ids = [
    ecloud_tag.env_prod.id,
    ecloud_tag.team_devops.id,
    ecloud_tag.monitoring_enabled.id
  ]
}
```

## Argument Reference

- `name`: (Required) Name of the tag. This is the primary identifier and label for the tag
- `scope`: Scope of the tag. Tag names must be unique within a given scope.

## Attributes Reference

- `id`: ID of the tag
- `name`: Name of the tag
- `scope`: Scope of the tag

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

### Tag Scopes
The `scope` field is optional but recommended. Some examples:

- Use `"environment"` with tag names like `"production"` or `"development"` to indicate the environment
- Use `"team"` for team-related tags, such as `"devops"` or `"engineering"`
- Use `"monitoring"` with values like `"enabled"` or `"disabled"` to indicate monitoring status
- Use `"cost_center"` for financial tracking, such as `"marketing"` or `"sales"`
- Use `"backups"` for backup-related tags, such as `"enabled"` or `"disabled"`

Custom scopes can be defined based on your organizational needs