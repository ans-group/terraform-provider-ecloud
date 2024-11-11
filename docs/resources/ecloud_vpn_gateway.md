# ecloud_vpn_gateway Resource

This resource represents an eCloud VPN gateway

## Example Usage

```hcl
data "ecloud_vpn_gateway_specification" "small" {
    name = "Small"
}

resource "ecloud_vpn_gateway" "gateway-1" {
    router_id        = "rt-abcd1234"
    name             = "example-gateway"
    specification_id = data.ecloud_vpn_gateway_specification.small.id
}
```

## Argument Reference

* `router_id` - (Required) ID of router
* `name` - (Optional) Name of VPN gateway
* `specification_id` - (Required) ID of VPN gateway specification

**Note:** The `router_id` and `specification_id` cannot be changed once the gateway is created

## Attributes Reference

`id` is set to VPN gateway ID

* `name` - Name of VPN gateway
* `router_id` - ID of router
* `specification_id` - ID of VPN gateway specification
* `fqdn` - Fully Qualified Domain Name for the VPN gateway
