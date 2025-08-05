---
page_title: "freshservice_asset Resource - freshservice"
subcategory: ""
description: |-
  Manages a Freshservice asset with custom type fields
---

# freshservice_asset (Resource)

Manages a Freshservice asset with custom type fields. This resource allows you to create and manage assets of any type in Freshservice with their specific type fields.

## Example Usage

```terraform
# Create a basic asset
resource "freshservice_asset" "laptop" {
  name          = "Dell Laptop"
  description   = "Dell Latitude 5520 laptop"
  asset_type_id = 25
  impact        = "medium"
  usage_type    = "permanent"
  
  type_fields = {
    "product"               = "10"
    "vendor"                = "14"
    "cost"                  = "5000"
    "salvage"               = "100"
    "depreciation_id"       = "30"
    "warranty"              = "20"
    "acquisition_date"      = "2023-07-26T12:25:04+05:30"
    "warranty_expiry_date"  = "2026-07-26T12:25:04+05:30"
    "domain"                = "1"
    "asset_state"           = "In Use"
    "serial_number"         = "SW12131133"
    "last_audit_date"       = "2023-07-26T12:25:04+05:30"
  }
}

# Create an asset with assignment
resource "freshservice_asset" "server" {
  name          = "Production Server"
  description   = "Main production web server"
  asset_type_id = 30
  impact        = "high"
  usage_type    = "permanent"
  user_id       = 1001
  location_id   = 5
  department_id = 10
  agent_id      = 2001
  group_id      = 15
  
  type_fields = {
    "hostname"     = "prod-web-01"
    "ip_address"   = "192.168.1.100"
    "os_version"   = "Ubuntu 22.04"
    "cpu_cores"    = "8"
    "memory_gb"    = "32"
    "storage_gb"   = "1000"
  }
}

# Output asset information
output "asset_details" {
  value = {
    id         = freshservice_asset.laptop.id
    display_id = freshservice_asset.laptop.display_id
    asset_tag  = freshservice_asset.laptop.asset_tag
    name       = freshservice_asset.laptop.name
  }
}
```

## Schema

### Required

- `name` (String) Name of the asset
- `asset_type_id` (Number) Asset type ID. Cannot be changed after creation.

### Optional

- `description` (String) Description of the asset
- `impact` (String) Impact level of the asset (low, medium, high). Default: "low"
- `usage_type` (String) Usage type of the asset (permanent, loaner). Default: "permanent"
- `user_id` (Number) User ID assigned to the asset
- `location_id` (Number) Location ID of the asset
- `department_id` (Number) Department ID of the asset
- `agent_id` (Number) Agent ID assigned to the asset
- `group_id` (Number) Group ID assigned to the asset
- `type_fields` (Map of String) Custom type fields specific to the asset type. Field names will automatically have the asset type ID appended (e.g., 'product' becomes 'product_25')

### Read-Only

- `id` (String) ID of the asset (contains display_id value)
- `display_id` (Number) Display ID of the asset
- `asset_tag` (String) Asset tag
- `author_type` (String) Author type of the asset
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
- `item_id` (String) Item ID of the asset
- `imei_number` (String) IMEI number of the asset

## Import

Assets can be imported using their display_id:

```bash
terraform import freshservice_asset.laptop 3567
```

## Notes

### Type Fields

The `type_fields` map allows you to set custom fields specific to each asset type. You only need to specify the field name, and the provider will automatically append the asset type ID.

For example, if you have an asset type with ID `25` and specify a field named `product`, the provider will automatically convert it to `product_25` when sending to the Freshservice API.

**Example:**
```terraform
# You specify:
type_fields = {
  "product" = "Dell Laptop"
  "vendor"  = "14"
}

# Provider automatically converts to:
# product_25 = "Dell Laptop"
# vendor_25 = "14"
```

### Supported Field Types

Type fields support various data types:
- **Strings**: Text values (e.g., "Dell Laptop")
- **Numbers**: Integer and decimal values (e.g., "5000", "32.5")
- **Booleans**: true/false values (e.g., "true", "false")
- **Dates**: ISO 8601 formatted dates (e.g., "2023-07-26T12:25:04+05:30")

The provider automatically converts string values to the appropriate type when sending to the API.

### Asset Type Restrictions

The `asset_type_id` cannot be changed after the asset is created. If you need to change the asset type, you must destroy and recreate the resource.

### Display ID vs Internal ID

Freshservice uses both internal IDs and display IDs for assets:
- The `id` field contains the display_id value (used for API operations)
- The `display_id` field shows the same value as a number
- Import operations use the display_id

This ensures proper state management and avoids 404 errors during resource operations.
