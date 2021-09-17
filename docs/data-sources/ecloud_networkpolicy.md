# ecloud_networkpolicy Data Source

This resource represents an eCloud Network Policy

## Example Usage

```hcl
data "ecloud_networkpolicy" "networkpolicy-1" {
  name = "my-policy"
}
```

## Argument Reference

- `network_policy_id`: ID of network policy
- `network_id`: ID of network policy network
- `vpc_id`: ID of network policy VPC
- `name`: Name of network policy

## Attributes Reference

`id` is set to network policy ID

- `network_id`: ID of network policy network
- `vpc_id`: ID of VPC
- `name`: Name of network policy