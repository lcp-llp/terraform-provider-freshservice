---
page_title: "freshservice_asset_type Resource - freshservice"
subcategory: ""
description: |-
  Manages a Freshservice asset type. Asset types are used to categorize and organize assets within Freshservice.
---

# freshservice_asset_type (Resource)

Manages a Freshservice asset type. Asset types are used to categorize and organize assets within Freshservice.

## Example Usage

```terraform
# Create a basic asset type
resource "freshservice_asset_type" "laptop" {
  name        = "Laptop"
  description = "Corporate laptops and portable computers"
}

# Create an asset type with a parent
resource "freshservice_asset_type" "macbook" {
  name                 = "MacBook"
  description          = "Apple MacBook laptops"
  parent_asset_type_id = freshservice_asset_type.laptop.id
  visible              = true
}
```

## Schema

### Required

- `name` (String) Name of the asset type

### Optional

- `description` (String) Short description of the asset type
- `parent_asset_type_id` (Number) ID of the parent asset type
- `visible` (Boolean) Visibility of the asset type. Custom asset types are set to true by default

### Read-Only

- `id` (String) ID of the asset type
- `created_at` (String) Creation timestamp of the asset type
- `updated_at` (String) Last update timestamp of the asset type

## Import

Import is supported using the following syntax:

```shell
terraform import freshservice_asset_type.example 123
```
