---
page_title: "Freshservice Provider"
subcategory: ""
description: |-
  The Freshservice provider is used to interact with Freshservice assets and asset types through the Freshservice API.
---

# Freshservice Provider

The Freshservice provider is used to interact with [Freshservice](https://freshservice.com/) assets and asset types through the Freshservice API. It allows you to manage cloud infrastructure assets like Azure subscriptions, AWS accounts, and GCP projects as assets in Freshservice.

## Example Usage

```terraform
terraform {
  required_providers {
    freshservice = {
      source  = "lcp-llp/freshservice"
      version = "~> 1.0"
    }
  }
}

provider "freshservice" {
  api_key = "your-api-key"
  domain  = "your-domain.freshservice.com"
}

# Create an Azure subscription asset
resource "freshservice_azure_subscription" "example" {
  subscription_name = "Production Subscription"
  subscription_id   = "12345678-1234-5678-9012-123456789012"
  tenant_id        = "87654321-4321-8765-2109-876543210987"
  po_number        = "PO-2024-001"
  owner           = "john.doe@company.com"
  approver        = "jane.smith@company.com"
  environment     = "Production"
  description     = "Main production Azure subscription"
}

# Search for existing assets
data "freshservice_asset" "laptop" {
  name = "Dell laptop"
}
```

## Schema

### Required

- `api_key` (String) Your Freshservice API key
- `domain` (String) Your Freshservice domain (e.g., 'yourdomain.freshservice.com')

## Resources

- [freshservice_asset_type](docs/resources/asset_type.md) - Manage Freshservice asset types
- [freshservice_azure_subscription](docs/resources/azure_subscription.md) - Manage Azure subscription assets
- [freshservice_aws_account](docs/resources/aws_account.md) - Manage AWS account assets  
- [freshservice_gcp_project](docs/resources/gcp_project.md) - Manage GCP project assets

## Data Sources

- [freshservice_asset](docs/data-sources/asset.md) - Search for existing assets
- [freshservice_asset_type](docs/data-sources/asset_type.md) - Retrieve asset type information

## API Rate Limits

Please be aware of Freshservice API rate limits when using this provider. The provider automatically handles authentication and request formatting according to Freshservice API specifications.
