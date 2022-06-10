# ecloud_affinityrule Resource

This resource is for managing eCloud Affinity Rule

An affinity rule resource is used to create the base affinity rule.

**NOTE** Please note that there are 2 separate methods for defining which instances this `ecloud_affinityrule` resource applies to.
- You can define each instance member using the `ecloud_affinityrule_member` resource.
- You can use the `instance_ids` property on this `ecloud_affinityrule` resource for a more concise definition. 

These methods of managing the rule members **cannot** be used together.

## Example Usage

```hcl
resource "ecloud_affinityrule" "affinityrule-1" {
  vpc_id               = "vpc-abcdef12"
  name                 = "my-affinity-rule"
  availability_zone_id = "az-abcdef12"
  type                 = "anti-affinity"

  instance_ids = [
    "i-abcdef12",
    "i-abcdef23"
  ]
}
```

## Argument Reference

- `vpc_id`: (Required) ID of VPC
- `name`: Name of Affinity Rule
- `availability_zone_id`: (Required) ID of availability zone.
- `type`: (Required) Type of rule. Accepted types: ["anti-affinity", "affinity"]
- `instance_ids`: IDs of instances to associate with the affinity rule. 
