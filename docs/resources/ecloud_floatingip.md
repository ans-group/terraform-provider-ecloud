# ecloud_floatingip Resource

This resource is for managing eCloud Floating IPs

**NOTE** Be aware that there are 2 methods of assigning a floating IP to an `ecloud_instance`.
- Specifying an `instance_id` on the floating ip resource
- Setting `requires_floating_ip` to `true` on the instance resource

These methods of managing the floating IP cannot be used together.

## Example Usage

```hcl
resource "ecloud_floatingip" "fip-1" {
  vpc_id = "vpc-abcdef12"
  name = "tf-fip-1"

  instance_id = "i-abcef12"
}
```

## Argument Reference

- `vpc_id`: (Required) ID of VPC the floating IP belongs to
- `name`: Name of floating ip
- `instance_id`: ID of eCloud instance to assign the floating IP to

## Attribute Reference

- `vpc_id`: ID of VPC the floating IP 
- `name`: Name of floating ip
- `ip_address`: IP Address of the resource