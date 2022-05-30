# ecloud_ipaddress Data Source

This resource represents an eCloud IP address

## Example Usage

```hcl
data "ecloud_ipaddress" "ipaddress-1" {
  name = "some-ip-address"
}
```

## Argument Reference

- `ip_address_id`: ID of IP address
- `name`: Name of IP address
- `ip_address`: Assigned IP address
- `network_id`: ID of network
- `type`: Type of IP address

## Attributes Reference

`id` is set to IP address ID

- `name`: Name of IP address
- `ip_address`: Assigned IP address
- `network_id`: ID of network
- `type`: Type of IP address