package ecloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceSite(t *testing.T) {
	var site ecloudservice.Site
	config, err := testAccCheckDataSourceSiteConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_site.test-site"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceSiteExists(resourceName, &site),
					resource.TestCheckResourceAttr(resourceName, "solution_id", UKF_TEST_SOLUTION_ID),
					resource.TestCheckResourceAttr(resourceName, "pod_id", UKF_TEST_SOLUTION_SITE_POD_ID),
				),
			},
		},
	})
}

func testAccCheckDataSourceSiteExists(n string, site *ecloudservice.Site) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No site ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		siteID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		getSite, err := service.GetSite(siteID)
		if err != nil {
			return err
		}

		*site = getSite

		return nil
	}
}

var testAccCheckDataSourceSiteConfigTemplate = `
data "ecloud_site" "test-site" {
	pod_id = "{{ .UKF_TEST_SOLUTION_SITE_POD_ID }}"
	solution_id = "{{ .UKF_TEST_SOLUTION_ID }}"
}`

func testAccCheckDataSourceSiteConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_SOLUTION_SITE_POD_ID": UKF_TEST_SOLUTION_SITE_POD_ID,
		"UKF_TEST_SOLUTION_ID":          UKF_TEST_SOLUTION_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceSiteConfigTemplate, data)
}
