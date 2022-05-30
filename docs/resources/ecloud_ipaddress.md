# ecloud_ipaddress Resource

This resource is for managing eCloud IP addresses

## Example Usage

```hcl
resource "ecloud_ipaddress" "ipaddress-1" {
  network_id          = "net-abcdef12"
  name                = "my-ip"
}
```

## Argument Reference

- `network_id`: (Required) ID of network
- `name`: Name of IP address
- `ip_address`: IP address to assign
