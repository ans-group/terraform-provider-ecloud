# ecloud_host Data Source

This resource represents an eCloud Host. 

## Example Usage

```hcl
data "ecloud_host" "host-1" {
  name = "example-host"
}
```

## Argument Reference

- `host_id`: ID of host
- `name`: Name of host 

## Attributes Reference

`id` is set to host ID

- `name`: Name of host
- `host_group_id`: ID of the host group the host is a member of