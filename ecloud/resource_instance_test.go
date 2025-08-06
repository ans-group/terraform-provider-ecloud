package ecloud

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"text/template"

	ecloudservice "github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccInstance_basic(t *testing.T) {
	instanceName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_instance.test-instance"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInstanceConfig_basic(instanceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
				),
			},
		},
	})
}

func TestAccInstance_withTags(t *testing.T) {
	instanceName := acctest.RandomWithPrefix("tftest")
	tagName := acctest.RandomWithPrefix("tftest-tag")
	resourceName := "ecloud_instance.test-instance"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceInstanceConfig_withTags(instanceName, tagName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
					resource.TestCheckResourceAttr(resourceName, "tag_ids.#", "1"),
				),
			},
		},
	})
}

func testAccCheckInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Instance ID is set")
		}

		service := testAccProvider.Meta().(ecloudservice.ECloudService)

		_, err := service.GetInstance(rs.Primary.ID)
		if err != nil {
			if _, ok := err.(*ecloudservice.InstanceNotFoundError); ok {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccCheckInstanceDestroy(s *terraform.State) error {
	service := testAccProvider.Meta().(ecloudservice.ECloudService)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecloud_instance" {
			continue
		}

		_, err := service.GetInstance(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Instance with ID [%s] still exists", rs.Primary.ID)
		}

		if _, ok := err.(*ecloudservice.InstanceNotFoundError); ok {
			return nil
		}

		return err
	}

	return nil
}

func testAccResourceInstanceConfig_basic(instanceName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
}

data "ecloud_image" "centos7" {
	name = "CentOS 7"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "tftest-network"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "%s"
	image_id = data.ecloud_image.centos7.id
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
}
`, instanceName)
}

func testAccResourceInstanceConfig_withTags(instanceName, tagName string) string {
	return fmt.Sprintf(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
}

data "ecloud_image" "centos7" {
	name = "CentOS 7"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "tftest-network"
}

resource "ecloud_tag" "test-tag" {
	name = "%s"
	scope = "instance"
}

resource "ecloud_instance" "test-instance" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "%s"
	image_id = data.ecloud_image.centos7.id
	volume_capacity = 20
	ram_capacity = 1024
	vcpu_cores = 1
	tag_ids = [ecloud_tag.test-tag.id]
}
`, tagName, instanceName)
}

type vcpuTestConfig struct {
	Name               string
	VCPUCores          int
	VCPUSockets        int
	VCPUCoresPerSocket int
}

