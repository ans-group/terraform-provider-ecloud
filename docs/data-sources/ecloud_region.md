# ecloud_region Data Source

This resource represents an eCloud Region. This data source can be used to retrieve the resource ID of a particular Region by name.

## Example Usage

```hcl
data "ecloud_region" "man-region" {
    name = "Manchester"
}
```

## Argument Reference

- `region_id`: ID of region
- `name`: Name of region


## Attributes Reference

`id` is set to Availability Zone ID

- `name`: Name of region