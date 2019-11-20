package ecloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceActiveDirectory(t *testing.T) {
	var domain ecloudservice.ActiveDirectoryDomain
	config, err := testAccCheckDataSourceActiveDirectoryConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_active_directory.test-active-directory"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceActiveDirectoryExists(resourceName, &domain),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_ACTIVE_DIRECTORY_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourceActiveDirectoryExists(n string, domain *ecloudservice.ActiveDirectoryDomain) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Active Directory domain ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		var domainID, err = strconv.Atoi(rs.Primary.ID)

		getActiveDirectoryDomain, err := service.GetActiveDirectoryDomain(domainID)
		if err != nil {
			return err
		}

		*domain = getActiveDirectoryDomain

		return nil
	}
}

var testAccCheckDataSourceActiveDirectoryConfigTemplate = `
data "ecloud_active_directory" "test-active-directory" {
    name = "{{ .UKF_TEST_ACTIVE_DIRECTORY_NAME }}"
}`

func testAccCheckDataSourceActiveDirectoryConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_ACTIVE_DIRECTORY_NAME": UKF_TEST_ACTIVE_DIRECTORY_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourceActiveDirectoryConfigTemplate, data)
}
