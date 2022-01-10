# ecloud_volumegroup Resource

This resource is for managing eCloud Volumegroups

## Example Usage

```hcl
data "ecloud_availability_zone" "az-man4" {
	name = "Manchester West"
}

resource "ecloud_volumegroup" "vg-1" {
  vpc_id = "vpc-abcdef12"
  availability_zone_id = data.ecloud_availability_zone.az-man4.id
  name = "tf-volumegroup"
}
```

## Argument Reference

- `vpc_id`: (Required) ID of volumegroup VPC
- `availability_zone_id`: (Required) ID of volumegroup Availability Zone
- `name`: Name of volumegroup

## Attribute Reference

- `id`: ID of volumegroup
- `vpc_id`: ID of volumegroup VPC
- `availability_zone_id`:  ID of volumegroup Availability Zone
- `name`: Name of volumegroup
