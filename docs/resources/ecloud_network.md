# ecloud_network Data Source

This resource is for managing eCloud Networks

## Example Usage

```hcl
resource "ecloud_network" "network-1" {
    router_id = "rtr-abcdef12"
    subnet    = "10.0.0.0/24"
}
```

## Argument Reference

- `router_id`: ID of network router
- `subnet`: Subnet of network
- `name`: Name of network