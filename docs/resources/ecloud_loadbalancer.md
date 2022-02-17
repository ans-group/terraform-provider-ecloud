# ecloud_loadbalancer Resource

This resource is for managing eCloud LoadBalancers. 

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
  region_id = "reg-abcdef12"
  name = "example-vpc"
}

data "ecloud_loadbalancer_spec" "lb-medium" {
	name = "Medium"
}

data "ecloud_availability_zone" "az-man4" {
  name = "Manchester West"
}

resource "ecloud_loadbalancer" "app-lb" {
   vpc_id = ecloud_vpc.vpc-1.id
   availability_zone_id = data.ecloud_availability_zone.az-man4.id
   name = "app lb"
   load_balancer_spec_id = data.ecloud_loadbalancer_spec.lb-medium.id
}
```

## Argument Reference

- `vpc_id`: (Required) ID of VPC
- `load_balancer_spec_id`: (Required) ID of the load balancer spec, which determines the spec of loadbalancer to be created.
- `availability_zone_id`: (Required) ID of the availability zone where the load balancer will be created.
- `name`: Name of loadbalancer

## Attributes Reference

- `id`: ID of loadbalancer
- `vpc_id`: ID of VPC
- `name`: Name of load balancer
- `load_balancer_spec_id`: ID of the loadbalancer spec used by the loadbalancer
- `config_id`: Configuration ID of the loadbalancer