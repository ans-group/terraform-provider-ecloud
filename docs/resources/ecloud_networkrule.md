# ecloud_networkrule Resource

This resource is for managing eCloud Network Rules

`advanced_networking` must be enabled on the VPC resource in order to create Network Policies and Rules

## Example Usage

```hcl
resource "ecloud_vpc" "vpc-1" {
	region_id           = "reg-abcdef12"
	name                = "my-vpc"
	advanced_networking = true
}

resource "ecloud_router" "router-1" {
	vpc_id = ecloud_vpc.test-vpc.id
	name   = "test-router"
}

resource "ecloud_network" "network-1" {
	router_id = ecloud_router.router-1.id
	subnet    = "10.0.0.0/24"
}

resource "ecloud_networkpolicy" "test-np" {
    network_id           = ecloud_network.network-1.id
    name                 = "my-networkpolicy"
    catchall_rule_action = "REJECT"
}

resource "ecloud_networkrule" "networkrule-web" {
  network_policy_id = ecloud_networkpolicy.test-np.id
  sequence           = 10
  name               = "to-network-1"
  direction          = "IN"
  source             = "172.16.0.5"
  destination        = ecloud_network.network-1.subnet
  action             = "ALLOW"
  enabled            = true

  port {
    protocol    = "TCP"
    source      = "ANY"
    destination = "3306"
    name        = "allow-mysql"
  }

}
```

## Argument Reference

- `network_policy_id`: (Required) ID of network policy for rule
- `sequence`: (Required) Sequence / ordering of network rule
- `name`: Name of network rule
- `direction`: (Required) Direction of network rule. One of: `IN`, `OUT`, `IN_OUT`
- `action`: (Required) Action of network rule. One of: `ALLOW`, `DROP`, `REJECT`
- `source`: (Required) Source of network rule. Accepts IP range / CIDR or `ANY`. Examples: `192.168.1.1`, `192.168.1.0/24`, `192.168.1.0-192.168.1.100`, `ANY`
- `destination`: (Required) Destination of network rule. Accepts IP range / CIDR or `ANY`. Examples: `192.168.1.1`, `192.168.1.0/24`, `192.168.1.0-192.168.1.100`, `ANY`
- `enabled`: Specifies whether network rule is enabled
- `port`: Map of ports for rule
  - `name`:  Name of network port rule
  - `protocol`: (Required) Protocol of port/service. One of: `TCP`, `UDP`
  - `source`: (Required) Source port / port-range. Accepts `ANY` or comma-separated list of ports / port-ranges
  - `destination`: (Required) Destination port / port-range. Accepts `ANY` or comma-separated list of ports / port-ranges.

## Attribute Reference

- `id`: ID of network rule
- `network_policy_id`: ID of network policy for rule
- `sequence`:  Sequence / ordering of network rule
- `name`: Name of network rule
- `direction`: Direction of network rule. 
- `action`: Action of network rule.
- `source`: Source of network rule.
- `destination`: Destination of network rule. 
- `enabled`: Specifies whether network rule is enabled
- `port`: Map of ports for rule
  - `name`:  Name of network port rule
  - `protocol`: Protocol of port/service. 
  - `source`:  Source port / port-range. 
  - `destination`: Destination port / port-range.