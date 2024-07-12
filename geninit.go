package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func createGrafanaTerraformStructure(baseDir string) error {
	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("failed to create base directory: %w", err)
	}

	// Create the module directory if it doesn't exist
	moduleDir := filepath.Join(baseDir, "modules")
	if err := os.MkdirAll(moduleDir, 0755); err != nil {
		return fmt.Errorf("failed to create modules directory: %w", err)
	}

	// Create the teams module directory if it doesn't exist
	teamsDir := filepath.Join(moduleDir, "teams")
	if err := os.MkdirAll(teamsDir, 0755); err != nil {
		return fmt.Errorf("failed to create teams module directory: %w", err)
	}

	// File contents
	files := map[string]string{
		filepath.Join(baseDir, "main.tf"): `# Main Terraform configuration file

module "teams" {
  source = "./modules/teams"
  # Add any necessary variables here
}
`,
		filepath.Join(baseDir, "terraform.tf"): `# Terraform settings and provider configurations

terraform {
  required_providers {
    grafana = {
      source  = "grafana/grafana"
      version = "~> 1.28.0"
    }
  }
  required_version = ">= 0.13"
}

provider "grafana" {
  url  = var.grafana_url
  auth = var.grafana_auth
}

variable "grafana_url" {
  type        = string
  description = "The URL of your Grafana instance"
}

variable "grafana_auth" {
  type        = string
  description = "The API key or auth token for Grafana"
  sensitive   = true
}
`,
		filepath.Join(teamsDir, "main.tf"): `# Main configuration for teams module

resource "grafana_team" "example_team" {
  name  = "Example Team"
  email = "example@team.com"
}
`,
		filepath.Join(teamsDir, "variables.tf"): `# Input variables for teams module

variable "team_name" {
  type        = string
  description = "The name of the team to create"
  default     = "Default Team"
}

variable "team_email" {
  type        = string
  description = "The email associated with the team"
  default     = "default@team.com"
}
`,
		filepath.Join(teamsDir, "outputs.tf"): `# Outputs for teams module

output "team_id" {
  value       = grafana_team.example_team.id
  description = "The ID of the created team"
}
`,
	}

	// Create or update files
	for filePath, content := range files {
		if err := createOrUpdateFile(filePath, content); err != nil {
			return fmt.Errorf("failed to create or update %s: %w", filePath, err)
		}
	}

	fmt.Println("Grafana Terraform folder structure created/updated successfully.")
	return nil
}

func createOrUpdateFile(filePath, content string) error {
	// Check if file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// Create new file with content
		return os.WriteFile(filePath, []byte(content), 0644)
	} else if err != nil {
		return err
	}

	// File exists, update its content
	return os.WriteFile(filePath, []byte(content), 0644)
}
