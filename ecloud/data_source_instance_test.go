package ecloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceInstance(t *testing.T) {
	config, err := testAccCheckDataSourceInstanceConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_instance.test-instance"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_REFERENCE_INSTANCE_NAME),
				),
			},
		},
	})
}

var testAccCheckDataSourceInstanceConfigTemplate = `
data "ecloud_instance" "test-instance" {
    instance_id = "{{ .UKF_TEST_REFERENCE_INSTANCE_ID }}"
}`

func testAccCheckDataSourceInstanceConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_REFERENCE_INSTANCE_ID": UKF_TEST_REFERENCE_INSTANCE_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceInstanceConfigTemplate, data)
}
