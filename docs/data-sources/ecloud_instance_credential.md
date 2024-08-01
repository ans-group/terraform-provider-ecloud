# ecloud_credential Data Source

This resource represents an eCloud 

## Example Usage

```hcl
data "ecloud_credential" "instance1_root" {
  instance_id = "i-abcdef12"
  username = "root"
}
```

## Argument Reference

- `instance_id`: (required) ID of the instance
- `username`: Username of credential
- `name`:   Name of credential

## Attributes Reference

`id` is set to host ID

- `name`: Name of host
- `username`: credential username
- `password`: Credential password