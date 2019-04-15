package ecloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccVirtualMachine_basic(t *testing.T) {
	var vm ecloudservice.VirtualMachine
	vmName := acctest.RandomWithPrefix("tftest")
	config, err := testAccCheckVirtualMachineConfig_basic(vmName)
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "ecloud_virtualmachine.test-vm"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualMachineDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMachineExists(resourceName, &vm),
					resource.TestCheckResourceAttr(resourceName, "environment", UKF_TEST_SOLUTION_ENVIRONMENT),
					resource.TestCheckResourceAttr(resourceName, "template", "CentOS 7 64-bit"),
					resource.TestCheckResourceAttr(resourceName, "cpu", "2"),
					resource.TestCheckResourceAttr(resourceName, "ram", "2"),
					resource.TestCheckResourceAttr(resourceName, "os_disk", "20"),
					resource.TestCheckResourceAttr(resourceName, "solution_id", UKF_TEST_SOLUTION_ID),
					resource.TestCheckResourceAttr(resourceName, "power_status", ecloudservice.VirtualMachinePowerStatusOnline.String()),
					resource.TestCheckResourceAttr(resourceName, "name", vmName),
				),
			},
		},
	})
}

func testAccCheckVirtualMachineExists(n string, vm *ecloudservice.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No virtual machine ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		vmID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		getVM, err := service.GetVirtualMachine(vmID)
		if err != nil {
			if _, ok := err.(*ecloudservice.VirtualMachineNotFoundError); ok {
				return nil
			}
			return err
		}

		*vm = getVM

		return nil
	}
}

func testAccCheckVirtualMachineDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_virtualmachine" {
			continue
		}

		vmID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = service.GetVirtualMachine(vmID)
		if err == nil {
			return fmt.Errorf("Virtual machine with ID [%d] still exists", vmID)
		}

		if _, ok := err.(*ecloudservice.VirtualMachineNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

var testAccCheckVirtualMachineConfigTemplate_basic = `
resource "ecloud_virtualmachine" "test-vm" {
  cpu = 2
  ram = 2
  os_disk = 20
  template = "CentOS 7 64-bit"
  name = "{{ .VMName }}"
  environment = "{{ .UKF_TEST_SOLUTION_ENVIRONMENT }}"
  solution_id = {{ .UKF_TEST_SOLUTION_ID }}
}`

func testAccCheckVirtualMachineConfig_basic(vmName string) (string, error) {
	data := map[string]interface{}{
		"VMName":                        vmName,
		"UKF_TEST_SOLUTION_ENVIRONMENT": UKF_TEST_SOLUTION_ENVIRONMENT,
		"UKF_TEST_SOLUTION_ID":          UKF_TEST_SOLUTION_ID,
	}

	return testAccTemplateConfig(testAccCheckVirtualMachineConfigTemplate_basic, data)
}
