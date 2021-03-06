package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudOssBucketObjectsDataSource_basic(t *testing.T) {
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOssBucketObjectsDataSourceBasic(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_oss_bucket_objects.objects"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.key", fmt.Sprintf("tf-sample/tf-testacc-bucket-object-ds-basic-%d-object", randInt)),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.acl", "default"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.content_type", "text/plain"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.content_length", "14"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.cache_control", "max-age=0"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.content_disposition", "attachment; filename=\"my-object\""),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.content_encoding", "gzip"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.content_md5", "1STMBJqp4X5QEQsYTbRmkQ=="),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.expires", "Wed, 06 May 2020 00:00:00 GMT"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_bucket_objects.objects", "objects.0.etag"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.storage_class", "Standard"),
					resource.TestCheckResourceAttrSet("data.alicloud_oss_bucket_objects.objects", "objects.0.last_modification_time"),
				),
			},
		},
	})
}

func TestAccAlicloudOssBucketObjectsDataSource_filterByPrefix(t *testing.T) {
	randInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudOssBucketObjectsDataSourceFilterByPrefix(randInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_oss_bucket_objects.objects"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.#", "1"),
					resource.TestCheckResourceAttr("data.alicloud_oss_bucket_objects.objects", "objects.0.key", fmt.Sprintf("tf-prefix1/tf-testacc-bucket-object-ds-prefix-%d-object", randInt)),
				),
			},
		},
	})
}

func testAccCheckAlicloudOssBucketObjectsDataSourceBasic(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-bucket-object-ds-basic-%d"
}

resource "alicloud_oss_bucket" "sample_bucket" {
	bucket = "${var.name}"
	acl = "private"
}

resource "alicloud_oss_bucket_object" "sample_object" {
	bucket = "${alicloud_oss_bucket.sample_bucket.bucket}"
	key = "tf-sample/${var.name}-object"
	content = "sample content"
	content_type = "text/plain"
	cache_control = "max-age=0"
	content_disposition = "attachment; filename=\"my-object\""
	content_encoding = "gzip"
	expires = "Wed, 06 May 2020 00:00:00 GMT"
}

data "alicloud_oss_bucket_objects" "objects" {
	bucket_name = "${var.name}"
    key_regex = "${alicloud_oss_bucket_object.sample_object.key}"
}
`, randInt)
}

func testAccCheckAlicloudOssBucketObjectsDataSourceFilterByPrefix(randInt int) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-bucket-object-ds-prefix-%d"
}

resource "alicloud_oss_bucket" "sample_bucket" {
	bucket = "${var.name}"
	acl = "private"
}

resource "alicloud_oss_bucket_object" "sample_prefix1_object" {
	bucket = "${alicloud_oss_bucket.sample_bucket.bucket}"
	key = "tf-prefix1/${var.name}-object"
	content = "sample content"
	content_type = "text/plain"
}

resource "alicloud_oss_bucket_object" "sample_prefix2_object" {
	bucket = "${alicloud_oss_bucket.sample_bucket.bucket}"
	key = "tf-prefix2/${var.name}-object"
	content = "sample content"
	content_type = "text/plain"
}

data "alicloud_oss_bucket_objects" "objects" {
	bucket_name = "${var.name}"
    key_regex = "(${alicloud_oss_bucket_object.sample_prefix1_object.key}|${alicloud_oss_bucket_object.sample_prefix2_object.key})"
    key_prefix = "tf-prefix1/"
}
`, randInt)
}
