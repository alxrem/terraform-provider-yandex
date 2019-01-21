package yandex

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/yandex-cloud/go-genproto/yandex/cloud/compute/v1"
)

func dataSourceYandexComputeImage() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceYandexComputeImageRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"family"},
			},
			"family": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"image_id"},
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceYandexComputeImageRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := context.Background()
	var image *compute.Image

	if v, ok := d.GetOk("image_id"); ok {
		imageID := v.(string)
		resp, err := config.sdk.Compute().Image().Get(ctx, &compute.GetImageRequest{
			ImageId: imageID,
		})

		if err != nil {
			return handleNotFoundError(err, d, fmt.Sprintf("image with ID %q", imageID))
		}

		image = resp
	} else if v, ok := d.GetOk("family"); ok {
		familyName := v.(string)

		folderID := StandardImagesFolderID
		if f, ok := d.GetOk("folder_id"); ok {
			folderID = f.(string)
		}

		resp, err := config.sdk.Compute().Image().GetLatestByFamily(ctx, &compute.GetImageLatestByFamilyRequest{
			FolderId: folderID,
			Family:   familyName,
		})

		if err != nil {
			return fmt.Errorf("failed to find latest image with family \"%s\": %s", familyName, err)
		}

		image = resp
	} else {
		return fmt.Errorf("one of 'image_id' or 'family' must be set")
	}

	createdAt, err := getTimestamp(image.CreatedAt)
	if err != nil {
		return err
	}

	d.Set("image_id", image.Id)
	d.Set("created_at", createdAt)
	d.Set("family", image.Family)
	d.Set("folder_id", image.FolderId)
	d.Set("name", image.Name)
	d.Set("description", image.Description)
	d.Set("status", strings.ToLower(image.Status.String()))
	d.Set("os_type", strings.ToLower(image.Os.Type.String()))
	d.Set("min_disk_size", toGigabytes(image.MinDiskSize))
	d.Set("size", toGigabytes(image.StorageSize))

	if err := d.Set("labels", image.Labels); err != nil {
		return err
	}

	if err := d.Set("product_ids", image.ProductIds); err != nil {
		return err
	}

	d.SetId(image.Id)

	return nil
}
