package ecloud

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceNic_basic(t *testing.T) {
	config := testAccDataSourceNicConfig_basic(UKF_TEST_VPC_REGION_ID)
	resourceName := "data.ecloud_nic.test-nic"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "ip_address", regexp.MustCompile(`^(10\.0\.0\.)(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`)),
				),
			},
		},
	})
}

func testAccDataSourceNicConfig_basic(regionID string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "test-network"
	subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "test-instance"
	image_id = "img-4ac43e58"
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
}

data "ecloud_nic" "test-nic" {
    nic_id = ecloud_instance.test-instance.nic_id
}
`, regionID)
}
