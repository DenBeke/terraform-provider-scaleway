package scaleway

import (
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("scaleway_bucket", &resource.Sweeper{
		Name: "scaleway_bucket",
		F:    testSweepBucket,
	})
}

func testSweepBucket(region string) error {

	log.Println(region)

	s3client, err := sharedS3Client(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	listBucketResponse, err := s3client.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return fmt.Errorf("couldn't list buckets: %s", err)
	}

	for _, bucket := range listBucketResponse.Buckets {
		log.Println(*bucket.Name)
		if strings.HasPrefix(*bucket.Name, "terraform-test") {
			_, err := s3client.DeleteBucket(&s3.DeleteBucketInput{
				Bucket: bucket.Name,
			})
			if err != nil {
				return fmt.Errorf("Error deleting bucket in Sweeper: %s", err)
			}
		}

	}

	return nil
}

var testBucketName = fmt.Sprintf("terraform-test-%d", time.Now().Unix())

func TestAccScalewayBucket(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckScalewayBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckScalewayBucket,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("scaleway_bucket.base", "name", testBucketName),
				),
			},
		},
	})
}

func testAccCheckScalewayBucketDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Meta).deprecatedClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "scaleway" {
			continue
		}

		_, err := client.ListObjects(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Bucket still exists")
		}
	}

	return nil
}

var testAccCheckScalewayBucket = fmt.Sprintf(`
resource "scaleway_bucket" "base" {
  name = "%s"
}
`, testBucketName)
