package ecloud

import (
	"fmt"
	"testing"

	"github.com/ans-group/sdk-go/pkg/connection"
	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIPAddressNICBinding_basic(t *testing.T) {
	nicIPAddressBindingResourceName := "ecloud_ipaddress.test-ipaddress"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAddressNICBindingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIPAddressNICBindingConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIPAddressNICBindingExists(nicIPAddressBindingResourceName),
				),
			},
		},
	})
}

func testAccCheckIPAddressNICBindingExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No IP address ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		nicID := rs.Primary.Attributes["nic_id"]
		ipAddressID := rs.Primary.Attributes["ip_address_id"]

		ipAddresses, err := service.GetNICIPAddresses(nicID, *connection.NewAPIRequestParameters().WithFilter(
			*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{ipAddressID}),
		))
		if err != nil {
			return fmt.Errorf("Failed to retrieve IP addresses for NIC: %s", err)
		}
		if len(ipAddresses) != 1 {
			return fmt.Errorf("IP address with ID [%s] still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccCheckIPAddressNICBindingDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_ipaddress" {
			continue
		}

		nicID := rs.Primary.Attributes["nic_id"]
		ipAddressID := rs.Primary.Attributes["ip_address_id"]

		ipAddresses, err := service.GetNICIPAddresses(nicID, *connection.NewAPIRequestParameters().WithFilter(
			*connection.NewAPIRequestFiltering("id", connection.EQOperator, []string{ipAddressID}),
		))
		if err != nil {
			return fmt.Errorf("Failed to retrieve IP addresses for NIC: %s", err)
		}
		if len(ipAddresses) != 1 {
			return fmt.Errorf("IP address with ID [%s] still exists", rs.Primary.ID)
		}

		return nil
	}

	return nil
}

func testAccResourceIPAddressNICBindingConfig_basic() string {
	return `
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

data "ecloud_image" "centos7" {
	name = "CentOS 7"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	image_id = data.ecloud_image.centos7.id
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
}

data "ecloud_nic" "nic-1" {
  instance_id = ecloud_instance.instance-1.id
}

resource "ecloud_ipaddress" "ipaddress-1" {
	network_id = ecloud_network.test-network.id
}

resource "ecloud_nic_ipaddress_binding" "ipaddress-1-binding-1" {
  nic_id = data.ecloud_nic.nic-1.id
  ip_address_id = ecloud_ipaddress.ipaddress-1.id
}
`
}
