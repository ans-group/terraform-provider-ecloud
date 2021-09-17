# ecloud_firewallrule Data Source

This resource represents an eCloud Firewall Rule

## Example Usage

```hcl
data "ecloud_firewallrule" "firewallrule-1" {
  name = "my-rule"
}
```

## Argument Reference

- `firewall_rule_id`: ID of firewall rule
- `firewall_policy_id`: ID of firewall policy for rule
- `sequence`: Sequence / ordering of firewall rule
- `name`: Name of firewall rule
- `source`: Source of firewall rule
- `destination`: Destination of firewall rule
- `action`: Action of firewall rule
- `direction`: Direction of firewall rule
- `enabled`: Specifies whether firewall rule is enabled

## Attributes Reference

`id` is set to firewall rule ID

- `firewall_policy_id`: ID of firewall policy for rule
- `sequence`: Sequence / ordering of firewall rule
- `name`: Name of firewall rule
- `source`: Source of firewall rule
- `destination`: Destination of firewall rule
- `action`: Action of firewall rule
- `direction`: Direction of firewall rule
- `enabled`: Specifies whether firewall rule is enabled