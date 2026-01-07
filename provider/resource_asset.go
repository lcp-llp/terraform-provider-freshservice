package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AssetResponse represents the API response for asset operations
type AssetResponse struct {
	Asset Asset `json:"asset"`
}

type Asset struct {
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	DisplayID   int                    `json:"display_id"`
	AssetTypeID int                    `json:"asset_type_id"`
	Description string                 `json:"description,omitempty"`
	AssetTag    string                 `json:"asset_tag"`
	Impact      string                 `json:"impact"`
	TypeFields  map[string]interface{} `json:"type_fields"`
}

// AssetRequest represents the request body for asset operations
type AssetRequest struct {
	Name        string                 `json:"name"`
	AssetTypeID int                    `json:"asset_type_id"`
	Description *string                `json:"description,omitempty"`
	AssetTag    *string                `json:"asset_tag,omitempty"`
	Impact      *string                `json:"impact,omitempty"`
	TypeFields  map[string]interface{} `json:"type_fields,omitempty"`
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
				Description: "ID of the asset",
			},
			"display_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Display ID of the asset",
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
			"asset_tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Asset tag",
			},
			"impact": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Impact of the asset",
			},
			"type_fields": {
				Type:        schema.TypeMap,
				Elem:        schema.TypeString,
				Optional:    true,
				Description: "Custom type fields for the asset",
			},
		},
	}
}

func resourceAssetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	assetReq := &AssetRequest{
		Name:        d.Get("name").(string),
		AssetTypeID: d.Get("asset_type_id").(int),
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		assetReq.Description = &desc
	}
	if v, ok := d.GetOk("asset_tag"); ok {
		tag := v.(string)
		assetReq.AssetTag = &tag
	}
	if v, ok := d.GetOk("impact"); ok {
		imp := v.(string)
		assetReq.Impact = &imp
	}
	if v, ok := d.GetOk("type_fields"); ok {
		assetReq.TypeFields = expandTypeFields(v.(map[string]interface{}))
	}

	jsonData, err := json.Marshal(assetReq)
	if err != nil {
		return diag.FromErr(err)
	}

	req, err := config.NewRequest(ctx, "POST", "/assets", bytes.NewReader(jsonData))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	var assetResp AssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(assetResp.Asset.DisplayID))
	return setAssetData(d, &assetResp.Asset)
}

func resourceAssetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	displayID := d.Id()
	endpoint := fmt.Sprintf("/assets/%s", displayID)
	req, err := config.NewRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := config.DoRequest(req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	var assetResp AssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(assetResp.Asset.DisplayID))
	return setAssetData(d, &assetResp.Asset)
}

func resourceAssetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	displayID := d.Id()
	assetReq := &AssetRequest{
		Name:        d.Get("name").(string),
		AssetTypeID: d.Get("asset_type_id").(int),
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		assetReq.Description = &desc
	}
	if v, ok := d.GetOk("asset_tag"); ok {
		tag := v.(string)
		assetReq.AssetTag = &tag
	}
	if v, ok := d.GetOk("impact"); ok {
		imp := v.(string)
		assetReq.Impact = &imp
	}
	if v, ok := d.GetOk("type_fields"); ok {
		assetReq.TypeFields = expandTypeFields(v.(map[string]interface{}))
	}

	jsonData, err := json.Marshal(assetReq)
	if err != nil {
		return diag.FromErr(err)
	}

	endpoint := fmt.Sprintf("/assets/%s", displayID)
	req, err := config.NewRequest(ctx, "PUT", endpoint, bytes.NewReader(jsonData))
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := config.DoRequest(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	var assetResp AssetResponse
	if err := json.NewDecoder(resp.Body).Decode(&assetResp); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(assetResp.Asset.DisplayID))
	return setAssetData(d, &assetResp.Asset)
}

func resourceAssetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)

	displayID := d.Id()
	endpoint := fmt.Sprintf("/assets/%s", displayID)
	req, err := config.NewRequest(ctx, "DELETE", endpoint, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := config.DoRequest(req)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	d.SetId("")
	return nil
}

func setAssetData(d *schema.ResourceData, asset *Asset) diag.Diagnostics {
	var diags diag.Diagnostics
	if err := d.Set("id", strconv.Itoa(asset.ID)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", asset.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("display_id", asset.DisplayID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", asset.Description); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("asset_type_id", asset.AssetTypeID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("asset_tag", asset.AssetTag); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("impact", asset.Impact); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type_fields", flattenTypeFields(asset.TypeFields)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func expandTypeFields(input map[string]interface{}) map[string]interface{} {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]interface{}, len(input))
	for key, value := range input {
		if value == nil {
			continue
		}
		result[key] = value
	}
	return result
}

func flattenTypeFields(fields map[string]interface{}) map[string]interface{} {
	if len(fields) == 0 {
		return nil
	}
	result := make(map[string]interface{}, len(fields))
	for key, value := range fields {
		result[key] = fmt.Sprint(value)
	}
	return result
}
