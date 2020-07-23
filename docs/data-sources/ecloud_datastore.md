# ecloud_datastore Data Source

This resource represents an eCloud Datastore

## Example Usage

```hcl
data "ecloud_datastore" "datastore-1" {
    name = "datastore-1"
    solution_id = 12345
}
```

## Argument Reference

* `name`: (Required) Name of datastore
* `solution_id`: (Required) ID of solution which the datastore is a member of
* `site_id`: ID of site which the datastore is a member of
* `status`: Status of datastore
* `capacity`: Capacity of datastore