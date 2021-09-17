# ecloud_vpn_endpoint Resource

This resource is for managing VPN endpoints

## Example Usage

```hcl
resource "ecloud_vpn_endpoint" "endpoint-1" {
  vpn_service_id = "vpn-abcdef12"
  name           = "endpoint-1"
}
```

## Argument Reference

- `name`: Name of VPN endpoint
- `vpn_service_id`: ID of VPN service
- `floating_ip_id`: Floating IP ID to assign

## Attributes Reference

- `id`: ID of VPN endpoint
- `name`: Name of VPN endpoint
- `vpn_service_id`: ID of VPN service
- `floating_ip_id`: Floating IP ID assigned