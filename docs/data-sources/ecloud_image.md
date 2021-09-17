# ecloud_image Data Source

This resource represents an eCloud Image

## Example Usage

```hcl
data "ecloud_image" "centos7" {
  name = "CentOS 7"
}
```

## Argument Reference

- `image_id`: ID of image
- `name`: Name of image

## Attributes Reference

`id` is set to image ID

- `name`: Name of image