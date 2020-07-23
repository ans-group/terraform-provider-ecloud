# ecloud_appliance Data Source

This resource represents an eCloud marketplace appliance

## Example Usage

```hcl
data "ecloud_appliance" "appliance-wordpress" {
    name = "WordPress"
}
```

## Argument Reference

* `name`: (Required) Name of appliance