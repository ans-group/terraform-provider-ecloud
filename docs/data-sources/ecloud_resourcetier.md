# ecloud_resourcetier Data Source

This resource represents an eCloud Resource Tier

## Example Usage

```hcl

data "ecloud_availability_zone" "man-south" {
    name = "Manchester South"
}

data "ecloud_resourcetier" "rt-standard" {
  name = "Standard CPU"
  availability_zone_id = data.ecloud_availability_zone.man-south.id
}
```

## Argument Reference

- `resource_tier_id`: ID of resource tier
- `name`: Name of image
- `availability_zone_id` : ID of availability zone

## Attributes Reference

`id` is set to resource tier ID

- `name`: Name of resource tier
- `availability_zone_id`: ID of availability zone