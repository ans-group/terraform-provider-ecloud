# ecloud_router Resource

This resource is for managing eCloud Routers

## Example Usage

```hcl
data "ecloud_availability_zone" "az-man4" {
  name = "Manchester West"
}

resource "ecloud_router" "router-1" {
  vpc_id = "vpc-abcdef12"
  availability_zone_id = data.ecloud_availability_zone.az-man4.id
  name   = "my-router"
}
```

## Argument Reference

- `vpc_id`: (Required) ID of router VPC
- `name`: Name of router
- `availability_zone_id`: (Required) ID of router availability zone
- `router_throughput_id`: ID of router throughput

## Attribute Reference

- `id`: ID of the router
- `vpc_id`: ID of router VPC
- `name`: Name of router
- `availability_zone_id`: ID of router availability zone
- `router_throughput_id`: ID of router throughput