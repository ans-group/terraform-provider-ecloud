# ecloud_vpn_gateway_user Resource

This resource represents an eCloud VPN gateway user

## Example Usage

```hcl
resource "ecloud_vpn_gateway_user" "user-1" {
    vpn_gateway_id = ecloud_vpn_gateway.gateway-1.id
    name           = "example-user"
    username       = "vpnuser1"
    password       = "Password123!"
}
```

## Argument Reference

* `vpn_gateway_id` - (Required) ID of VPN gateway
* `name` - (Required) Friendly name of VPN gateway user
* `username` - (Required) Username for VPN gateway user
* `password` - (Required) Password for VPN gateway user

**Note:** The `vpn_gateway_id` and `username` cannot be changed once the user is created.

## Attributes Reference

`id` is set to VPN gateway user ID

* `name` - Name of VPN gateway user
* `vpn_gateway_id` - ID of VPN gateway
* `username` - Username of VPN gateway user