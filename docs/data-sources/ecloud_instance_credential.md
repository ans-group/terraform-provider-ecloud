# ecloud_instance_credential Data Source

This resource represents an eCloud 

## Example Usage

```hcl
data "ecloud_instance_credential" "instance1_root" {
  instance_id = "i-abcdef12"
  username = "root"
}
```

## Argument Reference

- `instance_id`: (required) ID of the instance
- `username`: Username of credential
- `name`:   Name of credential

## Attributes Reference

- `name`: Name of credential
- `username`: Credential username
- `password`: Credential password