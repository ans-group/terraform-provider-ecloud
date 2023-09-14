package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAffinityRuleMember_basic(t *testing.T) {
	affinityRuleName := acctest.RandomWithPrefix("tftest")
	affinityRuleMemberInstanceID := "ecloud_instance.test-instance.id"
	armResourceName := "ecloud_affinityrule_member.test-arm"
	arResourceName := "ecloud_affinityrule.test-ar"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckAffinityRuleMemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAffinityRuleMemberConfig_basic(ANS_TEST_VPC_REGION_ID, affinityRuleName, affinityRuleMemberInstanceID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAffinityRuleMemberExists(armResourceName),
					resource.TestCheckResourceAttr(armResourceName, "name", affinityRuleName),
					resource.TestCheckResourceAttr(armResourceName, "instance_id", affinityRuleMemberInstanceID),
					resource.TestCheckResourceAttrPair(arResourceName, "id", armResourceName, "affinity_rule_id"),
				),
			},
		},
	})
}

func testAccCheckAffinityRuleMemberExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Affinity rule member ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetAffinityRuleMember(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.AffinityRuleMemberNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckAffinityRuleMemberDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_affinityrule_member" {
			continue
		}

		_, err := service.GetAffinityRuleMember(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Affinity rule member with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.AffinityRuleMemberNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceAffinityRuleMemberConfig_basic(regionID string, affinityRuleName string, affinityRuleMemberInstanceID string) string {
	return fmt.Sprintf(`
resource "ecloud_vpc" "test-vpc" {
	region_id = "%[1]s"
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

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "test-instance"
	image_id = data.ecloud_image.centos7.id
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
}

resource "ecloud_affinityrule" "test-ar" {
   vpc_id = ecloud_vpc.test-vpc.id
   availability_zone_id = data.ecloud_availability_zone.test-az.id
   name = "%[2]s"
   type = "anti-affinity"
}

resource "ecloud_affinityrule_member" "test-arm" {
	affinity_rule_id = ecloud_affinityrule.test-ar.id
	instance_id = %[3]s
}
`, regionID, affinityRuleName, affinityRuleMemberInstanceID)
}
