# ecloud_affinityrule_member Resource

This resource is for managing an eCloud Affinity Rule Member

## Example Usage

```hcl
resource "ecloud_affinityrule_member" "rule-member-1" {
  instance_id        = "i-abcdef12"
  affinity_rule_id   = "ar-abcdef12"
}
```

## Argument Reference

- `instance_id`: (Required) ID of instance
- `affinity_rule_id`: (Required) ID of the associated affinity rule.
