# ecloud_vpc Data Source

This resource represents an eCloud VPC

## Example Usage

```hcl
data "ecloud_vpc" "vpc-1" {
    name = "my-vpc"
}
```

## Argument Reference

- `vpc_id`: ID of VPC
- `name`: Name of VPC

## Attributes Reference

`id` is set to VPC ID

- `name`: Name of VPC