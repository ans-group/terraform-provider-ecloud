package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceBackupGateway_basic(t *testing.T) {
	gatewayName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceBackupGatewayConfig_basic(gatewayName)
	resourceName := "data.ecloud_backup_gateway.test-gateway"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", gatewayName),
					resource.TestCheckResourceAttrSet(resourceName, "vpc_id"),
					resource.TestCheckResourceAttrSet(resourceName, "availability_zone_id"),
					resource.TestCheckResourceAttrSet(resourceName, "gateway_spec_id"),
				),
			},
		},
	})
}

func testAccDataSourceBackupGatewayConfig_basic(gatewayName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

data "ecloud_backup_gateway_spec" "test-spec" {
	name = "Medium"
}

resource "ecloud_backup_gateway" "test-gateway" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "%s"
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	gateway_spec_id = data.ecloud_backup_gateway_spec.test-spec.id
}

data "ecloud_backup_gateway" "test-gateway" {
    name = ecloud_backup_gateway.test-gateway.name
}
`, gatewayName)
}
