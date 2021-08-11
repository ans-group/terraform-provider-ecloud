# ecloud_hostgroup Resource

This resource is for managing eCloud Host Groups. 

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
  region_id = "reg-abcdef12"
  name = "example-vpc"
}

data "ecloud_hostspec" "hs-1" {
  name = "DUAL-E5-2620--32GB"
}

data "ecloud_availability_zone" "az-man4" {
  name = "Manchester West"
}

resource "ecloud_hostgroup" "hg-1" {
  vpc_id = ecloud_vpc.vpc-1.id
  host_spec_id = data.ecloud_hostspec.hs-1.id
  availability_zone_id = data.ecloud_availability_zone.az-man4.id
  name = "example-hostgroup"
  windows_enabled = false
}
```

## Argument Reference

- `vpc_id`: (Required) ID of VPC
- `host_spec_id`: (Required) ID of the host spec, which determines the type of hosts in the host group
- `availability_zone_id`: (Required) ID of the availability zone where the host group will be created.
- `windows_enabled`: (Required) Whether the host group will need to support instances running Windows OS. 
- `name`: Name of host group

## Attributes Reference

- `id`: ID of host group
- `vpc_id`: ID of VPC
- `name`: Name of host group
- `host_spec_id`: ID of the host spec used by the host group