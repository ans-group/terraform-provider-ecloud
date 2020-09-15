package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceVPC(t *testing.T) {
	var vpc ecloudservice.VPC
	config, err := testAccCheckDataSourceVPCConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceVPCExists(resourceName, &vpc),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_VPC_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourceVPCExists(n string, vpc *ecloudservice.VPC) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No VPC ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		getVPC, err := service.GetVPC(rs.Primary.ID)
		if err != nil {
			return err
		}

		*vpc = getVPC

		return nil
	}
}

var testAccCheckDataSourceVPCConfigTemplate = `
data "ecloud_vpc" "test-vpc" {
    name = "{{ .UKF_TEST_VPC_NAME }}"
}`

func testAccCheckDataSourceVPCConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_VPC_NAME": UKF_TEST_VPC_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourceVPCConfigTemplate, data)
}
