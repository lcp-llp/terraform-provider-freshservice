package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API key for Freshservice",
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain for Freshservice (e.g., 'yourdomain.freshservice.com')",
			},
		},
		ConfigureContextFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"freshservice_asset_type":        resourceAssetType(),
			"freshservice_azure_subscription": resourceAzureSubscription(),
			"freshservice_aws_account":       resourceAWSAccount(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"freshservice_asset":      dataSourceAsset(),
			"freshservice_asset_type": dataSourceAssetType(),
		},
	}
}

// configureProvider configures the provider with authentication
func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiKey := d.Get("api_key").(string)
	domain := d.Get("domain").(string)

	config, err := NewConfig(apiKey, domain)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return config, nil
}
