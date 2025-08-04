---
page_title: "freshservice_asset_type Data Source - freshservice"
subcategory: ""
description: |-
  Use this data source to retrieve information about an existing Freshservice asset type.
---

# freshservice_asset_type (Data Source)

Use this data source to retrieve information about an existing Freshservice asset type.

## Example Usage

```terraform
# Get asset type by ID
data "freshservice_asset_type" "laptop" {
  id = 123
}

# Get asset type by name
data "freshservice_asset_type" "by_name" {
  name = "Laptop"
}

# Use the asset type data
resource "freshservice_asset_type" "sub_type" {
  name                 = "Gaming Laptop"
  description          = "High-performance gaming laptops"
  parent_asset_type_id = data.freshservice_asset_type.laptop.id
}

output "asset_type_details" {
  value = {
    id          = data.freshservice_asset_type.laptop.id
    name        = data.freshservice_asset_type.laptop.name
    description = data.freshservice_asset_type.laptop.description
    visible     = data.freshservice_asset_type.laptop.visible
  }
}
```

## Schema

### Optional

One of the following must be provided:

- `id` (Number) ID of the asset type
- `name` (String) Name of the asset type to search for

### Read-Only

- `description` (String) Short description of the asset type
- `parent_asset_type_id` (Number) ID of the parent asset type
- `visible` (Boolean) Visibility of the asset type
- `created_at` (String) Creation timestamp of the asset type
- `updated_at` (String) Last update timestamp of the asset type

## Notes

- If both `id` and `name` are provided, `id` takes precedence.
- The search by name is case-sensitive and must match exactly.
- If no asset type is found with the specified criteria, the data source will return an error.
