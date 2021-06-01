# ecloud_volume Resource

This resource is for managing eCloud Volumes

## Example Usage

```hcl
resource "ecloud_volume" "volume-1" {
  vpc_id = "vpc-abcdef12"
  name = "tf-volume"
  capacity = 1
  iops = 300
}
```

## Argument Reference

- `vpc_id`: (Required) ID of volume VPC
- `name`: Name of volume
- `capacity`: (Required) Volume size in GiB
- `iops`: IOPS of volume
