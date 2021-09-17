# ecloud_networkrule Data Source

This resource represents an eCloud Network Rule

## Example Usage

```hcl
data "ecloud_networkrule" "networkrule-1" {
  name = "my-rule"
}
```

## Argument Reference

- `network_rule_id`: ID of network rule
- `network_policy_id`: ID of network policy for rule
- `sequence`: Sequence / ordering of network rule
- `name`: Name of network rule
- `source`: Source of network rule
- `destination`: Destination of network rule
- `action`: Action of network rule
- `direction`: Direction of network rule
- `enabled`: Specifies whether network rule is enabled

## Attributes Reference

`id` is set to network rule ID

- `network_policy_id`: ID of network policy for rule
- `sequence`: Sequence / ordering of network rule
- `name`: Name of network rule
- `source`: Source of network rule
- `destination`: Destination of network rule
- `action`: Action of network rule
- `direction`: Direction of network rule
- `enabled`: Specifies whether network rule is enabled