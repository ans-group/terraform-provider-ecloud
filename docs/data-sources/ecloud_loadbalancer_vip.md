# ecloud_loadbalancer_vip Data Source

This resource represents an eCloud LoadBalancer VIP. 
This data source can be used to retrieve the ID of a particular loadbalancer VIP.

## Example Usage

```hcl
data "ecloud_loadbalancer_vip" "lb-vip" {
  name = "my-lb-vip"
}
```

## Argument Reference

- `load_balancer_vip_id`: ID of loadbalancer vip
- `name`: Name of loadbalancer vip
- `loadbalancer_id`: ID of loadbalancer

## Attributes Reference

`id` is set to loadbalancer vip ID

- `name`: Name of loadbalancer vip
- `loadbalancer_id`: ID of loadbalancer