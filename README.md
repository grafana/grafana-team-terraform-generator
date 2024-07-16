# Grafana Team Terraform Generator

This tool generates Terraform configurations for Grafana teams, their external group mappings, and associated folders. It supports fetching group information from identity providers and generates a modular Terraform structure.

![image](https://github.com/user-attachments/assets/3740ad5d-f11c-42a4-8abe-f336e8ed83a1)

## Features

- Generates Terraform configurations for Grafana teams and folders
- Supports external group mappings
- Creates a modular Terraform structure
- Configurable through YAML file and environment variables

## Build Instructions

To build the Grafana Terraform Generator, follow these steps:

1. Ensure you have Go installed on your system (version 1.16 or later recommended).
2. Clone this repository:
   ```
   git clone https://github.com/yourusername/grafana-terraform-generator.git
   cd grafana-terraform-generator
   ```
3. Build the binary:
   ```
   go build -o grafana-tf-gen
   ```

## Configuration

Create a `config.yaml` file in the same directory as the binary with the following content:

```yaml
provider: azure
log:
  level: info
# Azure AD configuration (if using Azure AD)
azure:
  client_id: your_azure_client_id
  client_secret: your_azure_client_secret
  tenant_id: your_azure_tenant_id
```

You can also use environment variables to override the configuration:

- `GTF_PROVIDER`: Set the identity provider (e.g., "azure")
- `GTF_LOG_LEVEL`: Set the log level (e.g., "debug", "info")
- `GTF_AZURE_CLIENT_ID`: Set your Azure client ID
- `GTF_AZURE_CLIENT_SECRET`: Set your Azure client secret
- `GTF_AZURE_TENANT_ID`: Set your Azure tenant ID

## Usage

Run the tool to generate the Terraform files after configuring your `config.yaml`:

```
./grafana-tf-gen
```

This will create a `grafana_tf` directory with the following structure:

```
grafana_tf/
├── main.tf
├── terraform.tf
└── modules/
    ├── teams/
    │   ├── main.tf
    │   └── terraform.tf
    └── folders/
        ├── main.tf
        └── terraform.tf
```

## Using the Generated Terraform Files

1. Navigate to the `grafana_tf` directory:
   ```
   cd grafana_tf
   ```

2. Initialize the Terraform working directory:
   ```
   terraform init
   ```

3. Set the required variables:
   - Create a `terraform.tfvars` file with the following content:
     ```
     grafana_url = "https://your-grafana-instance.com"
     grafana_auth = "your-grafana-api-key"
     ```

4. Review the planned changes:
   ```
   terraform plan
   ```

5. Apply the configuration:
   ```
   terraform apply
   ```

This will create the Grafana teams, external group mappings, and folders in your Grafana instance.

## Customizing the Generated Files

- The `main.tf` file in the root directory contains the team definitions and module calls. You can modify this file to add or remove teams.
- The `modules/teams/main.tf` file contains the resources for creating Grafana teams and external group mappings.
- The `modules/folders/main.tf` file contains the resources for creating Grafana folders and setting permissions.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Adding New Providers

To add support for a new identity provider:

1. Create a new file (e.g., `newprovider.go`) with a function to fetch groups from the provider.
2. Update the `fetchGroups` function in `main.go` to support the new provider.
3. Update the README to include the new supported provider.

## License

[MIT License](LICENSE)