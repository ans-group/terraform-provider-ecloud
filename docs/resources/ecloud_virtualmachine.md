# ecloud_virtualmachine Resource

This resource is for managing eCloud virtual machines

## Example Usage

```hcl
resource "ecloud_virtualmachine" "vm-1" {
    cpu = 2
    ram = 2
    disk {
      capacity = 20
    }
    template = "CentOS 7 64-bit"
    name = "vm-1"
    environment = "Hybrid"
    solution_id = 123
}
```

## Argument Reference

* `cpu`: (Required) CPU count
* `ram`: (Required) Amount of RAM in Gibibytes
* `disk`: (Required) Disk(s) to attach
  * `capacity` (Required) Capacity of disk
  * `name` Name of disk to target
* `template`: Template/OS name
* `template_password`: Password for template (if using custom Linux template)
* `appliance_id`: ID of Marketplace appliance to launch from (Mutually exclusive with `template`)
* `appliance_parameters`: An array of appliance parameters
  * `key`: Key of parameter
  * `value`: Value of parameter
* `name`: Name of VM
* `computername`: Computer name for VM
* `environment`: Environment for VM
* `solution_id`: ID of solution which the VM is a member of
* `datastore_id`: ID of datastore on which to launch VM
* `site_id`: ID of site on which to launch VM
* `network_id`: ID of network on which to launch VM
* `external_ip_required`: Specifies that an external IP address should be allocated
* `power_status`: Power status of VM. Valid values: `Online` (default), `Offline`
* `ssh_keys`: An array of SSH public keys to apply to VM
* `role`: Role for VM
* `bootstrap_script`: Script to be executed on first boot
* `activedirectory_domain_id`: ID of Active Directory Domain for VM (applicable to Windows-based virtual machines only)
* `pod_id`: ID of Pod on which to launch VM