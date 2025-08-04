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

// GCPProjectAsset represents a Freshservice GCP Project asset
type GCPProjectAsset struct {
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

// GCPProjectAssetResponse represents the API response for GCP project asset operations
type GCPProjectAssetResponse struct {
	Asset GCPProjectAsset `json:"asset"`
}

// GCPProjectAssetRequest represents the request body for GCP project asset operations
type GCPProjectAssetRequest struct {
	Name        string                 `json:"name"`
	AssetTypeID int                    `json:"asset_type_id"`
	Description string                 `json:"description,omitempty"`
	TypeFields  map[string]interface{} `json:"type_fields"`
}

func resourceGCPProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGCPProjectCreate,
		ReadContext:   resourceGCPProjectRead,
		UpdateContext: resourceGCPProjectUpdate,
		DeleteContext: resourceGCPProjectDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Manages a Freshservice GCP Project asset",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the GCP project asset",
			},
			"project_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the GCP project",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "GCP project ID",
			},
			"po_number": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Purchase order number",
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Owner of the GCP project",
			},
			"approver": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Approver for the GCP project",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Environment type (e.g., Production, Development, Test)",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of the GCP project asset",
			},
			"asset_type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     int64(56000979438),
				Description: "Asset type ID for GCP project (default: 56000979438)",
			},
			"active": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "Yes",
				Description: "Active status (default: Yes)",
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

func resourceGCPProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get asset type ID
	assetTypeID := d.Get("asset_type_id").(int)

	// Build type_fields with the specific field names based on asset type ID
	typeFields := map[string]interface{}{}

	if projectID := d.Get("project_id").(string); projectID != "" {
		typeFields[fmt.Sprintf("project_id_%d", assetTypeID)] = projectID
	}

	if projectName := d.Get("project_name").(string); projectName != "" {
		typeFields[fmt.Sprintf("project_name_%d", assetTypeID)] = projectName
	}

	if poNumber := d.Get("po_number").(string); poNumber != "" {
		typeFields[fmt.Sprintf("po_%d", assetTypeID)] = poNumber
	}

	if owner := d.Get("owner").(string); owner != "" {
		typeFields[fmt.Sprintf("owner_%d", assetTypeID)] = owner
	}

	if approver := d.Get("approver").(string); approver != "" {
		typeFields[fmt.Sprintf("approved_by_%d", assetTypeID)] = approver
	}

	if environment := d.Get("environment").(string); environment != "" {
		typeFields[fmt.Sprintf("environment_%d", assetTypeID)] = environment
	}

	if active := d.Get("active").(string); active != "" {
		typeFields[fmt.Sprintf("active_%d", assetTypeID)] = active
	}

	// Build request body
	assetReq := GCPProjectAssetRequest{
		Name:        d.Get("project_name").(string),
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
	var assetResp GCPProjectAssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	// Set the resource ID and other computed fields
	d.SetId(strconv.Itoa(assetResp.Asset.ID))

	return setGCPProjectAssetData(d, &assetResp.Asset)
}

func resourceGCPProjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset ID
	id := d.Id()

	// Create the request
	endpoint := fmt.Sprintf("/assets/%s", id)
	req, err := config.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Execute the request
	resp, err := config.DoRequest(req)
	if err != nil {
		// If asset not found, remove from state
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Parse response
	var assetResp GCPProjectAssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	return setGCPProjectAssetData(d, &assetResp.Asset)
}

func resourceGCPProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset ID
	id := d.Id()

	// Get asset type ID
	assetTypeID := d.Get("asset_type_id").(int)

	// Build type_fields with updated values
	typeFields := map[string]interface{}{}

	if d.HasChange("project_id") {
		typeFields[fmt.Sprintf("project_id_%d", assetTypeID)] = d.Get("project_id").(string)
	}

	if d.HasChange("project_name") {
		typeFields[fmt.Sprintf("project_name_%d", assetTypeID)] = d.Get("project_name").(string)
	}

	if d.HasChange("po_number") {
		typeFields[fmt.Sprintf("po_%d", assetTypeID)] = d.Get("po_number").(string)
	}

	if d.HasChange("owner") {
		typeFields[fmt.Sprintf("owner_%d", assetTypeID)] = d.Get("owner").(string)
	}

	if d.HasChange("approver") {
		typeFields[fmt.Sprintf("approved_by_%d", assetTypeID)] = d.Get("approver").(string)
	}

	if d.HasChange("environment") {
		typeFields[fmt.Sprintf("environment_%d", assetTypeID)] = d.Get("environment").(string)
	}

	if d.HasChange("active") {
		typeFields[fmt.Sprintf("active_%d", assetTypeID)] = d.Get("active").(string)
	}

	// Build request body with only changed fields
	assetReq := GCPProjectAssetRequest{}

	if d.HasChange("project_name") {
		assetReq.Name = d.Get("project_name").(string)
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
	endpoint := fmt.Sprintf("/assets/%s", id)
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
	var assetResp GCPProjectAssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	return setGCPProjectAssetData(d, &assetResp.Asset)
}

func resourceGCPProjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset ID
	id := d.Id()

	// Create the request
	endpoint := fmt.Sprintf("/assets/%s", id)
	req, err := config.NewRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Execute the request
	resp, err := config.DoRequest(req)
	if err != nil {
		// If asset not found, it's already deleted
		if resp != nil && resp.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Clear the resource ID
	d.SetId("")

	return nil
}

// setGCPProjectAssetData sets the GCP project asset data in the Terraform state
func setGCPProjectAssetData(d *schema.ResourceData, asset *GCPProjectAsset) diag.Diagnostics {
	if err := d.Set("project_name", asset.Name); err != nil {
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
		if projectID, ok := asset.TypeFields[fmt.Sprintf("project_id_%d", assetTypeID)].(string); ok {
			if err := d.Set("project_id", projectID); err != nil {
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

		if approver, ok := asset.TypeFields[fmt.Sprintf("approved_by_%d", assetTypeID)].(string); ok {
			if err := d.Set("approver", approver); err != nil {
				return diag.FromErr(err)
			}
		}

		if environment, ok := asset.TypeFields[fmt.Sprintf("environment_%d", assetTypeID)].(string); ok {
			if err := d.Set("environment", environment); err != nil {
				return diag.FromErr(err)
			}
		}

		if active, ok := asset.TypeFields[fmt.Sprintf("active_%d", assetTypeID)].(string); ok {
			if err := d.Set("active", active); err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return nil
}
