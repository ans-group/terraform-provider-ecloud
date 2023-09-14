package ecloud

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]func() (*schema.Provider, error)
var testAccProvider *schema.Provider
var (
	ANS_TEST_VPC_REGION_ID        = os.Getenv("ANS_TEST_VPC_REGION_ID")
	ANS_TEST_VPN_PROFILE_GROUP_ID = os.Getenv("ANS_TEST_VPN_PROFILE_GROUP_ID")
)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]func() (*schema.Provider, error){
		"ecloud": func() (*schema.Provider, error) { return testAccProvider, nil },
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	testAccPreCheckRequiredEnvVars(t)
}

func testAccPreCheckRequiredEnvVars(t *testing.T) {
	if ANS_TEST_VPC_REGION_ID == "" {
		t.Fatal("ANS_TEST_VPC_REGION_ID must be set for acceptance tests")
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
