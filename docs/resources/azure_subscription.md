---
page_title: "freshservice_azure_subscription Resource - freshservice"
subcategory: ""
description: |-
  Manages an Azure subscription asset in Freshservice with custom type fields for tracking Azure subscriptions.
---

# freshservice_azure_subscription (Resource)

Manages an Azure subscription asset in Freshservice with custom type fields for tracking Azure subscriptions.

## Example Usage

```terraform
resource "freshservice_azure_subscription" "production" {
  subscription_name = "Production Subscription"
  subscription_id   = "12345678-1234-5678-9012-123456789012"
  tenant_id        = "87654321-4321-8765-2109-876543210987"
  po_number        = "PO-2024-001"
  owner           = "john.doe@company.com"
  approver        = "jane.smith@company.com"
  environment     = "Production"
  description     = "Main production Azure subscription"
  eacsp           = "CSP"
  active          = "Yes"
  cloudockit      = "Yes"
}
```

## Schema

### Required

- `subscription_name` (String) Name of the Azure subscription
- `subscription_id` (String) Azure subscription ID
- `tenant_id` (String) Azure tenant ID

### Optional

- `po_number` (String) Purchase order number
- `owner` (String) Owner of the Azure subscription
- `approver` (String) Approver for the Azure subscription
- `environment` (String) Environment type (e.g., Production, Development, Test)
- `description` (String) Description of the Azure subscription asset
- `asset_type_id` (Number) Asset type ID for Azure subscription (default: 56000416566)
- `eacsp` (String) EA/CSP field (default: "CSP")
- `active` (String) Active status (default: "Yes")
- `cloudockit` (String) Cloudockit field (default: "Yes")

### Read-Only

- `id` (String) ID of the Azure subscription asset
- `display_id` (Number) Display ID of the asset
- `asset_tag` (String) Asset tag
- `created_at` (String) Creation timestamp of the asset
- `updated_at` (String) Last update timestamp of the asset
- `workspace_id` (Number) Workspace ID of the asset

## Import

Import is supported using the following syntax:

```shell
terraform import freshservice_azure_subscription.example 123
```
