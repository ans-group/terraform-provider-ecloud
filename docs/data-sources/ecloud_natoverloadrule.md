# ecloud_natoverloadrule Data Source

This resource represents an eCloud Network

## Example Usage

```hcl
data "ecloud_natoverloadrule" "rule-1" {
  name = "my-rule"
}
```

## Argument Reference

- `nat_overload_rule_id`: ID of NAT overload rule
- `network_id`: ID of rule network
- `floating_ip_id`: ID of rule floating IP
- `name`: Name of rule

## Attributes Reference

`id` is set to rule ID

- `network_id`: ID of rule network
- `floating_ip_id`: ID of rule floating IP
- `subnet`: Subnet for rule
- `action`: Action for rule
- `name`: Name of rule