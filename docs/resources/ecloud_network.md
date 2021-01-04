# ecloud_network Data Source

This resource is for managing eCloud Networks

## Example Usage

```hcl
resource "ecloud_network" "network-1" {
    router_id = "rtr-abcdef12"
}
```

## Argument Reference

- `router_id`: ID of network router
- `name`: Name of network