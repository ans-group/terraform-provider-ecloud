# ecloud_hostgroup Data Source

This resource represents an eCloud Host Group. 

## Example Usage

```hcl
data "ecloud_hostgroup" "hg-1" {
  name = "example-hostgroup"
}
```

## Argument Reference

- `host_group_id`: ID of host group
- `vpc_id`: ID of VPC
- `name`: Name of host group

## Attributes Reference

`id` is set to host group ID

- `vpc_id`: ID of VPC
- `name`: Name of host group
- `host_spec_id`: ID of the host spec used by the host group