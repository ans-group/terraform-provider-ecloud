# ecloud_instance Resource

This resource is for managing eCloud Networks

## Example Usage

```hcl
resource "ecloud_instance" "instance-1" {
  vpc_id          = "vpc-abcdef12"
  name            = "my instance"
  image_id        = "img-abcdef"
  vcpu_cores      = 2
  ram_capacity    = 2048
  volume_capacity = 20
  network_id      = "net-abcdef12"
}
```

## Argument Reference

- `vpc_id`: (Required) ID of instance VPC
- `name`: Name of instance
- `image_id`: (Required) ID of image
- `user_script`: Script to execute in-guest
- `vcpu_cores`: (Required) Count of vCPU cores for instance
- `ram_capacity`: (Required) Amount of RAM/Memory (in MiB) for instance
- `volume_capacity`: (Required) Size of volume (in GiB) to allocate for instance.
- `locked`: Specifies instance should be locked from update/delete
- `backup_enabled`: Specifies backup should be enabled
- `network_id`: (Required) ID of network to attach instance NIC to
- `floating_ip_id`: ID of floating IP address to assign to instance NIC
- `requires_floating_ip`: Specifies floating IP should be allocated and/or assigned