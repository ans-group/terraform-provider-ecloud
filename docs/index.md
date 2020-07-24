# eCloud Provider

Official UKFast eCloud Terraform provider, allowing for manipulation of eCloud environments

## Example Usage

```hcl
provider "ecloud" {
  api_key = "abc"
}

resource "ecloud_virtualmachine" "vm-1" {
    cpu = 2
    ram = 2
    disk {
      capacity = 20
    }
    template = "CentOS 7 64-bit"
    name = "vm-1"
    environment = "Hybrid"
    solution_id = 123
}
```

## Argument Reference

* `api_key`: UKFast API key - read/write permissions for `ecloud` service required. If omitted, will use `UKF_API_KEY` environment variable value