# ecloud_loadbalancer Data Source

This resource represents an eCloud LoadBalancer. 
This data source can be used to retrieve the ID of a particular load balancer.

## Example Usage

```hcl
data "ecloud_loadbalancer" "lb-1" {
  name = "My LoadBalancer"
}
```

## Argument Reference

- `load_balancer_id`: ID of loadbalancer
- `name`: Name of loadbalancer 
- `vpc_id`: ID of eCloud VPC
- `availability_zone_id`: ID of eCloud Availabilility Zone
- `load_balancer_spec_id`: ID of eCloud LoadBalancer Spec

## Attributes Reference

`id` is set to loadbalancer ID

- `name`: Name of loadbalancer
- `vpc_id`: ID of eCloud VPC
- `availability_zone_id`: ID of eCloud Availabilility Zone
- `load_balancer_spec_id`: ID of eCloud LoadBalancer Spec