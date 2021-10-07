package ec2_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccDataSourceAwsVpcDhcpOptions_basic(t *testing.T) {
	resourceName := "aws_vpc_dhcp_options.test"
	datasourceName := "data.aws_vpc_dhcp_options.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config:      testAccDataSourceAwsVpcDhcpOptionsConfig_Missing,
				ExpectError: regexp.MustCompile(`No matching EC2 DHCP Options found`),
			},
			{
				Config: testAccDataSourceAwsVpcDhcpOptionsConfig_DhcpOptionsID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "dhcp_options_id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain_name", resourceName, "domain_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain_name_servers.#", resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain_name_servers.0", resourceName, "domain_name_servers.0"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain_name_servers.1", resourceName, "domain_name_servers.1"),
					resource.TestCheckResourceAttrPair(datasourceName, "netbios_name_servers.#", resourceName, "netbios_name_servers.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "netbios_name_servers.0", resourceName, "netbios_name_servers.0"),
					resource.TestCheckResourceAttrPair(datasourceName, "netbios_node_type", resourceName, "netbios_node_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "ntp_servers.#", resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "ntp_servers.0", resourceName, "ntp_servers.0"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.Name", resourceName, "tags.Name"),
					resource.TestCheckResourceAttrPair(datasourceName, "owner_id", resourceName, "owner_id"),
					resource.TestCheckResourceAttrPair(datasourceName, "arn", resourceName, "arn"),
				),
			},
		},
	})
}

func TestAccDataSourceAwsVpcDhcpOptions_Filter(t *testing.T) {
	rInt := sdkacctest.RandInt()
	resourceName := "aws_vpc_dhcp_options.test.0"
	datasourceName := "data.aws_vpc_dhcp_options.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, ec2.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceAwsVpcDhcpOptionsConfig_Filter(rInt, 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "dhcp_options_id", resourceName, "id"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain_name", resourceName, "domain_name"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain_name_servers.#", resourceName, "domain_name_servers.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain_name_servers.0", resourceName, "domain_name_servers.0"),
					resource.TestCheckResourceAttrPair(datasourceName, "domain_name_servers.1", resourceName, "domain_name_servers.1"),
					resource.TestCheckResourceAttrPair(datasourceName, "netbios_name_servers.#", resourceName, "netbios_name_servers.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "netbios_name_servers.0", resourceName, "netbios_name_servers.0"),
					resource.TestCheckResourceAttrPair(datasourceName, "netbios_node_type", resourceName, "netbios_node_type"),
					resource.TestCheckResourceAttrPair(datasourceName, "ntp_servers.#", resourceName, "ntp_servers.#"),
					resource.TestCheckResourceAttrPair(datasourceName, "ntp_servers.0", resourceName, "ntp_servers.0"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(datasourceName, "tags.Name", resourceName, "tags.Name"),
					resource.TestCheckResourceAttrPair(datasourceName, "owner_id", resourceName, "owner_id"),
				),
			},
			{
				Config:      testAccDataSourceAwsVpcDhcpOptionsConfig_Filter(rInt, 2),
				ExpectError: regexp.MustCompile(`Multiple matching EC2 DHCP Options found`),
			},
			{
				// We have one last empty step here because otherwise we'll leave the
				// test case with resources in the state and an erroneous config, and
				// thus the automatic destroy step will fail. This ensures we end with
				// both an empty state and a valid config.
				Config: `/* this config intentionally left blank */`,
			},
		},
	})
}

const testAccDataSourceAwsVpcDhcpOptionsConfig_Missing = `
data "aws_vpc_dhcp_options" "test" {
  dhcp_options_id = "does-not-exist"
}
`

const testAccDataSourceAwsVpcDhcpOptionsConfig_DhcpOptionsID = `
resource "aws_vpc_dhcp_options" "incorrect" {
  domain_name = "tf-acc-test-incorrect.example.com"
}

resource "aws_vpc_dhcp_options" "test" {
  domain_name          = "service.consul"
  domain_name_servers  = ["127.0.0.1", "10.0.0.2"]
  netbios_name_servers = ["127.0.0.1"]
  netbios_node_type    = 2
  ntp_servers          = ["127.0.0.1"]

  tags = {
    Name = "tf-acc-test"
  }
}

data "aws_vpc_dhcp_options" "test" {
  dhcp_options_id = aws_vpc_dhcp_options.test.id
}
`

func testAccDataSourceAwsVpcDhcpOptionsConfig_Filter(rInt, count int) string {
	return fmt.Sprintf(`
resource "aws_vpc_dhcp_options" "incorrect" {
  domain_name = "tf-acc-test-incorrect.example.com"
}

resource "aws_vpc_dhcp_options" "test" {
  count = %[2]d

  domain_name          = "tf-acc-test-%[1]d.example.com"
  domain_name_servers  = ["127.0.0.1", "10.0.0.2"]
  netbios_name_servers = ["127.0.0.1"]
  netbios_node_type    = 2
  ntp_servers          = ["127.0.0.1"]

  tags = {
    Name = "tf-acc-test-%[1]d"
  }
}

data "aws_vpc_dhcp_options" "test" {
  filter {
    name   = "key"
    values = ["domain-name"]
  }

  filter {
    name   = "value"
    values = [aws_vpc_dhcp_options.test[0].domain_name]
  }
}
`, rInt, count)
}
