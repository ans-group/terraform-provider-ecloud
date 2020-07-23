# ecloud_solution_tag Resource

This resource is for managing eCloud solution tags

## Example Usage

```hcl
resource "ecloud_solution_tag" "solution-tag-1" {
    solution_id = 12345
    key = "some-key"
    value = "some-value"
}
```

## Argument Reference

* `solution_id`: (Required) ID of target solution
* `key`: (Required) Key for tag
* `value`: Value for tag