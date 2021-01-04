# eCloud Provider

Official UKFast eCloud Terraform provider, allowing for manipulation of eCloud environments

## Example Usage

```hcl
provider "ecloud" {
  api_key = "abc"
}

resource "ecloud_vpc" "vpc-1" {
    region_id = "reg-abcdef12"
}
```

## Argument Reference

* `api_key`: UKFast API key - read/write permissions for `ecloud` service required. If omitted, will use `UKF_API_KEY` environment variable value