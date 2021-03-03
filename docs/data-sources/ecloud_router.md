# ecloud_router Data Source

This resource represents an eCloud Router

## Example Usage

```hcl
data "ecloud_router" "router-1" {
  name = "my-router"
}
```

## Argument Reference

- `router_id`: ID of router
- `vpc_id`: ID of router VPC
- `name`: Name of router

## Attributes Reference

`id` is set to router ID

- `vpc_id`: ID of router VPC
- `name`: Name of router