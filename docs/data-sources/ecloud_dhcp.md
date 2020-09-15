# ecloud_dhcp Data Source

This resource represents an eCloud DHCP server/profile

## Example Usage

```hcl
data "ecloud_dhcp" "dhcp-1" {
    name = "my-dhcp"
}
```

## Argument Reference

- `dhcp_id`: ID of DHCP server/profile
- `name`: Name of DHCP server/profile

## Attributes Reference

`id` is set to DHCP profile/server ID

- `name`: Name of DHCP server/profile