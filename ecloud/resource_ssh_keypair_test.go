package ecloud

import (
	"fmt"
	"testing"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSshKeyPair_basic(t *testing.T) {
	keyPairName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_ssh_keypair.test-keypair"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSshKeyPairDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceSSHKeyPairConfig_basic(keyPairName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSshKeyPairExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", keyPairName),
				),
			},
		},
	})
}

func testAccCheckSshKeyPairExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ssh key pair ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetSSHKeyPair(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.SSHKeyPairNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckSshKeyPairDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_ssh_keypair" {
			continue
		}

		_, err := service.GetSSHKeyPair(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("ssh key pair with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.SSHKeyPairNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceSSHKeyPairConfig_basic(keyPairName string) string {
	return fmt.Sprintf(`
resource "ecloud_ssh_keypair" "test-keypair" {
	name = "%s"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDXIismybCTbE4p24LX/Aioi17UdLUrfolbwf1fKUD2a5Ps0xvZv3U19FRTo+x6yWux7kd78DpZ50CS4WkRs09QLP9K65hSZj/SJBXl+MNaz3pJ0FngZBXxTgxdJ82gcLCvY3iDBfn61PdrJTv6kLR4ZnZruj2kBND4yUZAyQKxfzrXD20UwlF1GWwE4lHuWXaEei4mGbHSeWVay0pOEf5d6uAWlsBm2JEdXkG7/LupdLh7z+RlEaTigHarlTbpcfCC82JX94IGWmiKToFr6+lX6y7QoVxd8pmEGIV/9dxPwWM/9RczSD2Oxum83ESPhVvQrBUTjE7T7fGoLlr31rQep+qgH5XdfCqkmiZ69NFDUEPIwiqpCKazli/Jdaxz6FsxlWZbmaMOW1cMAhAtxmpOxukbhB5hmJjzR3DTAEsv5euINFNxk8snY3b77JmDYX09yb+hT/fLyjBonc7I0RmFsUIV+H25yzh57iJoSuP9Qbz2RD4nIwxn5/PvKDbwElE= test-keypair"
}
`, keyPairName)
}
