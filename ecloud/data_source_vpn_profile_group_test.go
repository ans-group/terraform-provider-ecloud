package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceVPNProfileGroup_basic(t *testing.T) {
	config := testAccDataSourceVPNProfileGroupConfig_basic(ANS_TEST_VPN_PROFILE_GROUP_ID)
	resourceName := "data.ecloud_vpn_profile_group.test-vpnprofilegroup"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", ANS_TEST_VPN_PROFILE_GROUP_ID),
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
