---
page_title: "freshservice_asset Data Source - freshservice"
subcategory: ""
description: |-
  Use this data source to search for existing Freshservice assets by name, item_id, or asset_tag.
---

# freshservice_asset (Data Source)

Use this data source to search for existing Freshservice assets by name, item_id, or asset_tag.

## Example Usage

```terraform
# Search for an asset by name
data "freshservice_asset" "laptop" {
  name = "Dell laptop"
}

# Search for an asset by item_id
data "freshservice_asset" "by_item_id" {
  item_id = "ITEM-12345"
}

# Search for an asset by asset_tag
data "freshservice_asset" "by_tag" {
  asset_tag = "ASSET-67890"
}

# Search in trash
data "freshservice_asset" "trashed_asset" {
  name    = "Old laptop"
  trashed = true
}

# Use the asset data
output "asset_details" {
  value = {
    id          = data.freshservice_asset.laptop.id
    name        = data.freshservice_asset.laptop.name
    asset_tag   = data.freshservice_asset.laptop.asset_tag
    created_at  = data.freshservice_asset.laptop.created_at
  }
}
```

## Schema

### Optional

At least one of the following search parameters must be provided:

- `name` (String) Name of the asset to search for
- `item_id` (String) Item ID of the asset to search for
- `asset_tag` (String) Asset tag to search for
- `trashed` (Boolean) Include assets in trash (default: false)

### Read-Only

- `id` (String) Display ID of the asset (used for identification and API calls)
- `display_id` (Number) Display ID of the asset (same as id but as number)
- `description` (String) Description of the asset
- `asset_type_id` (Number) Asset type ID
- `impact` (String) Impact level of the asset
- `author_type` (String) Author type of the asset
- `usage_type` (String) Usage type of the asset
- `user_id` (Number) User ID assigned to the asset
- `location_id` (Number) Location ID of the asset
- `department_id` (Number) Department ID of the asset
- `agent_id` (Number) Agent ID assigned to the asset
- `assigned_on` (String) Date when the asset was assigned
- `created_at` (String) Creation timestamp of the asset
- `updated_at` (String) Last update timestamp of the asset
- `workspace_id` (Number) Workspace ID of the asset
- `created_by_source` (String) Source that created the asset
- `last_updated_by_source` (String) Source that last updated the asset
- `created_by_user` (Number) User who created the asset
- `last_updated_by_user` (Number) User who last updated the asset
- `sources` (List of String) List of sources for the asset
- `serial_number` (String) Serial number of the asset
- `mac_addresses` (List of String) MAC addresses of the asset
- `ip_addresses` (List of String) IP addresses of the asset
- `uuid` (String) UUID of the asset
- `imei_number` (String) IMEI number of the asset

## Notes

- The search must return exactly one asset. If multiple assets match the criteria, the data source will return an error.
- If no assets are found, the data source will return an error.
- Search results cannot be sorted and are returned by default sorted by created_at in descending order.
- Search queries are case-insensitive and support partial matching.
