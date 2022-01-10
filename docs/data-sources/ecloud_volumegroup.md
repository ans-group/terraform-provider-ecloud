# ecloud_volumegroup Data Source

This resource represents an eCloud Volumegroup

## Example Usage

```hcl
data "ecloud_volumegroup" "vg-1" {
  name = "my-volumegroup"
}
```

## Argument Reference

- `volume_group_id`: ID of volumegroup
- `vpc_id`: ID of volumegroup VPC
- `availability_zone_id`: ID of volumegroup availability zone
- `name`: Name of volumegroup

## Attributes Reference

`id` is set to volumegroup ID

- `vpc_id`: ID of volumegroup VPC
- `name`: Name of volumegroup
- `availability_zone_id`: ID of volumegroup availability zone
