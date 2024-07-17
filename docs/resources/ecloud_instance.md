# ecloud_instance Resource

This resource is for managing eCloud Instances

## Example Usage

```hcl
resource "ecloud_instance" "instance-1" {
  vcpu {
    sockets = 1
    cores_per_socket = 2
  }
  ram_capacity    = 2048
  vpc_id          = "vpc-abcdef12"
  name            = "my instance"
  image_id        = "img-abcdef"
  volume_capacity = 20
  volume_iops     = 600
  network_id      = "net-abcdef12"
  backup_enabled  = false
  encrypted       = false

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
- `vcpu`: (Required) Configuration block to configure the vCPUs for this instance. The total number of vCPU cores for this instance will be `sockets`*`cores_per_socket`. The following attributes are required in this block:
  - `sockets`: (Required) The number of vCPU sockets to allocate
  - `cores_per_socket`: (Required) The number of vCPU cores per socket
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
- `volume_group_id`: ID of the volumegroup to attach to the instance. There is a separate resource for handling the attachment (`ecloud_volumegroup_instance`) which will clash with this parameter
- `host_group_id`: ID of the dedicated host group to move the instance to. Cannot be used with `resource_tier_id`
- `resource_tier_id`: ID of the public resource tier to move the instance to. Cannot be used with `host_group_id`
- `ip_address`: DHCP IP address to allocate to instance
- `encrypted`: Whether instance should be encrypted at rest
- `vcpu_cores`: (Deprecated) Count of vCPU sockets for the instance, use the new `vcpu` block, with `vcpu.sockets` and `vcpu.cores_per_socket` instead. To migrate, set `vcpu.sockets` to the value of `vcpu_cores`, and `vcpu.cores_per_socket` to `1`. Once you have migrated to the new `vcpu` configuration block, you can no longer use `vcpu_cores` for this instance.


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
- `vcpu_cores`: (Deprecated) Count of vCPU cores for instance
- `vcpu`: Block for vCPU configuration. Both of the following attributes are required:
  - `sockets`: The number of vCPU sockets for this instance
  - `cores_per_socket`: The number of vCPU cores per socket.
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
- `volume_group_id`: ID of the volumegroup attached to the instance.
- `host_group_id`: ID of the host group the instance runs on, if defined.
- `resource_tier_id`: ID of the public resource tier the instance runs on.
- `encrypted`: Whether instance is encrypted