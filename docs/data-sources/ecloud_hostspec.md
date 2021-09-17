# ecloud_hostspec Data Source

This resource represents an eCloud Host Spec. 
This data source can be used to retrieve the ID of a particular host spec.

## Example Usage

```hcl
data "ecloud_hostspec" "hs-1" {
  name = "DUAL-4208--64GB"
}
```

## Argument Reference

- `host_spec_id`: ID of host spec
- `name`: Name of host spec

## Attributes Reference

`id` is set to host spec ID

- `name`: Name of host spec
- `cpu_sockets`: Number of CPU sockets in host spec
- `cpu_cores`: Number of CPU cores in host spec
- `ram_capacity`: RAM capacity of host spec