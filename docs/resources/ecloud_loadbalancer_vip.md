# ecloud_loadbalancer_vip Resource

This resource is for managing eCloud LoadBalancer VIPs. 

An eCloud LoadBalancer VIP resource must reference the ID of a LoadBalancer resource.

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
    region_id = "reg-abcdef12"
    name      = "example-vpc"
}

data "ecloud_loadbalancer_spec" "medium-lb" {
    name = "Medium
}

data "ecloud_availability_zone" "az-man4" {
    name = "Manchester South"
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
    name                  = "lb-1"
    load_balancer_spec_id = data.ecloud_loadbalancer_spec.medium-lb.id
    network_id            = ecloud_network.network-1.id
}

resource "ecloud_loadbalancer_vip" "lb-vip" {
    name             = "example loadbalancer"
    load_balancer_id = ecloud_loadbalancer.lb-1.id
}
```

## Argument Reference

- `load_balancer_id`: (Required) ID of the LoadBalancer resource with which to associate the VIP. 
- `name`: Name of LoadBalancer.
- `allocate_floating_ip`: Whether to allocate a floating IP to the LoadBalancer VIP on creation. (false if undefined)

## Attributes Reference

- `id`: ID of loadbalancer VIP
- `name`: Name of LoadBalancer
- `load_balancer_id`: Id of the LoadBalancer resource
- `floating_ip_id`: Id of the floating IP allocated to the VIP, if it exists
