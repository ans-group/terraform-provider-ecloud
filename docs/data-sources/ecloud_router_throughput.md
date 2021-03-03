# ecloud_router_throughput Data Source

This resource represents an eCloud Router throughput

## Example Usage

```hcl
data "ecloud_router_throughput" "throughput-1" {
  name = "some-throughput"
}
```

## Argument Reference

- `router_throughput_id`: ID of router throughput
- `availability_zone_id`: ID of availability zone
- `name`: Name of router throughput

## Attributes Reference

`id` is set to router throughput ID

- `availability_zone_id`: ID of availability zone
- `name`: Name of router throughput
- `committed_bandwidth`: Committed bandwidth configured for throughput
- `burst_size`: Burst size configured for throughput