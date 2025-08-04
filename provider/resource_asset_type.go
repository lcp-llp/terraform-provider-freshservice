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

// AssetType represents a Freshservice asset type
type AssetType struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	ParentAssetTypeID *int      `json:"parent_asset_type_id"`
	Description       string    `json:"description"`
	Visible           bool      `json:"visible"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

// AssetTypeResponse represents the API response for asset type operations
type AssetTypeResponse struct {
	AssetType AssetType `json:"asset_type"`
}

// AssetTypeRequest represents the request body for asset type operations
type AssetTypeRequest struct {
	Name              string `json:"name,omitempty"`
	ParentAssetTypeID *int   `json:"parent_asset_type_id,omitempty"`
	Description       string `json:"description,omitempty"`
	Visible           *bool  `json:"visible,omitempty"`
}

func resourceAssetType() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetTypeCreate,
		ReadContext:   resourceAssetTypeRead,
		UpdateContext: resourceAssetTypeUpdate,
		DeleteContext: resourceAssetTypeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "Manages a Freshservice asset type",

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the asset type",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the asset type",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Short description of the asset type",
			},
			"parent_asset_type_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the parent asset type",
			},
			"visible": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Visibility of the asset type. Custom asset types are set to true by default",
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

func resourceAssetTypeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Build request body
	assetTypeReq := AssetTypeRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if parentID, ok := d.GetOk("parent_asset_type_id"); ok {
		parentAssetTypeID := parentID.(int)
		assetTypeReq.ParentAssetTypeID = &parentAssetTypeID
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(assetTypeReq)
	if err != nil {
		return diag.Errorf("Failed to marshal request: %s", err)
	}

	// Create the request
	req, err := config.NewRequest(ctx, "POST", "/asset_types", bytes.NewReader(jsonData))
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
	var assetTypeResp AssetTypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetTypeResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	// Set the resource ID and other computed fields
	d.SetId(strconv.Itoa(assetTypeResp.AssetType.ID))

	return setAssetTypeData(d, &assetTypeResp.AssetType)
}

func resourceAssetTypeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset type ID
	id := d.Id()

	// Create the request
	endpoint := fmt.Sprintf("/asset_types/%s", id)
	req, err := config.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Execute the request
	resp, err := config.DoRequest(req)
	if err != nil {
		// If asset type not found, remove from state
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Parse response
	var assetTypeResp AssetTypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetTypeResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	return setAssetTypeData(d, &assetTypeResp.AssetType)
}

func resourceAssetTypeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset type ID
	id := d.Id()

	// Build request body with only changed fields
	assetTypeReq := AssetTypeRequest{}

	if d.HasChange("name") {
		assetTypeReq.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		assetTypeReq.Description = d.Get("description").(string)
	}

	if d.HasChange("parent_asset_type_id") {
		if parentID, ok := d.GetOk("parent_asset_type_id"); ok {
			parentAssetTypeID := parentID.(int)
			assetTypeReq.ParentAssetTypeID = &parentAssetTypeID
		}
	}

	if d.HasChange("visible") {
		visible := d.Get("visible").(bool)
		assetTypeReq.Visible = &visible
	}

	// Convert request to JSON
	jsonData, err := json.Marshal(assetTypeReq)
	if err != nil {
		return diag.Errorf("Failed to marshal request: %s", err)
	}

	// Create the request
	endpoint := fmt.Sprintf("/asset_types/%s", id)
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
	var assetTypeResp AssetTypeResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetTypeResp); err != nil {
		return diag.Errorf("Failed to decode response: %s", err)
	}

	return setAssetTypeData(d, &assetTypeResp.AssetType)
}

func resourceAssetTypeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	// Get the asset type ID
	id := d.Id()

	// Create the request
	endpoint := fmt.Sprintf("/asset_types/%s", id)
	req, err := config.NewRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Execute the request
	resp, err := config.DoRequest(req)
	if err != nil {
		// If asset type not found, it's already deleted
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

// setAssetTypeData sets the asset type data in the Terraform state
func setAssetTypeData(d *schema.ResourceData, assetType *AssetType) diag.Diagnostics {
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
