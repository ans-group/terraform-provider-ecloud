# ecloud_host Resource

This resource is for managing eCloud Hosts. 

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
  region_id = "reg-abcdef12"
  name = "example-vpc"
}

data "ecloud_hostspec" "hs-1" {
  name = "DUAL-E5-2620--32GB"
}

resource "ecloud_hostgroup" "hg-1" {
  vpc_id = ecloud_vpc.vpc-1.id
  host_spec_id = data.ecloud_hostspec.hs-1.id
  availability_zone_id = "az-abcd1234"
  name = "example-hostgroup"
  windows_enabled = false
}

resource "ecloud_host" "tf-host-1" {
  host_group_id = ecloud_hostgroup.hg-1.id
  name = "example-host"
# }
```

## Argument Reference

- `host_group_id`: (Required) ID of the host group to deploy host to
- `name`: Name of host

## Attributes Reference

- `id`: ID of host
- `name`: Name of host
- `host_group_id`: ID of the host group used by the host group