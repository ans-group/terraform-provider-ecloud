# ecloud_nic Data Source

This resource represents an eCloud NIC. This data source can be used to retrieve the internal IP address assigned to a particular NIC. 
## Example Usage

```hcl
data "ecloud_nic" "nic-1" {
  instance_id = "i-abcdef12"
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
