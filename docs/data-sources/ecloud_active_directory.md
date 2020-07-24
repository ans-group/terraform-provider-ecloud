# ecloud_active_directory Data Source

This resource represents an eCloud Active Directory domain

## Example Usage

```hcl
data "ecloud_active_directory" "domain-1" {
    name = "somedomain.testing"
}
```

## Argument Reference

- `name`: (Required) Name of Active Directory domain