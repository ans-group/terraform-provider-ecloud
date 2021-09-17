# ecloud_volume Resource

This resource is for managing eCloud Volumes

## Example Usage

```hcl
data "ecloud_availability_zone" "az-man4" {
	name = "Manchester West"
}

resource "ecloud_volume" "volume-1" {
  vpc_id = "vpc-abcdef12"
  availability_zone_id = data.ecloud_availability_zone.az-man4.id
  name = "tf-volume"
  capacity = 1
  iops = 300
}
```

## Argument Reference

- `vpc_id`: (Required) ID of volume VPC
- `availability_zone_id`: (Required) ID of volume Availability Zone
- `name`: Name of volume
- `capacity`: (Required) Volume size in GiB
- `iops`: IOPS of volume

## Attribute Reference

- `id`: ID of volume
- `vpc_id`: ID of volume VPC
- `availability_zone_id`:  ID of volume Availability Zone
- `name`: Name of volume
- `capacity`: Volume size in GiB
- `iops`: IOPS of volume