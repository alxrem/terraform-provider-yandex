package yandex

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceYandexIAMToken() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceYandexIAMTokenRead,
		Schema: map[string]*schema.Schema{
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceYandexIAMTokenRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	ctx := config.Context()

	response, err := config.sdk.CreateIAMToken(ctx)
	if err != nil {
		return err
	}

	token := response.GetIamToken()

	d.Set("token", token)
	d.SetId(fmt.Sprintf("%d", schema.HashString(token)))

	return nil
}
