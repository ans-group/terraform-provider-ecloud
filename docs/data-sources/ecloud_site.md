# ecloud_site Data Source

This resource represents an eCloud site

## Example Usage

```hcl
data "ecloud_site" "site-1" {
    pod_id = 12345
    solution_id = 12345
}
```

## Argument Reference

* `pod_id`: (Required) ID of Pod which site is a member of
* `solution_id`: (Required) ID of solution which the site is a member of
* `state`: State of the site