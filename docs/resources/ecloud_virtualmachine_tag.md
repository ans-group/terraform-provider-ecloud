# ecloud_virtualmachine_tag Resource

This resource is for managing eCloud virtual machine tags

## Example Usage

```hcl
resource "ecloud_virtualmachine_tag" "vm-tag-1" {
    virtualmachine_id = 12345
    key = "some-key"
    value = "some-value"
}
```

## Argument Reference

* `virtualmachine_id`: (Required) ID of target virtual machine
* `key`: (Required) Key for tag
* `value`: Value for tag