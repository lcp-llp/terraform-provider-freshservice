# Terraform Provider for Freshservice

A Terraform provider for managing assets and asset types in [Freshservice](https://freshservice.com/). This provider allows you to manage cloud infrastructure assets like Azure subscriptions, AWS accounts, and GCP projects as assets in Freshservice.

## Features

- **Asset Management**: Create, read, update, and delete assets in Freshservice
- **Asset Type Management**: Manage custom asset types
- **Cloud Asset Support**: Specialized resources for Azure subscriptions, AWS accounts, and GCP projects
- **Search Capabilities**: Find existing assets by name, item_id, or asset_tag
- **Custom Type Fields**: Support for asset-specific custom fields

## Supported Resources

- `freshservice_asset_type` - Manage Freshservice asset types
- `freshservice_azure_subscription` - Manage Azure subscription assets
- `freshservice_aws_account` - Manage AWS account assets
- `freshservice_gcp_project` - Manage GCP project assets

## Supported Data Sources

- `freshservice_asset` - Search for existing assets
- `freshservice_asset_type` - Retrieve asset type information

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19 (for development)

## Using the Provider

### Authentication

The provider uses HTTP Basic Authentication with your Freshservice API key:

1. Log in to your Freshservice account
2. Go to Admin → API Settings
3. Generate or copy your API key

### Example Usage

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
}

# Search for existing assets
data "freshservice_asset" "laptop" {
  name = "Dell laptop"
}
```

## Development

### Building the Provider

```bash
go build -o terraform-provider-freshservice
```

### Running Tests

```bash
go test ./...
```

### Local Installation

To install the provider locally for development:

1. Build the provider:
   ```bash
   go build -o terraform-provider-freshservice
   ```

2. Create a local provider directory:
   ```bash
   mkdir -p ~/.terraform.d/plugins/registry.terraform.io/lcp-llp/freshservice/1.0.0/windows_amd64/
   ```

3. Copy the binary:
   ```bash
   cp terraform-provider-freshservice.exe ~/.terraform.d/plugins/registry.terraform.io/lcp-llp/freshservice/1.0.0/windows_amd64/
   ```

4. Use in your Terraform configuration:
   ```terraform
   terraform {
     required_providers {
       freshservice = {
         source  = "registry.terraform.io/lcp-llp/freshservice"
         version = "1.0.0"
       }
     }
   }
   ```

### Documentation Generation

Generate documentation for the Terraform Registry:

```bash
go generate
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run tests and ensure they pass
6. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For issues and questions:

- Create an issue in this repository
- Check the [Freshservice API documentation](https://api.freshservice.com/)
- Review the provider documentation in the `docs/` directory
