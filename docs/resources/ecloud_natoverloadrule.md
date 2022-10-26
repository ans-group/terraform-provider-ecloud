# ecloud_natoverloadrule Resource

This resource is for managing eCloud NAT overload rules

## Example Usage

```hcl
resource "ecloud_natoverloadrule" "rule-1" {
  network_id = "net-abcdef12"
  floating_ip_id = "fip-abcdef12"
  subnet    = "10.0.0.0/24"
  action    = "allow"
}
```

## Argument Reference

- `network_id`: (Required) ID of rule network
- `floating_ip_id`: (Required) ID of floating IP for rule
- `subnet`: (Required) Subnet for rule
- `action`: (Required) Action for rule (`allow`/`deny`)
- `name`: Name of rule