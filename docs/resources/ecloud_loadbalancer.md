# ecloud_loadbalancer Resource

This resource is for managing eCloud LoadBalancers. 

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
  region_id = "reg-abcdef12"
  name      = "example-vpc"
}

data "ecloud_loadbalancer_spec" "lb-medium" {
	name = "Medium"
}

data "ecloud_availability_zone" "az-man4" {
  name = "Manchester West"
}

resource "ecloud_router" "router-1" {
  vpc_id               = ecloud_vpc.vpc-1.id
  availability_zone_id = data.ecloud_availability_zone.az-man4.id
  name                 = "example-router"
}

resource "ecloud_network" "network-1" {
  router_id = ecloud_router.router-1.id
  name      = "example-network"
  subnet    = "10.0.1.0/24"
}

resource "ecloud_loadbalancer" "lb-1" {
   vpc_id                = ecloud_vpc.vpc-1.id
   availability_zone_id  = data.ecloud_availability_zone.az-man4.id
   name                  = "app lb"
   load_balancer_spec_id = data.ecloud_loadbalancer_spec.lb-medium.id
   network_id            = ecloud_network.network-1.id
}
```

## Argument Reference

- `vpc_id`: (Required) ID of VPC
- `load_balancer_spec_id`: (Required) ID of the LoadBalancer spec, which determines the spec of LoadBalancer to be created.
- `availability_zone_id`: (Required) ID of the availability zone where the LoadBalancer will be created.
- `name`: Name of LoadBalancer
- `network_id`: ID of the network used by the LoadBalancer

## Attributes Reference

- `id`: ID of LoadBalancer
- `vpc_id`: ID of VPC
- `name`: Name of LoadBalancer
- `load_balancer_spec_id`: ID of the LoadBalancer spec used by the LoadBalancer
- `config_id`: Configuration ID of the LoadBalancer
- `network_id`: ID of the network used by the LoadBalancer