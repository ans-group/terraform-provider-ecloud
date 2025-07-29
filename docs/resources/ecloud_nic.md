# ecloud_nic Resource

This resource is for managing eCloud NICs

## Example Usage

```hcl
resource "ecloud_nic" "nic-1" {
  instance_id     = "i-abcdef12"
  network_id      = "net-abcdef12"
}
```

## Argument Reference

- `nic_id`: ID of NIC 
- `network_id`: ID of Network
- `instance_id`: ID of the Instance


## Attributes Reference

`id` is set to NIC ID

- `network_id`: ID of Network
- `instance_id`: ID of the Instance
- `ip_address`: Internal IP address of NIC
- `mac_address`: MAC address of the NIC
- `name`: Name of the NIC
