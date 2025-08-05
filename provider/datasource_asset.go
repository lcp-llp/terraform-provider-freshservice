package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AssetSearchResponse represents the API response for asset search
type AssetSearchResponse struct {
	Assets []Asset `json:"assets"`
}

func dataSourceAsset() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAssetRead,
		Description: "Data source to search for Freshservice assets by name, display_id, or asset_tag",

		Schema: map[string]*schema.Schema{
			// Search parameters (at least one required)
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the asset to search for",
			},
			"display_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Display ID of the asset to search for",
			},
			"asset_tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Asset tag to search for",
			},
			"trashed": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Include assets in trash (default: false)",
			},

			// Asset output fields
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ID of the asset",
			},
			"item_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Item ID of the asset",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the asset",
			},
			"asset_type_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Asset type ID",
			},
			"impact": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Impact level of the asset",
			},
			"author_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Author type of the asset",
			},
			"usage_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Usage type of the asset",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User ID assigned to the asset",
			},
			"location_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Location ID of the asset",
			},
			"department_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Department ID of the asset",
			},
			"agent_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Agent ID assigned to the asset",
			},
			"assigned_on": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date when the asset was assigned",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Creation timestamp of the asset",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last update timestamp of the asset",
			},
			"workspace_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Workspace ID of the asset",
			},
			"created_by_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source that created the asset",
			},
			"last_updated_by_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Source that last updated the asset",
			},
			"created_by_user": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User who created the asset",
			},
			"last_updated_by_user": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User who last updated the asset",
			},
			"sources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of sources for the asset",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"serial_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Serial number of the asset",
			},
			"mac_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "MAC addresses of the asset",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_addresses": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "IP addresses of the asset",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "UUID of the asset",
			},
			"imei_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IMEI number of the asset",
			},
		},
	}
}

func dataSourceAssetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Validate that at least one search parameter is provided
	name := d.Get("name").(string)
	displayID := d.Get("display_id").(int)
	assetTag := d.Get("asset_tag").(string)
	trashed := d.Get("trashed").(bool)

	if name == "" && displayID == 0 && assetTag == "" {
		return diag.Errorf("At least one of 'name', 'display_id', or 'asset_tag' must be provided")
	}

	// Build search query
	searchQuery := buildSearchQuery(name, displayID, assetTag)

	// URL encode the search query
	encodedQuery := url.QueryEscape(searchQuery)

	// Build the endpoint
	endpoint := fmt.Sprintf("/assets?search=%s", encodedQuery)
	if trashed {
		endpoint += "&trashed=true"
	}

	// Make the API request
	req, err := config.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	var searchResponse AssetSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	if len(searchResponse.Assets) == 0 {
		return diag.Errorf("No assets found matching the search criteria")
	}

	if len(searchResponse.Assets) > 1 {
		return diag.Errorf("Multiple assets found matching the search criteria. Please refine your search to return a single asset")
	}

	// Set the asset data
	asset := searchResponse.Assets[0]
	d.SetId(strconv.Itoa(asset.DisplayID))

	if err := d.Set("id", asset.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_id", asset.DisplayID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", asset.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", asset.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("asset_type_id", asset.AssetTypeID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("impact", asset.Impact); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("author_type", asset.AuthorType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("usage_type", asset.UsageType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("asset_tag", asset.AssetTag); err != nil {
		return diag.FromErr(err)
	}

	// Handle nullable fields
	if asset.UserID != nil {
		if err := d.Set("user_id", *asset.UserID); err != nil {
			return diag.FromErr(err)
		}
	}
	if asset.LocationID != nil {
		if err := d.Set("location_id", *asset.LocationID); err != nil {
			return diag.FromErr(err)
		}
	}
	if asset.DepartmentID != nil {
		if err := d.Set("department_id", *asset.DepartmentID); err != nil {
			return diag.FromErr(err)
		}
	}
	if asset.AgentID != nil {
		if err := d.Set("agent_id", *asset.AgentID); err != nil {
			return diag.FromErr(err)
		}
	}
	if asset.AssignedOn != nil {
		if err := d.Set("assigned_on", *asset.AssignedOn); err != nil {
			return diag.FromErr(err)
		}
	}
	if asset.CreatedByUser != nil {
		if err := d.Set("created_by_user", *asset.CreatedByUser); err != nil {
			return diag.FromErr(err)
		}
	}
	if asset.LastUpdatedByUser != nil {
		if err := d.Set("last_updated_by_user", *asset.LastUpdatedByUser); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("created_at", asset.CreatedAt.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_at", asset.UpdatedAt.Format(time.RFC3339)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("workspace_id", asset.WorkspaceID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("created_by_source", asset.CreatedBySource); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_updated_by_source", asset.LastUpdatedBySource); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("sources", asset.Sources); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("serial_number", asset.SerialNumber); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("mac_addresses", asset.MacAddresses); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("ip_addresses", asset.IPAddresses); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("uuid", asset.UUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("item_id", asset.ItemID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("imei_number", asset.IMEINumber); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// buildSearchQuery builds the search query string based on provided parameters
func buildSearchQuery(name string, displayID int, assetTag string) string {
	var queryParts []string

	if name != "" {
		// Escape single quotes in the name
		escapedName := escapeSearchValue(name)
		queryParts = append(queryParts, fmt.Sprintf("name:'%s'", escapedName))
	}

	if displayID != 0 {
		queryParts = append(queryParts, fmt.Sprintf("display_id:%d", displayID))
	}

	if assetTag != "" {
		escapedAssetTag := escapeSearchValue(assetTag)
		queryParts = append(queryParts, fmt.Sprintf("asset_tag:'%s'", escapedAssetTag))
	}

	// Join with AND operator if multiple criteria
	if len(queryParts) > 1 {
		query := ""
		for i, part := range queryParts {
			if i > 0 {
				query += " AND "
			}
			query += part
		}
		return fmt.Sprintf("\"%s\"", query)
	}

	// Single criterion
	return fmt.Sprintf("\"%s\"", queryParts[0])
}

// escapeSearchValue escapes single quotes in search values
func escapeSearchValue(value string) string {
	// Replace single quotes with escaped single quotes for Freshservice API
	return strings.ReplaceAll(value, "'", "\\'")
}
