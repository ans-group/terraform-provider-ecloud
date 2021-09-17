# ecloud_instance Resource

This resource is for managing eCloud Instances

## Example Usage

```hcl
resource "ecloud_instance" "instance-1" {
  vpc_id          = "vpc-abcdef12"
  name            = "my instance"
  image_id        = "img-abcdef"
  vcpu_cores      = 2
  ram_capacity    = 2048
  volume_capacity = 20
  volume_iops     = 600
  network_id      = "net-abcdef12"

  data_volume_ids = [
    "vol-abc12345"
  ]

  ssh_keypair_ids = [
    "ssh-abcd1234"
  ]
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
- `volume_iops`: IOPs of the operating system volume
- `locked`: Specifies instance should be locked from update/delete
- `backup_enabled`: Specifies backup should be enabled
- `network_id`: (Required) ID of network to attach instance NIC to
- `floating_ip_id`: ID of floating IP address to assign to instance NIC
- `requires_floating_ip`: Specifies floating IP should be allocated and assigned
- `data_volume_ids`: IDs of volumes to attach to the instance
- `image_data`: Any parameters needed for deploying an image 
- `ssh_keypair_ids`: IDs of any ssh keypairs to be added to the instance during creation 



**Note on Floating IPs** 

The optional argument `requires_floating_ip`, allows a user to quickly create and assign a floating IP address to the eCloud Instance resource without having to manage the floating IP resource independently.  

In cases where the floating IP needs to be managed separately (e.g. in order to assign to other resources, or to re-use a public IP address), please instead use the `ecloud_floatingip` managed resource to create and manage the floating IP.

If `requires_floating_ip` is set to `true` for an instance resource, **do not** use any other method of attaching a floating IP to the resource. This is to prevent floating IP conflicts.


## Attribute Reference

- `id`: ID of instance
- `vpc_id`: ID of instance VPC
- `name`: Name of instance
- `image_id`: ID of image
- `user_script`: Script to execute in-guest
- `vcpu_cores`: Count of vCPU cores for instance
- `ram_capacity`: Amount of RAM/Memory (in MiB) for instance
- `volume_capacity`: Size of OS volume (in GiB) for instance.
- `volume_iops`: IOPs of the operating system volume
- `locked`: Whether instance is locked from update/delete
- `backup_enabled`: Whether backup is be enabled
- `network_id`:  ID of instance network
- `floating_ip_id`: ID of assigned floating ip address
- `data_volume_ids`: IDs of attached data volumes
- `ssh_keypair_ids`: IDs of instance ssh keypairs 
- `host_group_id`: ID of host group