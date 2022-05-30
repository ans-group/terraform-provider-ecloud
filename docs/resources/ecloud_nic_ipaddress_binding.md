# ecloud_nic_ipaddress_binding Resource

This resource is for managing eCloud NIC IP address bindings

## Example Usage

```hcl
resource "ecloud_nic_ipaddress_binding" "nic-ipaddress-binding-1" {
  nic_id        = "nic-abcdef12"
  ip_address_id = "ip-abcdef12"
}
```

## Argument Reference

- `nic_id`: (Required) ID of NIC
- `ip_address_id`: (Required) ID of IP address
