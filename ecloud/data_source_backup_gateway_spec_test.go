package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBackupGatewaySpecification_basic(t *testing.T) {
	specName := "Medium"
	config := testAccDataSourceBackupGatewaySpecificationConfig_basic(specName)
	resourceName := "data.ecloud_backup_gateway_spec.test-bgwspec"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", specName),
				),
			},
		},
	})
}

func TestAccDataSourceBackupGatewaySpecification_byID(t *testing.T) {
	specID := "bgws-b84122e0"
	config := testAccDataSourceBackupGatewaySpecificationConfig_byID(specID)
	resourceName := "data.ecloud_backup_gateway_spec.test-bgwspec"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "description"),
				),
			},
		},
	})
}

func testAccDataSourceBackupGatewaySpecificationConfig_basic(specName string) string {
	return fmt.Sprintf(`
data "ecloud_backup_gateway_spec" "test-bgwspec" {
    name = "%s"
}
`, specName)
}

func testAccDataSourceBackupGatewaySpecificationConfig_byID(specID string) string {
	return fmt.Sprintf(`
data "ecloud_backup_gateway_spec" "test-bgwspec" {
    backup_gateway_specification_id = "%s"
}
`, specID)
}
