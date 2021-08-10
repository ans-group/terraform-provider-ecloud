# ecloud_availability_zone Data Source

This resource represents an eCloud Availability Zone. This data source can be used to retrieve the resource ID of a particular AZ by name.

## Example Usage

```hcl
data "ecloud_availability_zone" "man-az" {
    name = "Manchester West"
}
```

## Argument Reference

- `availability_zone_id`: ID of availability zone
- `name`: Name of availability zone
- `region_id`: Name of availability zone region
- `code`: Availability zone code 


## Attributes Reference

`id` is set to Availability Zone ID

- `name`: Name of availability zone
- `region_id`: Name of availability zone region
- `code`: Availability zone code 