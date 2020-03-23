package ecloud

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/customdiff"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/ukfast/sdk-go/pkg/ptr"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func resourceVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualMachineCreate,
		Read:   resourceVirtualMachineRead,
		Update: resourceVirtualMachineUpdate,
		Delete: resourceVirtualMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"environment": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"template_password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"appliance_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"appliance_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"cpu": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"ram": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"disk": {
				Type:     schema.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"computername": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"solution_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"datastore_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"ip_internal": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_external": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"site_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"network_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"power_status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Online",
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					_, err := ecloudservice.ParseVirtualMachinePowerStatus(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q validation failed: %s", key, err))
					}
					return
				},
			},
			"external_ip_required": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"ssh_keys": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"bootstrap_script": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"activedirectory_domain_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
			"pod_id": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
		CustomizeDiff: customdiff.All(
			customdiff.ForceNewIfChange("datastore_id", func(old, new, meta interface{}) bool {
				return old.(int) > 0
			}),
			customdiff.ForceNewIfChange("site_id", func(old, new, meta interface{}) bool {
				return old.(int) > 0
			}),
			customdiff.ForceNewIfChange("network_id", func(old, new, meta interface{}) bool {
				return old.(int) > 0
			}),
		),
	}
}

func resourceVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	createReq := ecloudservice.CreateVirtualMachineRequest{
		Environment:             d.Get("environment").(string),
		Template:                d.Get("template").(string),
		TemplatePassword:        d.Get("template_password").(string),
		ApplianceID:             d.Get("appliance_id").(string),
		Parameters:              expandCreateVirtualMachineRequestApplianceParameters(d.Get("appliance_parameters").(map[string]interface{})),
		CPU:                     d.Get("cpu").(int),
		RAM:                     d.Get("ram").(int),
		Disks:                   expandCreateVirtualMachineRequestDisks(d.Get("disk").([]interface{})),
		Name:                    d.Get("name").(string),
		ComputerName:            d.Get("computername").(string),
		SolutionID:              d.Get("solution_id").(int),
		DatastoreID:             d.Get("datastore_id").(int),
		SiteID:                  d.Get("site_id").(int),
		NetworkID:               d.Get("network_id").(int),
		Role:                    d.Get("role").(string),
		ExternalIPRequired:      d.Get("external_ip_required").(bool),
		SSHKeys:                 expandVirtualMachineSSHKeys(d.Get("ssh_keys").([]interface{})),
		BootstrapScript:         d.Get("bootstrap_script").(string),
		ActiveDirectoryDomainID: d.Get("activedirectory_domain_id").(int),
		PodID:                   d.Get("pod_id").(int),
	}

	log.Printf("Created CreateVirtualMachineRequest: %+v", createReq)

	log.Print("Creating virtual machine")
	vmID, err := service.CreateVirtualMachine(createReq)
	if err != nil {
		return fmt.Errorf("Error creating virtual machine: %s", err)
	}

	d.SetId(strconv.Itoa(vmID))

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Complete"},
		Refresh:    VirtualMachineStatusRefreshFunc(service, vmID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for virtual machine with ID [%d] to return status of [Complete]: %s", vmID, err)
	}

	s := d.Get("power_status").(string)
	if s != "" {
		powerStatus, err := ecloudservice.ParseVirtualMachinePowerStatus(s)
		if err != nil {
			return fmt.Errorf("Failed to parse power status [%s]", s)
		}

		if powerStatus == ecloudservice.VirtualMachinePowerStatusOffline {
			log.Printf("Powering off virtual machine with ID [%d]", vmID)
			err := service.PowerOffVirtualMachine(vmID)
			if err != nil {
				return fmt.Errorf("Error powering off virtual machine with ID [%d]: %s", vmID, err)
			}
		}
	}

	return resourceVirtualMachineRead(d, meta)
}

func resourceVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	vmID, _ := strconv.Atoi(d.Id())

	log.Printf("Retrieving virtual machine with ID [%d]", vmID)
	vm, err := service.GetVirtualMachine(vmID)
	if err != nil {
		switch err.(type) {
		case *ecloudservice.VirtualMachineNotFoundError:
			d.SetId("")
			return nil
		default:
			return err
		}
	}

	d.Set("environment", vm.Environment)
	d.Set("template", vm.Template)
	d.Set("cpu", vm.CPU)
	d.Set("ram", vm.RAM)
	d.Set("disk", flattenVirtualMachineDisks(d.Get("disk").([]interface{}), vm.Disks))
	d.Set("name", vm.Name)
	d.Set("computername", vm.ComputerName)
	d.Set("solution_id", vm.SolutionID)
	d.Set("ip_internal", vm.IPInternal)
	d.Set("ip_external", vm.IPExternal)
	d.Set("power_status", vm.PowerStatus)

	return nil
}

