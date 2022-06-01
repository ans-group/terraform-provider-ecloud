# ecloud_vpn_session Data Source

This resource represents an eCloud VPN session

## Example Usage

```hcl
data "ecloud_vpn_session" "session-1" {
  name = "some-session"
}
```

## Argument Reference

- `vpn_session_id`: ID of VPN session
- `vpn_service_id`: ID of VPN service
- `vpn_profile_group_id`: ID of profile group
- `vpn_endpoint_id`: ID of VPN endpoint
- `remote_ip`: IP address of remote
- `name`: Name of VPN session

## Attributes Reference

`id` is set to VPN session ID

- `vpn_service_id`: ID of VPN service
- `vpn_profile_group_id`: ID of profile group
- `vpn_endpoint_id`: ID of VPN endpoint
- `remote_ip`: IP address of remote
- `name`: Name of VPN session
- `remote_networks`: Comma seperated list of remote network CIDRs
- `local_networks`: Comma seperated list of local network CIDRs
- `psk`: Pre-shared key for VPN session