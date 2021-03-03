# ecloud_firewallrule Resource

This resource is for managing eCloud Firewall Rules

## Example Usage

```hcl
resource "ecloud_firewallrule" "firewallrule-web" {
  firewall_policy_id = "fwp-abcdef12"
  sequence           = 0
  name               = "firewallrule-web"
  direction          = "IN"
  source             = "ANY"
  destination        = "ANY"
  action             = "ALLOW"
  enabled            = true

  port {
    protocol    = "TCP"
    source      = "ANY"
    destination = "80"
  }

  port {
    protocol    = "TCP"
    source      = "ANY"
    destination = "443"
  }
}
```

## Argument Reference

- `firewall_policy_id`: (Required) ID of firewall policy for rule
- `sequence`: (Required) Sequence / ordering of firewall rule
- `name`: Name of firewall rule
- `direction`: (Required) Direction of firewall rule. One of: `IN`, `OUT`, `IN_OUT`
- `action`: (Required) Action of firewall rule. One of: `ALLOW`, `DROP`, `REJECT`
- `source`: (Required) Source of firewall rule. Accepts IP range / CIDR or `ANY`. Examples: `192.168.1.1`, `192.168.1.0/24`, `192.168.1.0-192.168.1.100`, `ANY`
- `destination`: (Required) Destination of firewall rule. Accepts IP range / CIDR or `ANY`. Examples: `192.168.1.1`, `192.168.1.0/24`, `192.168.1.0-192.168.1.100`, `ANY`
- `enabled`: Specifies whether firewall rule is enabled
- `port`: Map of ports for rule
  - `protocol`: (Required) Protocol of port/service. One of: `TCP`, `UDP`, `ICMPv4`
  - `source`: (Required if `protocol` is `TCP` or `UDP`)
  - `destination`: (Required if `protocol` is `TCP` or `UDP`)