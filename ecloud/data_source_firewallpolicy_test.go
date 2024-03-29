package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceFirewallPolicy_basic(t *testing.T) {
	params := map[string]string{
		"vpc_region_id":   ANS_TEST_VPC_REGION_ID,
		"policy_name":     acctest.RandomWithPrefix("tftest"),
		"policy_sequence": "0",
	}
	config := testAccDataSourceFirewallPolicyConfig_basic(params)
	resourceName := "data.ecloud_firewallpolicy.test-fwp"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", params["policy_name"]),
					resource.TestCheckResourceAttr(resourceName, "sequence", params["policy_sequence"]),
				),
			},
		},
	})
}

func testAccDataSourceFirewallPolicyConfig_basic(params map[string]string) string {
	str, _ := testAccTemplateConfig(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "{{ .vpc_region_id }}"
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "test-router"
}

resource "ecloud_firewallpolicy" "test-fwp" {
	router_id = ecloud_router.test-router.id
	name = "{{ .policy_name }}"
	sequence = {{ .policy_sequence }}
}

data "ecloud_firewallpolicy" "test-fwp" {
    name = ecloud_firewallpolicy.test-fwp.name
}
`, params)
	return str
}
