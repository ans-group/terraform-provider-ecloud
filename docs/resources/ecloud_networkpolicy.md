# ecloud_networkpolicy Resource

This resource is for managing eCloud Network Policies. 

`advanced_networking` must be enabled on the VPC resource in order to create Network Policies

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
```

## Argument Reference

- `network_id`: (Required) ID of the network policy network
- `name`: Name of network policy
- `catchall_rule_action`: The catchall rule action. One of "REJECT", "DROP", "ALLOW". If not specified, the default is "REJECT".


## Attributes Reference

- `id`: ID of network policy
- `vpc_id`: ID of VPC
- `name`: Name of network policy
- `catchall_rule_action`: The catchall rule action
- `catchall_rule_id`: The ID of the catchall network rule