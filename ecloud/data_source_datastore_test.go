package ecloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	ecloudservice "github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func TestAccDataSourceDatastore(t *testing.T) {
	var datastore ecloudservice.Datastore
	config, err := testAccCheckDataSourceDatastoreConfig()
	if err != nil {
		t.Fatalf("Failed to generate config: %s", err)
	}

	resourceName := "data.ecloud_datastore.test-datastore"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDatastoreExists(resourceName, &datastore),
					resource.TestCheckResourceAttr(resourceName, "name", UKF_TEST_SOLUTION_DATASTORE_NAME),
					resource.TestCheckResourceAttr(resourceName, "solution_id", UKF_TEST_SOLUTION_ID),
				),
			},
		},
	})
}

func testAccCheckDataSourceDatastoreExists(n string, datastore *ecloudservice.Datastore) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No datastore ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		datastoreID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		getDatastore, err := service.GetDatastore(datastoreID)
		if err != nil {
			return err
		}

		*datastore = getDatastore

		return nil
	}
}

var testAccCheckDataSourceDatastoreConfigTemplate = `
data "ecloud_datastore" "test-datastore" {
	name = "{{ .UKF_TEST_SOLUTION_DATASTORE_NAME }}"
	solution_id = "{{ .UKF_TEST_SOLUTION_ID }}"
}`

func testAccCheckDataSourceDatastoreConfig() (string, error) {
	data := map[string]interface{}{
		"UKF_TEST_SOLUTION_DATASTORE_NAME": UKF_TEST_SOLUTION_DATASTORE_NAME,
		"UKF_TEST_SOLUTION_ID":             UKF_TEST_SOLUTION_ID,
	}

	return testAccTemplateConfig(testAccCheckDataSourceDatastoreConfigTemplate, data)
}
