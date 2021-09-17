# ecloud_vpn_profile_group Data Source

This resource represents an eCloud VPN profile group

## Example Usage

```hcl
data "ecloud_vpn_profile_group" "group-1" {
  name = "some-group"
}
```

## Argument Reference

- `vpn_profile_group_id`: ID of VPN profile group
- `name`: Name of VPN profile group
- `availability_zone_id`: ID of availability zone

## Attributes Reference

`id` is set to VPN profile group ID

- `name`: Name of VPN profile group
- `availability_zone_id`: ID of availability zone
- `description`: Description of VPN profile group