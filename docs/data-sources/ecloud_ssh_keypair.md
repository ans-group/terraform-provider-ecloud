# ecloud_ssh_keypair Data Source

This resource represents an eCloud SSH Key Pair. 

## Example Usage

```hcl
data "ecloud_ssh_keypair" "keypair-1" {
  name = "my-public-key"
}
```

## Argument Reference

- `ssh_keypair_id`: ID of SSH key pair
- `name`: Name of SSH key pair 

## Attributes Reference

`id` is set to SSH key pair ID

- `name`: Name of SSH key pair
- `public_key`: The public key string