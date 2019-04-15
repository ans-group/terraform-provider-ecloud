package ecloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourcePod(t *testing.T) {
	var pod ecloudservice.Pod
	config, err := testAccCheckDataSourcePodConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_pod.test-pod"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourcePodExists(resourceName, &pod),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_POD_NAME),
				),
			},
		},
	})
}

func testAccCheckDataSourcePodExists(n string, pod *ecloudservice.Pod) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No pod ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		podID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		getPod, err := service.GetPod(podID)
		if err != nil {
			return err
		}

		*pod = getPod

		return nil
	}
}

var testAccCheckDataSourcePodConfigTemplate = `
data "ecloud_pod" "test-pod" {
	name = "{{ .UKF_TEST_POD_NAME }}"
}`

func testAccCheckDataSourcePodConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_POD_NAME": UKF_TEST_POD_NAME,
	}

	return testAccTemplateConfig(testAccCheckDataSourcePodConfigTemplate, data)
}
