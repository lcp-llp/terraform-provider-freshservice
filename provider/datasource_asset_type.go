package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AssetTypesListResponse represents the API response for listing asset types
type AssetTypesListResponse struct {
	AssetTypes []AssetType `json:"asset_types"`
}

func dataSourceAssetType() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetTypeRead,
		Description: "Data source to retrieve information about a Freshservice asset type",

		Schema: map[string]*schema.Schema{
			// Search parameters
			"id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the asset type to retrieve",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the asset type to search for",
			},

			// Output fields
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the asset type",
			},
			"parent_asset_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the parent asset type",
			},
			"visible": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Visibility of the asset type",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the asset type",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the asset type",
			},
		},
	}
}

func dataSourceAssetTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Check if ID is provided
	if id, ok := d.GetOk("id"); ok {
		return readAssetTypeByID(ctx, d, config, id.(int))
	}

	// Check if name is provided
	if name, ok := d.GetOk("name"); ok {
		return readAssetTypeByName(ctx, d, config, name.(string))
	}

	return diag.Errorf("Either 'id' or 'name' must be provided")
}

// readAssetTypeByID retrieves an asset type by its ID
func readAssetTypeByID(ctx context.Context, d *schema.ResourceData, config *Config, id int) diag.Diagnostics {
	endpoint := fmt.Sprintf("/asset_types/%d", id)
	req, err := config.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	var assetTypeResp AssetTypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetTypeResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	// Set the ID and data
	d.SetId(strconv.Itoa(assetTypeResp.AssetType.ID))
	return setAssetTypeDataSourceData(d, &assetTypeResp.AssetType)
}

// readAssetTypeByName retrieves an asset type by searching by name
func readAssetTypeByName(ctx context.Context, d *schema.ResourceData, config *Config, name string) diag.Diagnostics {
	// List all asset types and find the one with matching name
	endpoint := "/asset_types"
	req, err := config.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	var assetTypesResp AssetTypesListResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetTypesResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	// Find asset type with matching name
	var foundAssetType *AssetType
	for _, assetType := range assetTypesResp.AssetTypes {
		if assetType.Name == name {
			foundAssetType = &assetType
			break
		}
	}

	if foundAssetType == nil {
		return diag.Errorf("No asset type found with name: %s", name)
	}

	// Set the ID and data
	d.SetId(strconv.Itoa(foundAssetType.ID))
	return setAssetTypeDataSourceData(d, foundAssetType)
}

// setAssetTypeDataSourceData sets the asset type data for the data source
func setAssetTypeDataSourceData(d *schema.ResourceData, assetType *AssetType) diag.Diagnostics {
	if err := d.Set("id", assetType.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", assetType.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", assetType.Description); err != nil {
		return diag.FromErr(err)
	}
	if assetType.ParentAssetTypeID != nil {
		if err := d.Set("parent_asset_type_id", *assetType.ParentAssetTypeID); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("visible", assetType.Visible); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_at", assetType.CreatedAt.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_at", assetType.UpdatedAt.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
