# ecloud_affinityrule Data Resource

This resource represents an eCloud Affinity Rule

## Example Usage

```hcl
data "ecloud_affinityrule" "affinityrule-1" {
  type              = "anti-affinity"
  name              = "my-affinity-rule"
}
```

## Argument Reference

- `affinity_rule_id`: ID of Affinity Rule
- `vpc_id`: ID of VPC
- `name`: Name of Affinity Rule
- `availability_zone_id`:  ID of availability zone.
- `type`: Type of affinity rule. Accepted types: ["anti-affinity", "affinity"]


## Attributes Reference

`id` is set to Affinity Rule ID

- `name`: Name of affinity rule
- `availability_zone_id`: ID of availability zone
- `vpc_id`: ID of VPC
- `type`: Type of affinity rule
