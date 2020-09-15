# ecloud_network Data Source

This resource represents an eCloud Network

## Example Usage

```hcl
data "ecloud_network" "web-network" {
    name = "web"
}
```

## Argument Reference

- `network_id`: ID of network
- `name`: Name of network
- `router_id`: Router ID of network

## Attributes Reference

`id` is set to network ID

- `name`: Name of network
- `router_id`: Router ID of network