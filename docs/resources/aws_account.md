---
page_title: "freshservice_aws_account Resource - freshservice"
subcategory: ""
description: |-
  Manages an AWS account asset in Freshservice with custom type fields for tracking AWS accounts.
---

# freshservice_aws_account (Resource)

Manages an AWS account asset in Freshservice with custom type fields for tracking AWS accounts.

## Example Usage

```terraform
resource "freshservice_aws_account" "production" {
  account_name = "Production AWS Account"
  account_id   = "123456789012"
  po_number    = "PO-2024-002"
  owner        = "aws.admin@company.com"
  approver     = "finance@company.com"
  environment  = "Production"
  description  = "Main production AWS account"
}
```

## Schema

### Required

- `account_name` (String) Name of the AWS account
- `account_id` (String) AWS account ID

### Optional

- `po_number` (String) Purchase order number
- `owner` (String) Owner of the AWS account
- `approver` (String) Approver for the AWS account
- `environment` (String) Environment type (e.g., Production, Development, Test)
- `description` (String) Description of the AWS account asset
- `asset_type_id` (Number) Asset type ID for AWS account (default: 56000947175)

### Read-Only

- `id` (String) Display ID of the AWS account asset (used for API calls)
- `display_id` (Number) Display ID of the asset (same as id but as number)
- `asset_tag` (String) Asset tag
- `created_at` (String) Creation timestamp of the asset
- `updated_at` (String) Last update timestamp of the asset
- `workspace_id` (Number) Workspace ID of the asset

## Import

Import is supported using the display ID:

```shell
terraform import freshservice_aws_account.example 3567
```
