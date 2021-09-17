# ecloud_firewallpolicy Resource

This resource is for managing eCloud Firewall Policies

## Example Usage

```hcl
resource "ecloud_firewallpolicy" "firewallpolicy-1" {
  router_id = "rtr-abcdef12"
  sequence  = 0
  name      = "my-firewallpolicy"
}
```

## Argument Reference

- `router_id`: (Required) ID of firewall policy router
- `sequence`: (Required) Sequence / ordering of firewall policy
- `name`: Name of firewall policy