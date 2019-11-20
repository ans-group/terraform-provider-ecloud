package ecloud

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var (
	UKF_TEST_SOLUTION_ENVIRONMENT       = os.Getenv("UKF_TEST_SOLUTION_ENVIRONMENT")
	UKF_TEST_SOLUTION_ID                = os.Getenv("UKF_TEST_SOLUTION_ID")
	UKF_TEST_SOLUTION_NAME              = os.Getenv("UKF_TEST_SOLUTION_NAME")
	UKF_TEST_SOLUTION_POD_ID            = os.Getenv("UKF_TEST_SOLUTION_POD_ID")
	UKF_TEST_SOLUTION_DATASTORE_NAME    = os.Getenv("UKF_TEST_SOLUTION_DATASTORE_NAME")
	UKF_TEST_SOLUTION_SITE_POD_ID       = os.Getenv("UKF_TEST_SOLUTION_SITE_POD_ID")
	UKF_TEST_SOLUTION_NETWORK_NAME      = os.Getenv("UKF_TEST_SOLUTION_NETWORK_NAME")
	UKF_TEST_SOLUTION_TEMPLATE_NAME     = os.Getenv("UKF_TEST_SOLUTION_TEMPLATE_NAME")
	UKF_TEST_SOLUTION_TEMPLATE_PLATFORM = os.Getenv("UKF_TEST_SOLUTION_TEMPLATE_PLATFORM")
	UKF_TEST_TEMPLATE_NAME              = os.Getenv("UKF_TEST_TEMPLATE_NAME")
	UKF_TEST_TEMPLATE_PLATFORM          = os.Getenv("UKF_TEST_TEMPLATE_PLATFORM")
	UKF_TEST_POD_NAME                   = os.Getenv("UKF_TEST_POD_NAME")
	UKF_TEST_APPLIANCE_NAME             = os.Getenv("UKF_TEST_APPLIANCE_NAME")
	UKF_TEST_APPLIANCE_POD_ID           = os.Getenv("UKF_TEST_APPLIANCE_POD_ID")
	UKF_TEST_ACTIVE_DIRECTORY_NAME      = os.Getenv("UKF_TEST_ACTIVE_DIRECTORY_NAME")
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"ecloud": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	if UKF_TEST_SOLUTION_ENVIRONMENT == "" {
		t.Fatal("UKF_TEST_SOLUTION_ENVIRONMENT must be set for acceptance tests")
	}
	if UKF_TEST_SOLUTION_ID == "" {
		t.Fatal("UKF_TEST_SOLUTION_ID must be set for acceptance tests")
	}
	if UKF_TEST_SOLUTION_NAME == "" {
		t.Fatal("UKF_TEST_SOLUTION_NAME must be set for acceptance tests")
	}
	if UKF_TEST_SOLUTION_POD_ID == "" {
		t.Fatal("UKF_TEST_SOLUTION_POD_ID must be set for acceptance tests")
	}
	if UKF_TEST_SOLUTION_DATASTORE_NAME == "" {
		t.Fatal("UKF_TEST_SOLUTION_DATASTORE_NAME must be set for acceptance tests")
	}
	if UKF_TEST_SOLUTION_SITE_POD_ID == "" {
		t.Fatal("UKF_TEST_SOLUTION_SITE_POD_ID must be set for acceptance tests")
	}
	if UKF_TEST_SOLUTION_NETWORK_NAME == "" {
		t.Fatal("UKF_TEST_SOLUTION_NETWORK_NAME must be set for acceptance tests")
	}
	if UKF_TEST_SOLUTION_TEMPLATE_NAME == "" {
		t.Fatal("UKF_TEST_SOLUTION_TEMPLATE_NAME must be set for acceptance tests")
	}
	if UKF_TEST_SOLUTION_TEMPLATE_PLATFORM == "" {
		t.Fatal("UKF_TEST_SOLUTION_TEMPLATE_PLATFORM must be set for acceptance tests")
	}
	if UKF_TEST_TEMPLATE_NAME == "" {
		t.Fatal("UKF_TEST_TEMPLATE_NAME must be set for acceptance tests")
	}
	if UKF_TEST_TEMPLATE_PLATFORM == "" {
		t.Fatal("UKF_TEST_TEMPLATE_PLATFORM must be set for acceptance tests")
	}
	if UKF_TEST_POD_NAME == "" {
		t.Fatal("UKF_TEST_POD_NAME must be set for acceptance tests")
	}
	if UKF_TEST_APPLIANCE_NAME == "" {
		t.Fatal("UKF_TEST_APPLIANCE_NAME must be set for acceptance tests")
	}
}

func testAccTemplateConfig(t string, i interface{}) (string, error) {
	tmpl, err := template.New("output").Parse(t)
	if err != nil {
		return "", fmt.Errorf("failed to create template: %s", err.Error())
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, i)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %s", err.Error())
	}

	return buf.String(), nil
}
