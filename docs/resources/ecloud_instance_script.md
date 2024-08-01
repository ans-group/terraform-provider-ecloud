# ecloud_instance_script Resource

This resource is for managing eCloud Networks

## Example Usage

```hcl
resource "ecloud_instance_script" "instance1_provision" {
  instance_id = "i-abcdef12"
  username = data.ecloud_credential.instance1_root.username
  password = data.ecloud_credential.instance1_root.password
  script = "somescript"
}
```

## Argument Reference

- `instance_id`: (Required) ID of instance
- `username`: (Required) Instance user credential
- `password`: (Required) Instance password credential
- `script`: (Required) Script content