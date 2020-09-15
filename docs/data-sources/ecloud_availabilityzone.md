# ecloud_availabilityzone Data Source

This resource represents an eCloud Availability Zone

## Example Usage

```hcl
data "ecloud_availabilityzone" "az-1" {
    name = "Manchester West"
}
```

## Argument Reference

- `name`: (Required) Name of availability zone