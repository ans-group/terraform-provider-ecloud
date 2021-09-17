package ecloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceSshKeyPair_basic(t *testing.T) {
	keyPairName := acctest.RandomWithPrefix("tftest")
	config := testAccDataSourceSshKeyPairConfig_basic(keyPairName)
	resourceName := "data.ecloud_ssh_keypair.test-keypair"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", keyPairName),
				),
			},
		},
	})
}

func testAccDataSourceSshKeyPairConfig_basic(keyPairName string) string {
	return fmt.Sprintf(`

resource "ecloud_ssh_keypair" "test-keypair" {
	name = "%s"
	public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDXIismybCTbE4p24LX/Aioi17UdLUrfolbwf1fKUD2a5Ps0xvZv3U19FRTo+x6yWux7kd78DpZ50CS4WkRs09QLP9K65hSZj/SJBXl+MNaz3pJ0FngZBXxTgxdJ82gcLCvY3iDBfn61PdrJTv6kLR4ZnZruj2kBND4yUZAyQKxfzrXD20UwlF1GWwE4lHuWXaEei4mGbHSeWVay0pOEf5d6uAWlsBm2JEdXkG7/LupdLh7z+RlEaTigHarlTbpcfCC82JX94IGWmiKToFr6+lX6y7QoVxd8pmEGIV/9dxPwWM/9RczSD2Oxum83ESPhVvQrBUTjE7T7fGoLlr31rQep+qgH5XdfCqkmiZ69NFDUEPIwiqpCKazli/Jdaxz6FsxlWZbmaMOW1cMAhAtxmpOxukbhB5hmJjzR3DTAEsv5euINFNxk8snY3b77JmDYX09yb+hT/fLyjBonc7I0RmFsUIV+H25yzh57iJoSuP9Qbz2RD4nIwxn5/PvKDbwElE= test-keypair"
}

data "ecloud_ssh_keypair" "test-keypair" {
    name = ecloud_ssh_keypair.test-keypair.name
}
`, keyPairName)
}