// TestAccInstance_vcpu tests vcpu_cores / vcpu.sockets / vcpu.cores_per_socket options. May require test timeout
// increasing with `-timeout 15m`.
func TestAccInstance_vcpu(t *testing.T) {
	instanceName := acctest.RandomWithPrefix("tftest")
	resourceName := "ecloud_instance.test-instance-vcpu"
	vpcResourceName := "ecloud_vpc.test-vpc"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			// 1. Create instance using vcpu_cores
			{
				Config: testAccResourceInstanceConfig_vcpu(vcpuTestConfig{instanceName, 1, 0, 0}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", instanceName),
					resource.TestCheckResourceAttrPair(vpcResourceName, "id", resourceName, "vpc_id"),
					resource.TestCheckResourceAttr(resourceName, "vcpu_cores", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "vcpu.0.sockets"),
					resource.TestCheckNoResourceAttr(resourceName, "vcpu.0.cores_per_socket"),
				),
			},
			// 2. Increase with vcpu_cores
			{
				Config: testAccResourceInstanceConfig_vcpu(vcpuTestConfig{instanceName, 2, 0, 0}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vcpu_cores", "2"),
					resource.TestCheckNoResourceAttr(resourceName, "vcpu.0.sockets"),
					resource.TestCheckNoResourceAttr(resourceName, "vcpu.0.cores_per_socket"),
				),
			},
			// 3. Decrease with vcpu_cores
			{
				Config: testAccResourceInstanceConfig_vcpu(vcpuTestConfig{instanceName, 1, 0, 0}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vcpu_cores", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "vcpu.0.sockets"),
					resource.TestCheckNoResourceAttr(resourceName, "vcpu.0.cores_per_socket"),
				),
			},
			// 4. Switch from vcpu_cores to using vcpu.sockets/vcpu.cores_per_socket.
			{
				Config: testAccResourceInstanceConfig_vcpu(vcpuTestConfig{instanceName, 0, 1, 1}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vcpu.0.sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "vcpu.0.cores_per_socket", "1"),
					resource.TestCheckResourceAttr(resourceName, "vcpu_cores", "0"),
				),
			},
			// 5. Increase our CPU count with vcpu.sockets/vcpu.cores_per_socket
			{
				Config: testAccResourceInstanceConfig_vcpu(vcpuTestConfig{instanceName, 0, 2, 2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vcpu.0.sockets", "2"),
					resource.TestCheckResourceAttr(resourceName, "vcpu.0.cores_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceName, "vcpu_cores", "0"),
				),
			},
			// 6. Decrease our CPU count with vcpu.sockets/vcpu.cores_per_socket
			{
				Config: testAccResourceInstanceConfig_vcpu(vcpuTestConfig{instanceName, 0, 1, 2}),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vcpu.0.sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "vcpu.0.cores_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceName, "vcpu_cores", "0"),
				),
			},
			// 7. Attempt to increase our CPU count with vcpu_cores - this should error
			{
				Config:      testAccResourceInstanceConfig_vcpu(vcpuTestConfig{instanceName, 2, 0, 0}),
				ExpectError: regexp.MustCompile(`.*vcpu cores can't be updated.*`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "vcpu.0.sockets", "1"),
					resource.TestCheckResourceAttr(resourceName, "vcpu.0.cores_per_socket", "2"),
					resource.TestCheckResourceAttr(resourceName, "vcpu_cores", "0"),
				),
			},
		},
	})
}

func testAccResourceInstanceConfig_vcpu(config vcpuTestConfig) string {
	if config.VCPUCores > 0 && (config.VCPUSockets > 0 || config.VCPUCoresPerSocket > 0) {
		panic("bad test config, VCPUCores exclusive with VCPUSockets/VCPUCoresPerSocket")
	}

	tmpl := template.Must(template.New("tf1").Parse(`
data "ecloud_region" "test-region" {
	name = "Manchester"
}

resource "ecloud_vpc" "test-vpc" {
	region_id = data.ecloud_region.test-region.id
	name = "tftest-vpc"
}

data "ecloud_image" "ubuntu2204" {
	name = "Ubuntu Server 22.04 LTS"
}

data "ecloud_availability_zone" "test-az" {
	name = "Manchester West"
}

resource "ecloud_router" "test-router" {
	vpc_id = ecloud_vpc.test-vpc.id
	availability_zone_id = data.ecloud_availability_zone.test-az.id
	name = "tftest-router"
}

resource "ecloud_network" "test-network" {
	router_id = ecloud_router.test-router.id
	name = "tftest-network"
    subnet = "10.0.0.0/24"
}

resource "ecloud_instance" "test-instance-vcpu" {
	vpc_id = ecloud_vpc.test-vpc.id
	network_id = ecloud_network.test-network.id
	name = "{{.Name}}"
	image_id = data.ecloud_image.ubuntu2204.id
	volume_capacity = 40
	ram_capacity = 1024
{{if gt .VCPUCores 0 }}
	vcpu_cores = {{.VCPUCores}}
{{end}}{{if or (gt .VCPUSockets 0) (gt .VCPUCoresPerSocket 0)}}
	vcpu {
		sockets = {{.VCPUSockets}}
		cores_per_socket = {{.VCPUCoresPerSocket}}
	}
{{end}}
}
`))

	buf := bytes.Buffer{}
	err := tmpl.Execute(&buf, config)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
