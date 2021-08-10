# ecloud_floatingip Resource

This resource is for managing eCloud Floating IPs

**NOTE** Be aware that there are 2 methods of assigning a floating IP to an `ecloud_instance`.
- Specifying an instance nic ID as the `resource_id` on the floating ip resource.
- Setting `requires_floating_ip` to `true` on the instance resource.

These methods of managing the floating IP cannot be used together.

## Example Usage

```hcl
data "ecloud_availability_zone" "az-man4" {
  name = "Manchester West"
}

resource "ecloud_floatingip" "fip-1" {
  vpc_id = "vpc-abcdef12"
  availability_zone_id = data.ecloud_availability_zone.az-man4.id
  name = "tf-fip-1"

  resource_id = "nic-abcdef12"
}
```

## Argument Reference

- `vpc_id`: (Required) ID of VPC the floating IP belongs to
- `availability_zone_id`: (Required) ID of Availability Zone the floating IP belongs to
- `name`: Name of floating ip
- `resource_id`: ID of eCloud resource to assign the floating IP to. Currently this supports `ecloud_nic` resource IDs. 

## Attribute Reference

- `id`: ID of the floating IP
- `vpc_id`: ID of VPC the floating IP 
- `availability_zone_id`: ID of Availability Zone of floating IP
- `name`: Name of floating ip
- `resource_id`: ID of eCloud resource to assign the floating IP to
- `ip_address`: IP Address of the resource