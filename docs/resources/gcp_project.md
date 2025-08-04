---
page_title: "freshservice_gcp_project Resource - freshservice"
subcategory: ""
description: |-
  Manages a GCP project asset in Freshservice with custom type fields for tracking Google Cloud Platform projects.
---

# freshservice_gcp_project (Resource)

Manages a GCP project asset in Freshservice with custom type fields for tracking Google Cloud Platform projects.

## Example Usage

```terraform
resource "freshservice_gcp_project" "production" {
  project_name = "my-production-project"
  project_id   = "my-prod-project-123456"
  po_number    = "PO-2024-003"
  owner        = "gcp.admin@company.com"
  approver     = "project.manager@company.com"
  environment  = "Production"
  description  = "Main production GCP project"
  active       = "Yes"
}
```

## Schema

### Required

- `project_name` (String) Name of the GCP project
- `project_id` (String) GCP project ID

### Optional

- `po_number` (String) Purchase order number
- `owner` (String) Owner of the GCP project
- `approver` (String) Approver for the GCP project
- `environment` (String) Environment type (e.g., Production, Development, Test)
- `description` (String) Description of the GCP project asset
- `asset_type_id` (Number) Asset type ID for GCP project (default: 56000979438)
- `active` (String) Active status (default: "Yes")

### Read-Only

- `id` (String) ID of the GCP project asset
- `display_id` (Number) Display ID of the asset
- `asset_tag` (String) Asset tag
- `created_at` (String) Creation timestamp of the asset
- `updated_at` (String) Last update timestamp of the asset
- `workspace_id` (Number) Workspace ID of the asset

## Import

Import is supported using the following syntax:

```shell
terraform import freshservice_gcp_project.example 123
```
