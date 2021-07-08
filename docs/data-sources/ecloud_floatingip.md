# ecloud_floatingip Data Source

This resource represents an eCloud Floating IP. This data source can be used to retrieve the resource ID for a particular IP address. 
## Example Usage

```hcl
data "ecloud_floatingip" "fip-1" {
  name = "tf-fip-1"
}
```

## Argument Reference

- `vpc_id`: ID of VPC
- `name`: Name of floating IP resource
- `ip_address`: IP Address belonging to the floating IP.

## Attributes Reference

`id` is set to floating IP ID

- `vpc_id`: ID of VPC
- `name`: Name of floating IP
- `ip_address`: IP Address belonging to the floating IP.
