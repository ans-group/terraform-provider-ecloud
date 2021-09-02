# ecloud_vpn_service Data Source

This resource represents an eCloud VPN service

## Example Usage

```hcl
data "ecloud_vpn_service" "service-1" {
  name = "some-service"
}
```

## Argument Reference

- `vpn_service_id`: ID of VPN service
- `name`: Name of VPN service
- `vpc_id`: ID of VPC
- `router_id`: ID of router

## Attributes Reference

`id` is set to VPN service ID

- `name`: Name of VPN service
- `vpc_id`: ID of VPC
- `router_id`: ID of router