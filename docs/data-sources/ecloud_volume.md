# ecloud_volume Data Source

This resource represents an eCloud Volume

## Example Usage

```hcl
data "ecloud_volume" "volume-1" {
  name = "my-volume"
}
```

## Argument Reference

- `volume_id`: ID of volume
- `vpc_id`: ID of volume VPC
- `name`: Name of volume

## Attributes Reference

`id` is set to volume ID

- `vpc_id`: ID of volume VPC
- `name`: Name of volume