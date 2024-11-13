# ecloud_backup_gateway Resource

This resource is for managing eCloud Backup Gateways

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
  region_id = "reg-abcdef12"
  name = "example-vpc"
}

data "ecloud_availability_zone" "az-man4" {
  name = "Manchester West"
}

data "ecloud_backup_gateway_spec" "spec-1" {
  name = "Medium"
}

resource "ecloud_backup_gateway" "gateway-1" {
  vpc_id              = ecloud_vpc.vpc-1.id
  name                = "example-gateway"
  availability_zone_id = data.ecloud_availability_zone.az-man4.id
  gateway_spec_id     = data.ecloud_backup_gateway_spec.spec-1.id
}
```

## Argument Reference

- `vpc_id`: (Required) ID of VPC where the backup gateway will be created
- `name`: Name of backup gateway
- `availability_zone_id`: (Required) ID of the availability zone where the backup gateway will be created
- `gateway_spec_id`: (Required) ID of the backup gateway specification that determines the gateway's capabilities

## Attributes Reference

- `id`: ID of backup gateway
- `vpc_id`: ID of VPC
- `name`: Name of backup gateway
- `availability_zone_id`: ID of availability zone
- `gateway_spec_id`: ID of backup gateway specification