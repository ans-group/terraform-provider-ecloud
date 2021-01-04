# ecloud_vpc Data Source

This resource is for managing eCloud VPCs

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
    region_id = "reg-abcdef12"
    name = "my-vpc"
}
```

## Argument Reference

- `region_id`: ID of VPC region
- `name`: Name of VPC