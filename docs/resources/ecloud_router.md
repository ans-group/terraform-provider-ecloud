# ecloud_router Resource

This resource is for managing eCloud Routers

## Example Usage

```hcl
resource "ecloud_router" "router-1" {
  vpc_id = "vpc-abcdef12"
  name   = "my-router"
}
```

## Argument Reference

- `vpc_id`: (Required) ID of router VPC
- `name`: Name of router