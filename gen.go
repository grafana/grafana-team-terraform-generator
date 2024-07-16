package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func createGrafanaTerraformStructure(baseDir string) error {
	groups, err := fetchGroups()
	if err != nil {
		return err
	}

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

	// Create the teams module directory if it doesn't exist
	foldersDir := filepath.Join(moduleDir, "folders")
	if err := os.MkdirAll(foldersDir, 0755); err != nil {
		return fmt.Errorf("failed to create folders module directory: %w", err)
	}

	// File contents
	files := map[string]string{
		filepath.Join(baseDir, "main.tf"): generateMainTerraformFile(groups),
		filepath.Join(baseDir, "terraform.tf"): `# Terraform settings and provider configurations

terraform {
  required_providers {
    grafana = {
      source  = "grafana/grafana"
      version = "~> 3.4.0"
    }
  }
}

provider "grafana" {
  url  = var.grafana_url
  auth = var.grafana_auth
  retries = 5
  retry_wait = 10
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
		filepath.Join(teamsDir, "main.tf"):   TeamModuleMain,
		filepath.Join(foldersDir, "main.tf"): FolderModuleMain,
		filepath.Join(teamsDir, "terraform.tf"): `# Terraform settings and provider configurations for teams module

terraform {
  required_providers {
    grafana = {
      source  = "grafana/grafana"
      version = "~> 3.4.0"
    }
  }
}
`,
		filepath.Join(foldersDir, "terraform.tf"): `# Terraform settings and provider configurations for folders module

terraform {
  required_providers {
    grafana = {
      source  = "grafana/grafana"
      version = "~> 3.4.0"
    }
  }
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

func fetchGroups() ([]Group, error) {
	if len(cachedGroups) > 0 {
		slog.Debug("Using cached groups")
		return cachedGroups, nil
	}

	// Get the selected provider
	provider := viper.GetString("provider")
	var groups []Group
	var err error

	switch provider {
	case "azure":
		groups, err = getAzureGroups()
	// Add cases for other providers here
	default:
		slog.Error("Unsupported provider", "provider", provider)
		os.Exit(1)
	}

	if err != nil {
		slog.Error("Failed to fetch groups", "error", err)
		os.Exit(1)
	}

	cachedGroups = groups
	return groups, nil
}
