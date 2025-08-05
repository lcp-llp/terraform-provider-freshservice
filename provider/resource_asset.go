package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Asset represents a Freshservice asset
type Asset struct {
	ID                  int                    `json:"id"`
	DisplayID           int                    `json:"display_id"`
	Name                string                 `json:"name"`
	Description         string                 `json:"description"`
	AssetTypeID         int                    `json:"asset_type_id"`
	Impact              string                 `json:"impact"`
	AuthorType          string                 `json:"author_type"`
	UsageType           string                 `json:"usage_type"`
	AssetTag            string                 `json:"asset_tag"`
	UserID              *int                   `json:"user_id"`
	LocationID          *int                   `json:"location_id"`
	DepartmentID        *int                   `json:"department_id"`
	AgentID             *int                   `json:"agent_id"`
	GroupID             *int                   `json:"group_id"`
	AssignedOn          *string                `json:"assigned_on"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	WorkspaceID         int                    `json:"workspace_id"`
	CreatedBySource     string                 `json:"created_by_source"`
	LastUpdatedBySource string                 `json:"last_updated_by_source"`
	CreatedByUser       *int                   `json:"created_by_user"`
	LastUpdatedByUser   *int                   `json:"last_updated_by_user"`
	Sources             []string               `json:"sources"`
	SerialNumber        string                 `json:"serial_number,omitempty"`
	MacAddresses        []string               `json:"mac_addresses,omitempty"`
	IPAddresses         []string               `json:"ip_addresses,omitempty"`
	UUID                string                 `json:"uuid,omitempty"`
	ItemID              string                 `json:"item_id,omitempty"`
	IMEINumber          string                 `json:"imei_number,omitempty"`
	TypeFields          map[string]interface{} `json:"type_fields,omitempty"`
}

// AssetResponse represents the API response for asset operations
type AssetResponse struct {
	Asset Asset `json:"asset"`
}

// AssetRequest represents the request body for asset operations
type AssetRequest struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	AssetTypeID  int                    `json:"asset_type_id"`
	Impact       string                 `json:"impact,omitempty"`
	UsageType    string                 `json:"usage_type,omitempty"`
	UserID       *int                   `json:"user_id,omitempty"`
	LocationID   *int                   `json:"location_id,omitempty"`
	DepartmentID *int                   `json:"department_id,omitempty"`
	AgentID      *int                   `json:"agent_id,omitempty"`
	GroupID      *int                   `json:"group_id,omitempty"`
	TypeFields   map[string]interface{} `json:"type_fields,omitempty"`
}

func resourceAsset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetCreate,
		ReadContext:   resourceAssetRead,
		UpdateContext: resourceAssetUpdate,
		DeleteContext: resourceAssetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Manages a Freshservice asset with custom type fields",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the asset (contains display_id value)",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the asset",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the asset",
			},
			"asset_type_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true, // Asset type cannot be changed after creation
				Description: "Asset type ID",
			},
			"impact": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "low",
				Description: "Impact level of the asset (low, medium, high)",
			},
			"usage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "permanent",
				Description: "Usage type of the asset (permanent, loaner)",
			},
			"user_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "User ID assigned to the asset",
			},
			"location_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Location ID of the asset",
			},
			"department_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Department ID of the asset",
			},
			"agent_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Agent ID assigned to the asset",
			},
			"group_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Group ID assigned to the asset",
			},
			"type_fields": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Custom type fields specific to the asset type. Field names will automatically have the asset type ID appended (e.g., 'product' becomes 'product_25')",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			// Computed fields
			"display_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Display ID of the asset",
			},
			"asset_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Asset tag",
			},
			"author_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Author type of the asset",
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
			"item_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Item ID of the asset",
			},
			"imei_number": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IMEI number of the asset",
			},
		},
	}
}

func resourceAssetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Build type_fields from the type_fields map
	typeFields := make(map[string]interface{})
	if typeFieldsRaw, ok := d.GetOk("type_fields"); ok {
		assetTypeID := d.Get("asset_type_id").(int)
		for key, value := range typeFieldsRaw.(map[string]interface{}) {
			// Automatically append the asset type ID to the field name
			fieldKey := fmt.Sprintf("%s_%d", key, assetTypeID)
			// Convert string values to appropriate types based on common patterns
			typeFields[fieldKey] = convertTypeFieldValue(value.(string))
		}
	}

	// Build request body
	assetReq := AssetRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		AssetTypeID: d.Get("asset_type_id").(int),
		Impact:      d.Get("impact").(string),
		UsageType:   d.Get("usage_type").(string),
		TypeFields:  typeFields,
	}

	// Handle optional nullable fields
	if userID, ok := d.GetOk("user_id"); ok {
		uid := userID.(int)
		assetReq.UserID = &uid
	}
	if locationID, ok := d.GetOk("location_id"); ok {
		lid := locationID.(int)
		assetReq.LocationID = &lid
	}
	if departmentID, ok := d.GetOk("department_id"); ok {
		did := departmentID.(int)
		assetReq.DepartmentID = &did
	}
	if agentID, ok := d.GetOk("agent_id"); ok {
		aid := agentID.(int)
		assetReq.AgentID = &aid
	}
	if groupID, ok := d.GetOk("group_id"); ok {
		gid := groupID.(int)
		assetReq.GroupID = &gid
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(assetReq)
	if err != nil {
		return diag.Errorf("Failed to marshal request: %s", err)
	}

	// Create the request
	req, err := config.NewRequest(ctx, "POST", "/assets", bytes.NewReader(jsonData))
	if err != nil {
		return diag.FromErr(err)
	}

	// Execute the request
	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Parse response
	var assetResp AssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	// Set the resource ID using display_id
	d.SetId(strconv.Itoa(assetResp.Asset.DisplayID))

	return setAssetData(d, &assetResp.Asset)
}

func resourceAssetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset display ID (stored as Terraform resource ID)
	displayID := d.Id()

	// Create the request using display_id
	endpoint := fmt.Sprintf("/assets/%s", displayID)
	req, err := config.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return diag.Errorf("Failed to create request for asset %s: %s", displayID, err)
	}

	// Execute the request
	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.Errorf("Request failed for asset %s: %s", displayID, err)
	}
	defer resp.Body.Close()

	// Check for 404 specifically
	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	// Check for other non-200 status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("API request failed with status %d for asset %s", resp.StatusCode, displayID)
	}

	// Parse response
	var assetResp AssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response for asset %s: %s", displayID, err)
	}

	return setAssetData(d, &assetResp.Asset)
}

func resourceAssetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset display ID (stored as Terraform resource ID)
	displayID := d.Id()

	// Build type_fields from the type_fields map
	typeFields := make(map[string]interface{})
	if typeFieldsRaw, ok := d.GetOk("type_fields"); ok {
		assetTypeID := d.Get("asset_type_id").(int)
		for key, value := range typeFieldsRaw.(map[string]interface{}) {
			// Automatically append the asset type ID to the field name
			fieldKey := fmt.Sprintf("%s_%d", key, assetTypeID)
			// Convert string values to appropriate types based on common patterns
			typeFields[fieldKey] = convertTypeFieldValue(value.(string))
		}
	}

	// Build request body with all current values
	assetReq := AssetRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		AssetTypeID: d.Get("asset_type_id").(int),
		Impact:      d.Get("impact").(string),
		UsageType:   d.Get("usage_type").(string),
		TypeFields:  typeFields,
	}

	// Handle optional nullable fields
	if userID, ok := d.GetOk("user_id"); ok {
		uid := userID.(int)
		assetReq.UserID = &uid
	}
	if locationID, ok := d.GetOk("location_id"); ok {
		lid := locationID.(int)
		assetReq.LocationID = &lid
	}
	if departmentID, ok := d.GetOk("department_id"); ok {
		did := departmentID.(int)
		assetReq.DepartmentID = &did
	}
	if agentID, ok := d.GetOk("agent_id"); ok {
		aid := agentID.(int)
		assetReq.AgentID = &aid
	}
	if groupID, ok := d.GetOk("group_id"); ok {
		gid := groupID.(int)
		assetReq.GroupID = &gid
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(assetReq)
	if err != nil {
		return diag.Errorf("Failed to marshal request: %s", err)
	}

	// Create the request
	endpoint := fmt.Sprintf("/assets/%s", displayID)
	req, err := config.NewRequest(ctx, "PUT", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return diag.FromErr(err)
	}

	// Execute the request
	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Parse response
	var assetResp AssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	return setAssetData(d, &assetResp.Asset)
}

func resourceAssetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset display ID (stored as Terraform resource ID)
	displayID := d.Id()

	// Create the request
	endpoint := fmt.Sprintf("/assets/%s", displayID)
	req, err := config.NewRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Execute the request
	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// If asset not found, it's already deleted
	if resp.StatusCode == 404 {
		d.SetId("")
		return nil
	}

	// Check for other errors
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return diag.Errorf("Failed to delete asset: API returned status %d", resp.StatusCode)
	}

	// Clear the resource ID
	d.SetId("")

	return nil
}

// setAssetData sets the asset data in the Terraform state
func setAssetData(d *schema.ResourceData, asset *Asset) diag.Diagnostics {
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
	if err := d.Set("usage_type", asset.UsageType); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_id", asset.DisplayID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("asset_tag", asset.AssetTag); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("author_type", asset.AuthorType); err != nil {
		return diag.FromErr(err)
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
	if asset.GroupID != nil {
		if err := d.Set("group_id", *asset.GroupID); err != nil {
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

	// Set type_fields - convert back to map[string]string for Terraform
	// Strip the asset type ID suffix from field names
	if asset.TypeFields != nil {
		typeFieldsMap := make(map[string]string)
		assetTypeIDSuffix := fmt.Sprintf("_%d", asset.AssetTypeID)

		for key, value := range asset.TypeFields {
			// Remove the asset type ID suffix if present
			cleanKey := key
			if strings.HasSuffix(key, assetTypeIDSuffix) {
				cleanKey = strings.TrimSuffix(key, assetTypeIDSuffix)
			}
			typeFieldsMap[cleanKey] = fmt.Sprintf("%v", value)
		}
		if err := d.Set("type_fields", typeFieldsMap); err != nil {
			return diag.FromErr(err)
		}
	}

	return nil
}

// convertTypeFieldValue converts string values to appropriate types for the API
func convertTypeFieldValue(value string) interface{} {
	// Try to convert to int
	if intVal, err := strconv.Atoi(value); err == nil {
		return intVal
	}

	// Try to convert to float
	if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
		return floatVal
	}

	// Try to convert to bool
	if boolVal, err := strconv.ParseBool(value); err == nil {
		return boolVal
	}

	// Return as string if no conversion is possible
	return value
}
