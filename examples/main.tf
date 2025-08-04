terraform {
  required_providers {
    freshservice = {
      source  = "lcp-llp/freshservice"
      version = "~> 0.0.1"
    }
  }
}

provider "freshservice" {
  api_key = var.freshservice_api_key
  domain  = var.freshservice_domain
}

variable "freshservice_api_key" {
  description = "Freshservice API key"
  type        = string
  sensitive   = true
}

variable "freshservice_domain" {
  description = "Freshservice domain (e.g., 'yourdomain.freshservice.com')"
  type        = string
}

# Create a custom asset type
resource "freshservice_asset_type" "cloud_subscription" {
  name        = "Cloud Subscription"
  description = "Cloud provider subscriptions and accounts"
}

# Create an Azure subscription asset
resource "freshservice_azure_subscription" "production" {
  subscription_name = "Production Subscription"
  subscription_id   = "12345678-1234-5678-9012-123456789012"
  tenant_id        = "87654321-4321-8765-2109-876543210987"
  po_number        = "PO-2024-001"
  owner           = "john.doe@company.com"
  approver        = "jane.smith@company.com"
  environment     = "Production"
  description     = "Main production Azure subscription"
}

# Create an AWS account asset
resource "freshservice_aws_account" "production" {
  account_name = "Production AWS Account"
  account_id   = "123456789012"
  po_number    = "PO-2024-002"
  owner        = "aws.admin@company.com"
  approver     = "finance@company.com"
  environment  = "Production"
  description  = "Main production AWS account"
}

# Create a GCP project asset
resource "freshservice_gcp_project" "production" {
  project_name = "my-production-project"
  project_id   = "my-prod-project-123456"
  po_number    = "PO-2024-003"
  owner        = "gcp.admin@company.com"
  approver     = "project.manager@company.com"
  environment  = "Production"
  description  = "Main production GCP project"
}

# Search for existing assets
data "freshservice_asset" "existing_laptop" {
  name = "Dell laptop"
}

# Output asset information
output "azure_subscription_id" {
  description = "The ID of the created Azure subscription asset"
  value       = freshservice_azure_subscription.production.id
}

output "aws_account_id" {
  description = "The ID of the created AWS account asset"
  value       = freshservice_aws_account.production.id
}

output "gcp_project_id" {
  description = "The ID of the created GCP project asset"
  value       = freshservice_gcp_project.production.id
}

output "existing_laptop_details" {
  description = "Details of the existing laptop asset"
  value = {
    id         = data.freshservice_asset.existing_laptop.id
    name       = data.freshservice_asset.existing_laptop.name
    asset_tag  = data.freshservice_asset.existing_laptop.asset_tag
    created_at = data.freshservice_asset.existing_laptop.created_at
  }
}
