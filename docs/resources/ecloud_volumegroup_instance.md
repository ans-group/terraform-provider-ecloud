# ecloud_volumegroup_instance Resource

This resource is for managing eCloud Volume group instance attachments

## Example Usage

```hcl
resource "ecloud_volumegroup_instance" "vgi-1" {
  volume_group_id = "vg-abcdef12"
  instance_id     = "i-abcdef12"
}
```

## Argument Reference

- `volume_group_id`: (Required) ID of volume group
- `instance_id`: ID of the instance to attach to volume group

## Attribute Reference

- `volume_group_id`: ID of volume group
- `instance_id`: ID of the instance to attach to volume group
