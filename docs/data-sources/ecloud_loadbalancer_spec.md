# ecloud_loadbalancer_spec Data Source

This resource represents an eCloud LoadBalancer Spec. 
This data source can be used to retrieve the ID of a particular load balancer spec.

## Example Usage

```hcl
data "ecloud_loadbalancer_spec" "lbs-medium" {
  name = "Medium"
}
```

## Argument Reference

- `loadbalancer_spec_id`: ID of loadbalancer spec
- `name`: Name of loadbalancer spec

## Attributes Reference

`id` is set to loadbalancer spec ID

- `name`: Name of loadbalancer spec