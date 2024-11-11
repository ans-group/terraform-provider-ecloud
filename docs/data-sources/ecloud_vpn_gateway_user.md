# ecloud_vpn_gateway_user Data Source

This resource represents an eCloud VPN gateway user

## Example Usage

```hcl
data "ecloud_vpn_gateway_user" "user-1" {
    name = "example-user"
}
```

## Argument Reference

- `vpn_gateway_user_id`: ID of VPN gateway user
- `name`: Name of VPN gateway user
- `vpn_gateway_id`: ID of VPN gateway
- `username`: Username of VPN gateway user

## Attributes Reference

`id` is set to VPN gateway user ID

- `name`: Name of VPN gateway user
- `vpn_gateway_id`: ID of VPN gateway
- `username`: Username of VPN gateway user
