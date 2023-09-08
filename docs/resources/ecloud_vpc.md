# ecloud_vpc Resource

This resource is for managing eCloud VPCs

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
  region_id           = "reg-abcdef12"
  name                = "my-vpc"
  advanced_networking = false
}
```

## Argument Reference

- `region_id`: (Required) ID of VPC region
- `name`: Name of VPC
- `client_id`: ID of VPC client
- `advanced_networking`: Whether advanced networking is enabled or disabled for the VPC. Can only be set during VPC creation. When enabled, network policies and rules can be applied to restrict East-West traffic flow between networks.
