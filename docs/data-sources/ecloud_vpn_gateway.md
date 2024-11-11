# ecloud_vpn_gateway Data Source

This resource represents an eCloud VPN gateway

## Example Usage

```hcl
data "ecloud_vpn_gateway" "gateway-1" {
    name = "example-gateway"
}
```

## Argument Reference

- `vpn_gateway_id`: ID of VPN gateway
- `name`: Name of VPN gateway
- `router_id`: ID of router
- `specification_id`: ID of VPN gateway specification

## Attributes Reference

`id` is set to VPN gateway ID

- `name`: Name of VPN gateway
- `router_id`: ID of router
- `specification_id`: ID of VPN gateway specification
- `fqdn`: Fully Qualified Domain Name for the VPN gateway
