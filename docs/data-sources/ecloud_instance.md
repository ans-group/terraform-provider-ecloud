# ecloud_instance Data Source

This resource represents an eCloud Network

## Example Usage

```hcl
data "ecloud_instance" "instance-1" {
    name = "testinstance"
}
```

## Argument Reference

- `name`: (Required) Name of instance