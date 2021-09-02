package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceVPNProfileGroup_basic(t *testing.T) {
	config := testAccDataSourceVPNProfileGroupConfig_basic(UKF_TEST_VPN_PROFILE_GROUP_ID)
	resourceName := "data.ecloud_vpn_profile_group.test-vpnprofilegroup"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_VPN_PROFILE_GROUP_ID),
				),
			},
		},
	})
}

func testAccDataSourceVPNProfileGroupConfig_basic(vpnProfileGroupName string) string {
	return fmt.Sprintf(`
data "ecloud_vpn_profile_group" "test-vpnprofilegroup" {
	name = "%[1]s"
}
`, vpnProfileGroupName)
}
