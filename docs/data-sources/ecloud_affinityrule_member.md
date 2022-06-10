# ecloud_affinityrule_member Data Resource

This resource represents an eCloud Affinity Rule Member

## Example Usage

```hcl
data "ecloud_affinityrule_member" "arm-1" {
  instance_id        = "i-abcdef12"
}
```

## Argument Reference

- `affinity_rule_member_id`: ID of Affinity Rule Member
- `instance_id`: ID of Instance member
- `affinity_rule_id`:  ID of affinity rule


## Attributes Reference

`id` is set to Affinity Rule Member ID

- `instance_id`: ID of Instance member
- `affinity_rule_id`:  ID of affinity rule
