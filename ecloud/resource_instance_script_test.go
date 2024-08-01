package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccInstanceScript(t *testing.T) {
	scriptResourceName := "ecloud_instance_script.test-script"
	scriptContent := "hostname"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPAddressDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInstanceScriptConfig_basic(scriptContent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceScriptExists(scriptResourceName),
					resource.TestCheckResourceAttr(scriptResourceName, "script", scriptContent),
				),
			},
		},
	})
}

func testAccCheckInstanceScriptExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Instance ID is set")
		}

		return nil
	}
}

func testAccResourceInstanceScriptConfig_basic(scriptContent string) string {
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

resource "ecloud_router" "test-router-1" {
  vpc_id               = ecloud_vpc.test-vpc.id
  availability_zone_id = data.ecloud_availability_zone.test-az.id
  name                 = "test-router"
}

resource "ecloud_network" "network-1" {
	router_id = ecloud_router.test-router-1.id
	subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "instance-1" {
  vcpu {
    sockets          = 1
    cores_per_socket = 2
  }
  ram_capacity    = 2048
  vpc_id          = ecloud_vpc.test-vpc.id
  name            = "instance test"
  image_id        = "img-19cb94e5"
  volume_capacity = 40
  volume_iops     = 600
  network_id      = ecloud_network.network-1.id
  backup_enabled  = false
  encrypted       = false
}

data "ecloud_instance_credential" "root_creds" {
  instance_id = ecloud_instance.instance-1.id
  name        = "root"
}

resource "ecloud_instance_script" "test-script" {
  instance_id = ecloud_instance.instance-1.id
  username = data.ecloud_instance_credential.root_creds.username
  password = data.ecloud_instance_credential.root_creds.password
  script = "%[1]s"
}
`, scriptContent)
}
