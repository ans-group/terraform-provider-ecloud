# ecloud_network Data Source

This resource represents an eCloud Network

## Example Usage

```hcl
data "ecloud_network" "network-1" {
  name = "my-network"
}
```

## Argument Reference

- `network_id`: ID of network
- `router_id`: ID of network router
- `subnet`: Subnet of network
- `name`: Name of network

## Attributes Reference

`id` is set to network ID

- `router_id`: ID of network router
- `subnet`: Subnet of network
- `name`: Name of network