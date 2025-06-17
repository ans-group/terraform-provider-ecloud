# ecloud_monitoring_gateway Data Source

This data source represents an eCloud Monitoring Gateway

## Example Usage

```hcl
data "ecloud_monitoring_gateway" "gateway-1" {
  name = "my-monitoring-gateway"
}
```

## Argument Reference

- `monitoring_gateway_id`: ID of monitoring gateway
- `vpc_id`: ID of VPC
- `name`: Name of monitoring gateway
- `router_id`: ID of monitoring gateway router

## Attributes Reference

`id` is set to the monitoring gateway ID

- `vpc_id`: ID of VPC
- `name`: Name of monitoring gateway
- `router_id`: ID of monitoring gateway router
- `specification_id`: ID of monitoring gateway specification
