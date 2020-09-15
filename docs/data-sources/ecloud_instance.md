# ecloud_instance Data Source

This resource represents an eCloud Network

## Example Usage

```hcl
data "ecloud_instance" "instance-1" {
    name = "testinstance"
}
```

## Argument Reference

- `instance_id`: ID of instance
- `name`: Name of instance

## Attributes Reference

`id` is set to instance ID

- `name`: Name of instance