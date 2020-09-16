package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVPC(t *testing.T) {
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
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_REFERENCE_VPC_NAME),
				),
			},
		},
	})
}

var testAccCheckDataSourceVPCConfigTemplate = `
data "ecloud_vpc" "test-vpc" {
    vpc_id = "{{ .UKF_TEST_REFERENCE_VPC_ID }}"
}`

func testAccCheckDataSourceVPCConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_REFERENCE_VPC_ID": UKF_TEST_REFERENCE_VPC_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceVPCConfigTemplate, data)
}
