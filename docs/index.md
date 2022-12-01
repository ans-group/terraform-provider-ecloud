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

* `context`: Config context to use (overrides current context)
* `api_key`: API key - read/write permissions for `ecloud` service required

## Configuration

If `api_key` is omitted from the provider config, the provider will fall back to the default configuration (file / environment). Documentation for the configuration file / environment variables can be found within the [SDK repository](https://github.com/ans-group/sdk-go#configuration-file)