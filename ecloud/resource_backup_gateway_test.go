package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBackupGateway_basic(t *testing.T) {
	backupGatewayName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_backup_gateway.test-bg"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBackupGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBackupGatewayConfig_basic(backupGatewayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", backupGatewayName),
					resource.TestCheckResourceAttrPair(resourceName, "vpc_id", vpcResourceName, "id"),
				),
			},
		},
	})
}

func TestAccBackupGateway_update(t *testing.T) {
	backupGatewayName := acctest.RandomWithPrefix("tftest")
	updatedName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_backup_gateway.test-bg"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBackupGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceBackupGatewayConfig_basic(backupGatewayName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", backupGatewayName),
				),
			},
			{
				Config: testAccResourceBackupGatewayConfig_basic(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBackupGatewayExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", updatedName),
				),
			},
		},
	})
}

func testAccCheckBackupGatewayExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Backup Gateway ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetBackupGateway(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.BackupGatewayNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckBackupGatewayDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_backup_gateway" {
			continue
		}

		_, err := service.GetBackupGateway(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Backup Gateway with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.BackupGatewayNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceBackupGatewayConfig_basic(backupGatewayName string) string {
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
	name = "Small"
}

resource "ecloud_backup_gateway" "test-bg" {
	vpc_id = ecloud_vpc.test-vpc.id
	name = "%s"
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	gateway_spec_id = data.ecloud_backup_gateway_spec.test-spec.id
}
`, backupGatewayName)
}
