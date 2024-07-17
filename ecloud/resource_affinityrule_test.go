package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAffinityRule_basic(t *testing.T) {
	affinityRuleName := acctest.RandomWithPrefix("tftest")
	affinityRuleType := "anti-affinity"
	resourceName := "ecloud_affinityrule.test-ar"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAffinityRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAffinityRuleConfig_basic(affinityRuleName, affinityRuleType),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAffinityRuleExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", affinityRuleName),
					resource.TestCheckResourceAttr(resourceName, "type", affinityRuleType),
				),
			},
		},
	})
}

func testAccCheckAffinityRuleExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Affinity rule ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetAffinityRule(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.AffinityRuleNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckAffinityRuleDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_affinityrule" {
			continue
		}

		_, err := service.GetAffinityRule(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Affinity rule with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.AffinityRuleNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceAffinityRuleConfig_basic(affinityRuleName string, affinityRuleType string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "test-vpc"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_affinityrule" "test-ar" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "%[1]s"
	type = "%[2]s"
}
`, affinityRuleName, affinityRuleType)
}
