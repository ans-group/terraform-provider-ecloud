# terraform-provider-ecloud

## Getting Started

To get started, the `terraform-provider-ecloud` binary (`.exe` extension if Windows) should be downloaded from [Releases](https://github.com/ukfast/terraform-provider-ecloud/releases) and placed in a directory. For this example,
we'll place it at `/tmp/terraform-provider-ecloud`.

Next, we'll go ahead and create a new directory to hold our `terraform` file and state:

```console
mkdir /home/user/terraform
```

We'll then create an example terraform file `/home/user/terraform/test.tf`:

```console
cat <<EOF > /home/user/terraform/test.tf
provider "ecloud" {
  api_key = "abc"
}

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
EOF
```

We'll then need to initialise terraform with our provider (specifying `plugin-dir` as the path to where the provider was downloaded to earlier):

```console
terraform init -get-plugins=false -plugin-dir=/tmp/terraform-provider-ecloud
```

Finally, we can invoke `terraform apply` to apply our terraform configuration:

```console
terraform apply
```

## Provider

**Parameters**

- `api_key`: UKFast API key - read/write permissions for `ecloud` service required. If omitted, will use `UKF_API_KEY` environment variable value

## Resources

### ecloud_virtualmachine

**Schema**

- `cpu`: (Required) CPU count
- `ram`: (Required) Amount of RAM in Gibibytes
- `disk`: (Required) Disk(s) to attach
  - `capacity` (Required) Capacity of disk
  - `name` Name of disk to target
- `template`: (Required) Template/OS name
- `template_password`: Password for template (if using custom template)
- `name`: Name of VM
- `computername`: Computer name for VM
- `environment`: Environment for VM
- `solution_id`: ID of solution which the VM is a member of
- `datastore_id`: ID of datastore on which to launch VM
- `site_id`: ID of site on which to launch VM
- `network_id`: ID of network on which to launch VM
- `external_ip_required`: Specifies that an external IP address should be allocated
- `power_status`: Power status of VM. Valid values: `Online` (default), `Offline`
- `ssh_keys`: An array of SSH public keys to apply to VM

### ecloud_virtualmachine_tag

**Schema**

- `virtualmachine_id`: (Required) ID of target virtual machine
- `key`: (Required) Key for tag
- `value`: Value for tag

### ecloud_solution_tag

**Schema**

- `solution_id`: (Required) ID of target solution
- `key`: (Required) Key for tag
- `value`: Value for tag

### ecloud_solution_template

**Schema**

- `solution_id`: (Required) ID of target solution
- `virtualmachine_id`: (Required) ID of source virtual machine from which template will be created
- `name`: (Required) Name of template

## Data sources

### ecloud_datastore

**Schema**

- `name`: (Required) Name of datastore
- `solution_id`: (Required) ID of solution which the datastore is a member of
- `site_id`: ID of site which the datastore is a member of
- `status`: Status of datastore
- `capacity`: Capacity of datastore

### ecloud_site

**Schema**

- `pod_id`: (Required) ID of Pod which site is a member of
- `solution_id`: (Required) ID of solution which the site is a member of
- `state`: State of the site

### ecloud_solution

**Schema**

- `name`: (Required) Name of solution
- `environment`: Environment for solution

### ecloud_network

**Schema**

- `name`: (Required) Name of network
- `solution_id`: (Required) ID of solution which the network is a member of
