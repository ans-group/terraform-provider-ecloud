# ecloud_image Resource

This resource is for managing eCloud images

## Example Usage

```hcl
resource "ecloud_image" "image-1" {
  instance_id = "i-abcdef12"
  name   = "my-image"
}
```

## Argument Reference

- `instance_id`: (Required) ID of source instance
- `name`: Name of image

## Attribute Reference

- `id`: ID of the image
- `vpc_id`: ID of image VPC
- `name`: Name of image
- `availability_zone_id`: ID of image availability zone