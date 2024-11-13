# ecloud_backup_gateway_spec Data Source

This data source represents an eCloud Backup Gateway Specification

## Example Usage

```hcl
data "ecloud_backup_gateway_spec" "small" {
  name = "Medium"
}
```

## Argument Reference

- `backup_gateway_specification_id`: ID of backup gateway specification
- `name`: Name of specification

## Attributes Reference

`id` is set to the backup gateway specification ID

- `name`: Name of specification
- `description`: Description of the backup gateway specification