# ecloud_network Data Source

This resource represents an eCloud network

## Example Usage

```hcl
data "ecloud_network" "network-1" {
    name = "some network"
    solution_id = 12345
}
```

## Argument Reference

* `name`: (Required) Name of network
* `solution_id`: (Required) ID of solution which the network is a member of