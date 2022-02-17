# ecloud_loadbalancer_network Data Source

This resource represents an eCloud LoadBalancer Network. 
This data source can be used to retrieve the ID of a particular loadbalancer network.

## Example Usage

```hcl
data "ecloud_loadbalancer_network" "lbs-network" {
  name = "my-lb-network"
}
```

## Argument Reference

- `load_balancer_network_id`: ID of loadbalancer network
- `name`: Name of loadbalancer network
- `loadbalancer_id`: ID of loadbalancer
- `network_id`: ID of ecloud network

## Attributes Reference

`id` is set to loadbalancer network ID

- `name`: Name of loadbalancer network
- `network_id`: ID of ecloud network
- `loadbalancer_id`: ID of loadbalancer