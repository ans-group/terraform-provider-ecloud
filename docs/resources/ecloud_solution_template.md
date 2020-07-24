# ecloud_solution_template Resource

This resource is for managing eCloud solution templates

## Example Usage

```hcl
resource "ecloud_solution_template" "solution-template-1" {
    solution_id = 12345
    virtualmachine_id = 12345
    name = "some template"
}
```

## Argument Reference

* `solution_id`: (Required) ID of target solution
* `virtualmachine_id`: (Required) ID of source virtual machine from which template will be created
* `name`: (Required) Name of template