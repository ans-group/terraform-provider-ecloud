# ecloud_dhcp Data Source

This resource represents an eCloud DHCP server/profile

## Example Usage

```hcl
data "ecloud_dhcp" "dhcp-1" {
    name = "my-dhcp"
}
```

## Argument Reference

- `availability_zone_id`: (Required) ID of availability zone