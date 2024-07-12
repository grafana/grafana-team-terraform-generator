# Grafana Terraform Generator

This tool generates Terraform configurations for Grafana teams and their external group mappings. It supports multiple identity providers and can generate both resource definitions and outputs.

## Supported Providers

Currently, the tool supports the following identity providers:

- Azure Active Directory

(More providers can be added in the future)

## Build Instructions

To build the Grafana Terraform Generator, follow these steps:

1. Ensure you have Go installed on your system (version 1.16 or later recommended).
2. Clone this repository:
   ```
   git clone https://github.com/grafana/grafana-terraform-generator.git
   cd grafana-terraform-generator
   ```
3. Build the binary:
   ```
   go build -o grafana-tf-gen
   ```

## Installation

After building, you can move the binary to a location in your PATH for easy access:

```
sudo mv grafana-tf-gen /usr/local/bin/
```

## Configuration

Create a `config.yaml` file in the same directory as the binary with the following content:

```yaml
provider: azure
log:
  level: info

# Azure AD configuration
azure:
  client_id: your_azure_client_id
  client_secret: your_azure_client_secret
  tenant_id: your_azure_tenant_id
```

You can check the [`config.sample.yaml`](config.sample.yaml) file for an example configuration.

You can also use environment variables to override the configuration:

- `GTF_PROVIDER`: Set the identity provider (e.g., "azure")
- `GTF_LOG_LEVEL`: Set the log level (e.g., "debug", "info")
- `GTF_AZURE_CLIENT_ID`: Set your Azure client ID
- `GTF_AZURE_CLIENT_SECRET`: Set your Azure client secret
- `GTF_AZURE_TENANT_ID`: Set your Azure tenant ID

## Usage

The tool supports several commands:

1. Initialize a basic Terraform structure:
   ```
   grafana-tf-gen init
   ```

2. Generate Terraform resources for teams:
   ```
   grafana-tf-gen team
   ```

3. Generate Terraform outputs for teams:
   ```
   grafana-tf-gen team-outputs
   ```

### Examples

1. Initialize the Terraform structure:
   ```
   grafana-tf-gen init
   ```
   This will create a `grafana_tf` directory with the basic Terraform structure.

2. Generate team resources:
   ```
   grafana-tf-gen team
   ```
   This will create a `teams.tf` file (or `grafana_tf/modules/teams/main.tf` if initialized) with Grafana team and external group resources.

3. Generate team outputs:
   ```
   grafana-tf-gen team-outputs
   ```
   This will create a `teams-output.tf` file (or `grafana_tf/modules/teams/outputs.tf` if initialized) with outputs for the Grafana team IDs.

4. Run with debug logging:
   ```
   GTF_LOG_LEVEL=debug grafana-tf-gen team
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Adding New Providers

To add support for a new identity provider:

1. Create a new file (e.g., `newprovider.go`) with a function to fetch groups from the provider.
2. Add a case for the new provider in the `fetchGroups` function in `main.go`.
3. Update the README to include the new supported provider.


## License

[MIT License](LICENSE)