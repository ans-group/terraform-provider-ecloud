# ecloud_vpn_endpoint Data Source

This resource represents an eCloud VPN endpoint

## Example Usage

```hcl
data "ecloud_vpn_endpoint" "endpoint-1" {
  name = "some-endpoint"
}
```

## Argument Reference

- `vpn_endpoint_id`: ID of VPN endpoint
- `availability_zone_id`: ID of availability zone
- `name`: Name of VPN endpoint
- `vpn_service_id`: ID of VPN service
- `floating_ip_id`: ID of floating IP

## Attributes Reference

`id` is set to vpn endpoint ID

- `availability_zone_id`: ID of availability zone
- `name`: Name of VPN endpoint
- `vpn_service_id`: ID of VPN service
- `floating_ip_id`: ID of floating IP