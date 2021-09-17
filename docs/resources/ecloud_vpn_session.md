# ecloud_vpn_session Resource

This resource is for managing VPN sessions

## Example Usage

```hcl
resource "ecloud_vpn_session" "session-1" {
  vpn_service_id       = "vpn-abcdef12"
  vpn_endpoint_id      = "vpne-abcdef12"
  vpn_profile_group_id = "vpnpg-abcdef12"
  remote_ip            = "1.2.3.4"
  local_networks       = "10.0.0.0/24"
  remote_networks      = "10.0.1.0/24"
  name                 = "session-1"
}
```

## Argument Reference

- `name`: Name of VPN session
- `vpn_service_id`: ID of VPN service
- `vpn_profile_group_id`: ID of profile group
- `vpn_endpoint_id`: ID of VPN endpoint
- `remote_ip`: IP address of remote
- `remote_networks`: Comma seperated list of remote network CIDRs
- `local_networks`: Comma seperated list of local network CIDRs

## Attributes Reference

- `id`: ID of VPN session
- `name`: Name of VPN session
- `vpn_service_id`: ID of VPN service
- `vpn_profile_group_id`: ID of profile group
- `vpn_endpoint_id`: ID of VPN endpoint
- `remote_ip`: IP address of remote
- `remote_networks`: Comma seperated list of remote network CIDRs
- `local_networks`: Comma seperated list of local network CIDRs