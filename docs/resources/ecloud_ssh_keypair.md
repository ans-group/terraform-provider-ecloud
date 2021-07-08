# ecloud_ssh_keypair Resource

This resource is for managing SSH Key Pairs.

## Example Usage

```hcl
resource "ecloud_ssh_keypair" "test-keypair" {
	name = "test-keypair"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDXIismybCTbE4p24LX/Aioi17UdLUrfolbwf1fKUD2a5Ps0xvZv3U19FRTo+x6yWux7kd78DpZ50CS4WkRs09QLP9K65hSZj/SJBXl+MNaz3pJ0FngZBXxTgxdJ82gcLCvY3iDBfn61PdrJTv6kLR4ZnZruj2kBND4yUZAyQKxfzrXD20UwlF1GWwE4lHuWXaEei4mGbHSeWVay0pOEf5d6uAWlsBm2JEdXkG7/LupdLh7z+RlEaTigHarlTbpcfCC82JX94IGWmiKToFr6+lX6y7QoVxd8pmEGIV/9dxPwWM/9RczSD2Oxum83ESPhVvQrBUTjE7T7fGoLlr31rQep+qgH5XdfCqkmiZ69NFDUEPIwiqpCKazli/Jdaxz6FsxlWZbmaMOW1cMAhAtxmpOxukbhB5hmJjzR3DTAEsv5euINFNxk8snY3b77JmDYX09yb+hT/fLyjBonc7I0RmFsUIV+H25yzh57iJoSuP9Qbz2RD4nIwxn5/PvKDbwElE= test-keypair"
}
```

## Argument Reference

- `name`: Name of SSH key pair 
- `public_key`: (Required) The public key string for the SSH key pair

## Attributes Reference

- `name`: Name of SSH key pair
- `public_key`: The public key string for the SSH key pair