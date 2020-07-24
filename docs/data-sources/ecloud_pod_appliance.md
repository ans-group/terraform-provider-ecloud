# ecloud_pod_appliance Data Source

This resource represents an eCloud marketplace appliance within a pod

## Example Usage

```hcl
data "ecloud_pod_appliance" "appliance-wordpress" {
    name = "WordPress"
    pod_id = 12345
}
```

## Argument Reference

* `name`: (Required) Name of appliance
* `pod_id`: (Required) ID of pod which the appliance exists on