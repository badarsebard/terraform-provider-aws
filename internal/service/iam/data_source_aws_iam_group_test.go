package iam_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/iam"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/provider"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/verify"
)

func TestAccAWSDataSourceIAMGroup_basic(t *testing.T) {
	groupName := fmt.Sprintf("test-datasource-user-%d", sdkacctest.RandInt())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsIAMGroupDataSourceConfig(groupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws_iam_group.test", "group_id"),
					resource.TestCheckResourceAttr("data.aws_iam_group.test", "path", "/"),
					resource.TestCheckResourceAttr("data.aws_iam_group.test", "group_name", groupName),
					acctest.CheckResourceAttrGlobalARN("data.aws_iam_group.test", "arn", "iam", fmt.Sprintf("group/%s", groupName)),
				),
			},
		},
	})
}

func TestAccAWSDataSourceIAMGroup_users(t *testing.T) {
	groupName := fmt.Sprintf("test-datasource-group-%d", sdkacctest.RandInt())
	userName := fmt.Sprintf("test-datasource-user-%d", sdkacctest.RandInt())
	groupMemberShipName := fmt.Sprintf("test-datasource-group-membership-%d", sdkacctest.RandInt())
	userCount := 101

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:   func() { acctest.PreCheck(t) },
		ErrorCheck: acctest.ErrorCheck(t, iam.EndpointsID),
		Providers:  acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsIAMGroupDataSourceConfigWithUser(groupName, userName, groupMemberShipName, userCount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.aws_iam_group.test", "group_id"),
					resource.TestCheckResourceAttr("data.aws_iam_group.test", "path", "/"),
					resource.TestCheckResourceAttr("data.aws_iam_group.test", "group_name", groupName),
					acctest.CheckResourceAttrGlobalARN("data.aws_iam_group.test", "arn", "iam", fmt.Sprintf("group/%s", groupName)),
					resource.TestCheckResourceAttr("data.aws_iam_group.test", "users.#", fmt.Sprint(userCount)),
					resource.TestCheckResourceAttrSet("data.aws_iam_group.test", "users.0.arn"),
					resource.TestCheckResourceAttrSet("data.aws_iam_group.test", "users.0.user_id"),
					resource.TestCheckResourceAttrSet("data.aws_iam_group.test", "users.0.user_name"),
					resource.TestCheckResourceAttrSet("data.aws_iam_group.test", "users.0.path"),
				),
			},
		},
	})
}

func testAccAwsIAMGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "aws_iam_group" "group" {
  name = "%s"
  path = "/"
}

data "aws_iam_group" "test" {
  group_name = aws_iam_group.group.name
}
`, name)
}

func testAccAwsIAMGroupDataSourceConfigWithUser(groupName, userName, membershipName string, userCount int) string {
	return fmt.Sprintf(`
resource "aws_iam_group" "group" {
  name = "%s"
  path = "/"
}

resource "aws_iam_user" "user" {
  name  = "%s-${count.index}"
  count = %d
}

resource "aws_iam_group_membership" "team" {
  name  = "%s"
  users = aws_iam_user.user.*.name
  group = aws_iam_group.group.name
}

data "aws_iam_group" "test" {
  group_name = aws_iam_group_membership.team.group
}
`, groupName, userName, userCount, membershipName)
}
