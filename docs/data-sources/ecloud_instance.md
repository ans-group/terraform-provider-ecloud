# ecloud_instance Data Source

This resource represents an eCloud Instance

## Example Usage

```hcl
data "ecloud_instance" "instance-1" {
  name = "my-instance"
}
```

## Argument Reference

- `instance_id`: ID of instance
- `vpc_id`: ID of instance VPC
- `name`: Name of instance

## Attributes Reference

`id` is set to instance ID

- `vpc_id`: ID of instance VPC
- `name`: Name of instance