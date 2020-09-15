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
	UKF_TEST_VPC_NAME                 = os.Getenv("UKF_TEST_VPC_NAME")
	UKF_TEST_AVAILABILITYZONE_NAME    = os.Getenv("UKF_TEST_AVAILABILITYZONE_NAME")
	UKF_TEST_NETWORK_NAME             = os.Getenv("UKF_TEST_NETWORK_NAME")
	UKF_TEST_DHCP_AVAILABILITYZONE_ID = os.Getenv("UKF_TEST_DHCP_AVAILABILITYZONE_ID")
	UKF_TEST_INSTANCE_NAME            = os.Getenv("UKF_TEST_INSTANCE_NAME")
	// UKF_TEST_VPN_NAME                 = os.Getenv("UKF_TEST_VPN_NAME")

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
	if UKF_TEST_VPC_NAME == "" {
		t.Fatal("UKF_TEST_VPC_NAME must be set for acceptance tests")
	}
	if UKF_TEST_AVAILABILITYZONE_NAME == "" {
		t.Fatal("UKF_TEST_AVAILABILITYZONE_NAME must be set for acceptance tests")
	}
	if UKF_TEST_NETWORK_NAME == "" {
		t.Fatal("UKF_TEST_NETWORK_NAME must be set for acceptance tests")
	}
	if UKF_TEST_DHCP_AVAILABILITYZONE_ID == "" {
		t.Fatal("UKF_TEST_DHCP_AVAILABILITYZONE_ID must be set for acceptance tests")
	}
	if UKF_TEST_INSTANCE_NAME == "" {
		t.Fatal("UKF_TEST_INSTANCE_NAME must be set for acceptance tests")
	}
	// if UKF_TEST_VPN_NAME == "" {
	// 	t.Fatal("UKF_TEST_VPN_NAME must be set for acceptance tests")
	// }
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
