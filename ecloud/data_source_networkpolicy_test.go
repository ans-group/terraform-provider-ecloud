package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceNetworkPolicy_basic(t *testing.T) {
	policyName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceNetworkPolicyConfig_basic(UKF_TEST_VPC_REGION_ID, policyName)
	resourceName := "data.ecloud_networkpolicy.test-network"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", policyName),
				),
			},
		},
	})
}

func testAccDataSourceNetworkPolicyConfig_basic(regionID string, policyName string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
	name = "test-vpc"
	advanced_networking = true
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "test-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	subnet = "10.0.0.0/24"
}

resource "ecloud_networkpolicy" "test-np" {
	network_id = ecloud_router.test-router.id
	name = "%[2]s"
}

data "ecloud_networkpolicy" "test-np" {
    name = ecloud_networkpolicy.test-np.name
}
`, regionID, policyName)
}
