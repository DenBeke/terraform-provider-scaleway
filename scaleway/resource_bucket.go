package scaleway

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	api "github.com/nicolai86/scaleway-sdk"
)

func resourceScalewayBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceScalewayBucketCreate,
		Read:   resourceScalewayBucketRead,
		Delete: resourceScalewayBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the bucket",
			},
		},
	}
}

func resourceScalewayBucketRead(d *schema.ResourceData, m interface{}) error {

	s3client := m.(*Meta).s3Client

	_, err := s3client.ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(d.Get("name").(string)),
	})
	if err != nil {
		if serr, ok := err.(api.APIError); ok && serr.StatusCode == 404 {
			log.Printf("[DEBUG] Bucket %q was not found - removing from state!", d.Get("name").(string))
			d.SetId("")
			return nil
		}
	}
	return err
}

func resourceScalewayBucketCreate(d *schema.ResourceData, m interface{}) error {
	s3client := m.(*Meta).s3Client

	createBucketResponse, err := s3client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(d.Get("name").(string)),
	})
	if err != nil {
		return err
	}

	d.SetId(*createBucketResponse.Location)
	return nil
}

func resourceScalewayBucketDelete(d *schema.ResourceData, m interface{}) error {

	s3client := m.(*Meta).s3Client

	_, err := s3client.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(d.Get("name").(string)),
	})
	if err != nil {
		return err
	}
	return nil
}