func resourceVirtualMachineUpdate(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	vmID, _ := strconv.Atoi(d.Id())

	d.Partial(true)

	if d.HasChange("power_status") {
		s := d.Get("power_status").(string)
		powerStatus, err := ecloudservice.ParseVirtualMachinePowerStatus(s)
		if err != nil {
			return fmt.Errorf("Failed to parse power status [%s]", s)
		}

		switch powerStatus {
		case ecloudservice.VirtualMachinePowerStatusOnline:
			log.Printf("Powering on virtual machine with ID [%d]", vmID)
			err := service.PowerOnVirtualMachine(vmID)
			if err != nil {
				return fmt.Errorf("Error powering on virtual machine with ID [%d]: %s", vmID, err)
			}
		case ecloudservice.VirtualMachinePowerStatusOffline:
			log.Printf("Powering off virtual machine with ID [%d]", vmID)
			err := service.PowerOffVirtualMachine(vmID)
			if err != nil {
				return fmt.Errorf("Error powering off virtual machine with ID [%d]: %s", vmID, err)
			}
		default:
			return fmt.Errorf("Unsupported power status [%s]", s)
		}

		d.SetPartial("power_status")
	}

	patchRequest := ecloudservice.PatchVirtualMachineRequest{}

	hasChange := false
	if d.HasChange("name") {
		hasChange = true
		patchRequest.Name = ptr.String(d.Get("name").(string))
	}
	if d.HasChange("cpu") {
		hasChange = true
		patchRequest.CPU = d.Get("cpu").(int)
	}
	if d.HasChange("ram") {
		hasChange = true
		patchRequest.RAM = d.Get("ram").(int)
	}

	if d.HasChange("disk") {
		hasChange = true
		patchRequest.Disks = resourceVirtualMachineUpdateDisk(d.GetChange("disk"))
	}

	if d.HasChange("role") {
		hasChange = true
		patchRequest.Role = d.Get("role").(string)
	}

	log.Printf("Created PatchVirtualMachineRequest: %+v", patchRequest)

	if hasChange {
		log.Printf("Updating virtual machine with ID [%d]", vmID)
		err := service.PatchVirtualMachine(vmID, patchRequest)
		if err != nil {
			return fmt.Errorf("Error updating virtual machine with ID [%d]: %s", vmID, err)
		}

		stateConf := &resource.StateChangeConf{
			Target:     []string{"Complete"},
			Refresh:    VirtualMachineStatusRefreshFunc(service, vmID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for virtual machine with ID [%d] to return status of [Complete]: %s", vmID, err)
		}

		d.SetPartial("name")
		d.SetPartial("cpu")
		d.SetPartial("ram")
		d.SetPartial("role")
	}

	d.Partial(false)

	return resourceVirtualMachineRead(d, meta)
}

func resourceVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	service := meta.(ecloudservice.ECloudService)

	vmID, _ := strconv.Atoi(d.Id())

	log.Printf("Removing virtual machine with ID [%d]", vmID)
	err := service.DeleteVirtualMachine(vmID)
	if err != nil {
		if _, ok := err.(*ecloudservice.VirtualMachineNotFoundError); ok {
			log.Print("Virtual machine not found")
			return nil
		}

		return fmt.Errorf("Error removing virtual machine with ID [%d]: %s", vmID, err)
	}

	stateConf := &resource.StateChangeConf{
		Target:     []string{"Deleted"},
		Refresh:    VirtualMachineStatusRefreshFunc(service, vmID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for virtual machine with ID [%d] to return status of [Deleted]: %s", vmID, err)
	}

	return nil
}

// VirtualMachineStatusRefreshFunc returns a function with StateRefreshFunc signature for use
// with StateChangeConf
func VirtualMachineStatusRefreshFunc(service ecloudservice.ECloudService, vmid int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		vm, err := service.GetVirtualMachine(vmid)
		if err != nil {
			if _, ok := err.(*ecloudservice.VirtualMachineNotFoundError); ok {
				return vm, "Deleted", nil
			}
			return nil, "", err
		}

		return vm, vm.Status.String(), nil
	}
}
