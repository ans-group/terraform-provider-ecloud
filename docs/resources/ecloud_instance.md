# ecloud_network Resource

This resource is for managing eCloud Networks

## Example Usage

```hcl
resource "ecloud_network" "network-1" {
    router_id = "rtr-abcdef12"
    subnet    = "10.0.0.0/24"
}
```

## Argument Reference

- `vpc_id`: (Required) ID of instance VPC
- `name`: Name of instance
- `appliance_id`: (Required) ID of appliance
- `user_script`: Script to execute in-guest
- `vcpu_cores`: (Required) Count of vCPU cores for instance
- `ram_capacity`: (Required) Amount of RAM/Memory for instance
- `volume_capacity`: (Required) Size of volume to allocate for instance
- `locked`: Specifies instance should be locked from update/delete
- `backup_enabled`: Specifies backup should be enabled
- `network_id`: (Required) ID of network to attach instance NIC to
- `floating_ip_id`: ID of floating IP address to assign to instance NIC
- `requires_floating_ip`: Specifies floating IP should be allocated and/or assigned