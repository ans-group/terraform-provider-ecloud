# ecloud_vpc Data Source

This resource represents an eCloud VPC

## Example Usage

```hcl
data "ecloud_vpc" "vpc-1" {
    name = "my-vpc"
}
```

## Argument Reference

- `name`: (Required) Name of VPC