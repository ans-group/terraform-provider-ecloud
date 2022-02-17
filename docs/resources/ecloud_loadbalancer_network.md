# ecloud_loadbalancer_network Resource

This resource is for managing eCloud LoadBalancer Networks. 

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
  region_id = "reg-abcdef12"
  name = "example-vpc"
}

data "ecloud_loadbalancer_spec" "medium-lb" {
	name = "Medium
}

data "ecloud_availability_zone" "az-man4" {
  name = "Manchester West"
}

resource "ecloud_router" "router-1" {
	vpc_id = ecloud_vpc.vpc-1.id
	availability_zone_id = data.ecloud_availability_zone.az-man4.id
	name = "example-router"
}

resource "ecloud_network" "network-1" {
	router_id = ecloud_router.router-1.id
	name = "example-network"
	subnet = "10.0.1.0/24"
}

resource "ecloud_loadbalancer" "lb-1" {
	vpc_id = ecloud_vpc.vpc-1.id
	availability_zone_id = data.ecloud_availability_zone.az-man4.id
	name = "lb-1"
	load_balancer_spec_id = data.ecloud_loadbalancer_spec.medium-lb.id
}

resource "ecloud_loadbalancer_network" "lb-network" {
	network_id= ecloud_network.network-1.id
	name = "example loadbalancer"
	load_balancer_id = ecloud_loadbalancer.lb-1.id
}
```

## Argument Reference

- `network_id`: (Required) ID of eCloud Network
- `load_balancer_id`: (Required) ID of the loadbalancer resource with which to associate the network. 
- `name`: Name of host group

## Attributes Reference

- `id`: ID of loadbalancer network
- `network_id`: ID of eCloud network
- `name`: Name of host group
- `load_balancer_id`: Id of the loadbalancer resource