# ecloud_firewallpolicy Data Source

This resource represents an eCloud Firewall Policy

## Example Usage

```hcl
data "ecloud_firewallpolicy" "firewallpolicy-1" {
  name = "my-policy"
}
```

## Argument Reference

- `firewall_policy_id`: ID of firewall policy
- `router_id`: ID of firewall policy router
- `sequence`: Sequence / ordering of firewall policy
- `name`: Name of firewall policy

## Attributes Reference

`id` is set to firewall policy ID

- `router_id`: ID of firewall policy router
- `sequence`: Sequence / ordering of firewall policy
- `name`: Name of firewall policy