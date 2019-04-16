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

func TestAccVirtualMachineTag_basic(t *testing.T) {
	var tag ecloudservice.Tag
	vmName := acctest.RandomWithPrefix("tftest")
	tagKey := acctest.RandomWithPrefix("tftest")
	config, err := testAccCheckVirtualMachineTagConfig_basic(vmName, tagKey, "bar")
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "ecloud_virtualmachine_tag.test-tag"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckVirtualMachineTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVirtualMachineTagExists(resourceName, &tag),
					resource.TestCheckResourceAttr(resourceName, "key", tagKey),
					resource.TestCheckResourceAttr(resourceName, "value", "bar"),
				),
			},
		},
	})
}

func testAccCheckVirtualMachineTagExists(n string, tag *ecloudservice.Tag) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No virtual machine tag ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		tagKey := rs.Primary.ID

		vmID, err := strconv.Atoi(rs.Primary.Attributes["virtualmachine_id"])
		if err != nil {
			return err
		}

		getTag, err := service.GetVirtualMachineTag(vmID, tagKey)
		if err != nil {
			if _, ok := err.(*ecloudservice.TagNotFoundError); ok {
				return nil
			}
			return err
		}

		*tag = getTag

		return nil
	}
}

func testAccCheckVirtualMachineTagDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_virtualmachine_tag" {
			continue
		}

		tagKey := rs.Primary.ID

		vmID, err := strconv.Atoi(rs.Primary.Attributes["virtualmachine_id"])
		if err != nil {
			return err
		}

		_, err = service.GetVirtualMachineTag(vmID, tagKey)
		if err == nil {
			return fmt.Errorf("Tag with key [%s] for virtual machine with ID [%d] still exists", tagKey, vmID)
		}

		if _, ok := err.(*ecloudservice.TagNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

var testAccCheckVirtualMachineTagConfigTemplate_basic = `
resource "ecloud_virtualmachine" "test-vm" {
  cpu = 2
  ram = 2
  os_disk = 20
  template = "CentOS 7 64-bit"
  name = "{{ .VMName }}"
  environment = "{{ .UKF_TEST_SOLUTION_ENVIRONMENT }}"
  solution_id = {{ .UKF_TEST_SOLUTION_ID }}
}

resource "ecloud_virtualmachine_tag" "test-tag" {
	virtualmachine_id = "${ecloud_virtualmachine.test-vm.id}"
    key = "{{ .TagKey }}"
    value = "{{ .TagValue }}"
}`

func testAccCheckVirtualMachineTagConfig_basic(vmName string, tagKey string, tagValue string) (string, error) {
	data := map[string]interface{}{
		"VMName":                        vmName,
		"UKF_TEST_SOLUTION_ENVIRONMENT": UKF_TEST_SOLUTION_ENVIRONMENT,
		"UKF_TEST_SOLUTION_ID":          UKF_TEST_SOLUTION_ID,
		"TagKey":                        tagKey,
		"TagValue":                      tagValue,
	}

	return testAccTemplateConfig(testAccCheckVirtualMachineTagConfigTemplate_basic, data)
}
