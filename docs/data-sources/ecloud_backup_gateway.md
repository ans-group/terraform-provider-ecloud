# ecloud_backup_gateway Data Source

This data source represents an eCloud Backup Gateway

## Example Usage

```hcl
data "ecloud_backup_gateway" "gateway-1" {
  name = "my-backup-gateway"
}
```

## Argument Reference

- `backup_gateway_id`: ID of backup gateway
- `vpc_id`: ID of VPC
- `name`: Name of backup gateway
- `availability_zone_id`: ID of availability zone

## Attributes Reference

`id` is set to the backup gateway ID

- `vpc_id`: ID of VPC
- `name`: Name of backup gateway
- `availability_zone_id`: ID of availability zone
- `gateway_spec_id`: ID of backup gateway specification