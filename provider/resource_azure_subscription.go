package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AzureSubscriptionAsset represents a Freshservice Azure Subscription asset
type AzureSubscriptionAsset struct {
	ID           int                    `json:"id"`
	DisplayID    int                    `json:"display_id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	AssetTypeID  int                    `json:"asset_type_id"`
	Impact       string                 `json:"impact"`
	AuthorType   string                 `json:"author_type"`
	UsageType    string                 `json:"usage_type"`
	AssetTag     string                 `json:"asset_tag"`
	UserID       *int                   `json:"user_id"`
	LocationID   *int                   `json:"location_id"`
	DepartmentID *int                   `json:"department_id"`
	AgentID      *int                   `json:"agent_id"`
	AssignedOn   *string                `json:"assigned_on"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	WorkspaceID  int                    `json:"workspace_id"`
	TypeFields   map[string]interface{} `json:"type_fields"`
}

// AzureSubscriptionAssetResponse represents the API response for Azure subscription asset operations
type AzureSubscriptionAssetResponse struct {
	Asset AzureSubscriptionAsset `json:"asset"`
}

// AzureSubscriptionAssetRequest represents the request body for Azure subscription asset operations
type AzureSubscriptionAssetRequest struct {
	Name        string                 `json:"name"`
	AssetTypeID int                    `json:"asset_type_id"`
	Description string                 `json:"description,omitempty"`
	TypeFields  map[string]interface{} `json:"type_fields"`
}

func resourceAzureSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAzureSubscriptionCreate,
		ReadContext:   resourceAzureSubscriptionRead,
		UpdateContext: resourceAzureSubscriptionUpdate,
		DeleteContext: resourceAzureSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Manages a Freshservice Azure Subscription asset",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Azure subscription asset",
			},
			"subscription_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Azure subscription",
			},
			"subscription_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Azure subscription ID",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Azure tenant ID",
			},
			"po_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Purchase order number",
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Owner of the Azure subscription",
			},
			"approver": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Approver for the Azure subscription",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Environment type (e.g., Production, Development, Test)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the Azure subscription asset",
			},
			"asset_type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     int64(56000416566),
				Description: "Asset type ID for Azure subscription (default: 56000416566)",
			},
			"eacsp": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "CSP",
				Description: "EA/CSP field (default: CSP)",
			},
			"active": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Yes",
				Description: "Active status (default: Yes)",
			},
			"cloudockit": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Yes",
				Description: "Cloudockit field (default: Yes)",
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
		},
	}
}

func resourceAzureSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get asset type ID
	assetTypeID := d.Get("asset_type_id").(int)

	// Build type_fields with the specific field names based on asset type ID
	typeFields := map[string]interface{}{}

	if tenantID := d.Get("tenant_id").(string); tenantID != "" {
		typeFields[fmt.Sprintf("tenant_id_%d", assetTypeID)] = tenantID
	}

	if subscriptionID := d.Get("subscription_id").(string); subscriptionID != "" {
		typeFields[fmt.Sprintf("subscription_id_%d", assetTypeID)] = subscriptionID
	}

	if poNumber := d.Get("po_number").(string); poNumber != "" {
		typeFields[fmt.Sprintf("po_%d", assetTypeID)] = poNumber
	}

	if owner := d.Get("owner").(string); owner != "" {
		typeFields[fmt.Sprintf("owner_%d", assetTypeID)] = owner
	}

	if approver := d.Get("approver").(string); approver != "" {
		typeFields[fmt.Sprintf("approver_object_%d", assetTypeID)] = approver
	}

	if environment := d.Get("environment").(string); environment != "" {
		typeFields[fmt.Sprintf("environment_%d", assetTypeID)] = environment
	}

	if eacsp := d.Get("eacsp").(string); eacsp != "" {
		typeFields[fmt.Sprintf("eacsp_%d", assetTypeID)] = eacsp
	}

	if active := d.Get("active").(string); active != "" {
		typeFields[fmt.Sprintf("active_%d", assetTypeID)] = active
	}

	if cloudockit := d.Get("cloudockit").(string); cloudockit != "" {
		typeFields[fmt.Sprintf("cloudockit_%d", assetTypeID)] = cloudockit
	}

	// Build request body
	assetReq := AzureSubscriptionAssetRequest{
		Name:        d.Get("subscription_name").(string),
		AssetTypeID: assetTypeID,
		Description: d.Get("description").(string),
		TypeFields:  typeFields,
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
	var assetResp AzureSubscriptionAssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	// Set the resource ID using display_id (which is used for API calls) and other computed fields
	d.SetId(strconv.Itoa(assetResp.Asset.DisplayID))

	return setAzureSubscriptionAssetData(d, &assetResp.Asset)
}

func resourceAzureSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	var assetResp AzureSubscriptionAssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response for asset %s: %s", displayID, err)
	}

	return setAzureSubscriptionAssetData(d, &assetResp.Asset)
}

func resourceAzureSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset display ID (stored as Terraform resource ID)
	displayID := d.Id()

	// Get asset type ID
	assetTypeID := d.Get("asset_type_id").(int)

	// Build type_fields with updated values
	typeFields := map[string]interface{}{}

	if d.HasChange("tenant_id") {
		typeFields[fmt.Sprintf("tenant_id_%d", assetTypeID)] = d.Get("tenant_id").(string)
	}

	if d.HasChange("subscription_id") {
		typeFields[fmt.Sprintf("subscription_id_%d", assetTypeID)] = d.Get("subscription_id").(string)
	}

	if d.HasChange("po_number") {
		typeFields[fmt.Sprintf("po_%d", assetTypeID)] = d.Get("po_number").(string)
	}

	if d.HasChange("owner") {
		typeFields[fmt.Sprintf("owner_%d", assetTypeID)] = d.Get("owner").(string)
	}

	if d.HasChange("approver") {
		typeFields[fmt.Sprintf("approver_object_%d", assetTypeID)] = d.Get("approver").(string)
	}

	if d.HasChange("environment") {
		typeFields[fmt.Sprintf("environment_%d", assetTypeID)] = d.Get("environment").(string)
	}

	if d.HasChange("eacsp") {
		typeFields[fmt.Sprintf("eacsp_%d", assetTypeID)] = d.Get("eacsp").(string)
	}

	if d.HasChange("active") {
		typeFields[fmt.Sprintf("active_%d", assetTypeID)] = d.Get("active").(string)
	}

	if d.HasChange("cloudockit") {
		typeFields[fmt.Sprintf("cloudockit_%d", assetTypeID)] = d.Get("cloudockit").(string)
	}

	// Build request body with only changed fields
	assetReq := AzureSubscriptionAssetRequest{}

	if d.HasChange("subscription_name") {
		assetReq.Name = d.Get("subscription_name").(string)
	}

	if d.HasChange("description") {
		assetReq.Description = d.Get("description").(string)
	}

	if len(typeFields) > 0 {
		assetReq.TypeFields = typeFields
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
	var assetResp AzureSubscriptionAssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	return setAzureSubscriptionAssetData(d, &assetResp.Asset)
}

func resourceAzureSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

// setAzureSubscriptionAssetData sets the Azure subscription asset data in the Terraform state
func setAzureSubscriptionAssetData(d *schema.ResourceData, asset *AzureSubscriptionAsset) diag.Diagnostics {
	if err := d.Set("subscription_name", asset.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", asset.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("asset_type_id", asset.AssetTypeID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_id", asset.DisplayID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("asset_tag", asset.AssetTag); err != nil {
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

	// Extract values from type_fields
	assetTypeID := asset.AssetTypeID
	if asset.TypeFields != nil {
		if tenantID, ok := asset.TypeFields[fmt.Sprintf("tenant_id_%d", assetTypeID)].(string); ok {
			if err := d.Set("tenant_id", tenantID); err != nil {
				return diag.FromErr(err)
			}
		}

		if subscriptionID, ok := asset.TypeFields[fmt.Sprintf("subscription_id_%d", assetTypeID)].(string); ok {
			if err := d.Set("subscription_id", subscriptionID); err != nil {
				return diag.FromErr(err)
			}
		}

		if poNumber, ok := asset.TypeFields[fmt.Sprintf("po_%d", assetTypeID)].(string); ok {
			if err := d.Set("po_number", poNumber); err != nil {
				return diag.FromErr(err)
			}
		}

		if owner, ok := asset.TypeFields[fmt.Sprintf("owner_%d", assetTypeID)].(string); ok {
			if err := d.Set("owner", owner); err != nil {
				return diag.FromErr(err)
			}
		}

		if approver, ok := asset.TypeFields[fmt.Sprintf("approver_object_%d", assetTypeID)].(string); ok {
			if err := d.Set("approver", approver); err != nil {
				return diag.FromErr(err)
			}
		}

		if environment, ok := asset.TypeFields[fmt.Sprintf("environment_%d", assetTypeID)].(string); ok {
			if err := d.Set("environment", environment); err != nil {
				return diag.FromErr(err)
			}
		}

		if eacsp, ok := asset.TypeFields[fmt.Sprintf("eacsp_%d", assetTypeID)].(string); ok {
			if err := d.Set("eacsp", eacsp); err != nil {
				return diag.FromErr(err)
			}
		}

		if active, ok := asset.TypeFields[fmt.Sprintf("active_%d", assetTypeID)].(string); ok {
			if err := d.Set("active", active); err != nil {
				return diag.FromErr(err)
			}
		}

		if cloudockit, ok := asset.TypeFields[fmt.Sprintf("cloudockit_%d", assetTypeID)].(string); ok {
			if err := d.Set("cloudockit", cloudockit); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}
