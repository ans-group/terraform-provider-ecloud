# ecloud_floatingip Data Source

This resource represents an eCloud Floating IP. This data source can be used to retrieve the resource ID for a particular IP address. 
## Example Usage

```hcl
data "ecloud_floatingip" "fip-1" {
  name = "tf-fip-1"
}
```

## Argument Reference

- `floating_ip_id`: ID of floating IP resource
- `vpc_id`: ID of VPC
- `availability_zone_id`: ID of availability zone
- `name`: Name of floating IP resource
- `ip_address`: IP Address belonging to the floating IP.

## Attributes Reference

`id` is set to floating IP ID

- `vpc_id`: ID of VPC
- `availability_zone_id`: ID of availability zone
- `name`: Name of floating IP
- `ip_address`: IP Address belonging to the floating IP.
