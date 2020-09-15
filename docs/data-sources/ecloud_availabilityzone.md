# ecloud_availability_zone Data Source

This resource represents an eCloud Availability Zone

## Example Usage

```hcl
data "ecloud_availability_zone" "az-1" {
    name = "Manchester West"
}
```

## Argument Reference

- `availability_zone_id`: ID of availability zone
- `name`: Name of availability zone
- `datacentre_site_id`: Datacentre site ID for availability zone

## Attributes Reference

`id` is set to availability zone ID

- `name`: Name of availability zone
- `datacentre_site_id`: Datacentre site ID for availability zone