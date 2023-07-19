# ecloud_iops Data Source

This resource represents an eCloud IOPS tier

## Example Usage

```hcl
data "ecloud_iops" "iops-300" {
  number = 300
}
```

## Argument Reference

- `iops_id`: ID of IOPS tier
- `availability_zone_id`: ID of Availability Zone that tier is available in
- `name`: Name of IOPS tier
- `level`: IOPS level/limit

## Attributes Reference

`id` is set to IOPS tier ID

- `name`: Name of IOPS tier
- `level`: IOPS level/limit