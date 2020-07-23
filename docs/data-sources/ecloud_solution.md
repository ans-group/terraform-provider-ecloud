# ecloud_solution Data Source

This resource represents an eCloud solution

## Example Usage

```hcl
data "ecloud_solution" "solution-1" {
    name = "some solution"
}
```

## Argument Reference

* `name`: (Required) Name of solution
* `environment`: Environment for solution