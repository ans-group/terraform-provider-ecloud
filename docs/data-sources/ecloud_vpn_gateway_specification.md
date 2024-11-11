# ecloud_vpn_gateway_specification Data Source

This resource represents an eCloud VPN gateway specification

## Example Usage

```hcl
data "ecloud_vpn_gateway_specification" "spec-1" {
    name = "Small"
}
```

## Argument Reference

- `vpn_gateway_specification_id`: ID of VPN gateway specification
- `name`: Name of VPN gateway specification

## Attributes Reference

`id` is set to VPN gateway specification ID

- `name`: Name of VPN gateway specification
- `description`: Description of VPN gateway specification