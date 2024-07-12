package main

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type Group struct {
	Name       string
	Identifier string
}

// Define a struct to hold the function's return values
type fetchGroupsResult struct {
	Groups   []Group
	NextLink *string
}

func main() {
	// Set up Viper configuration
	err := setupConfig()
	if err != nil {
		slog.Error("Failed to set up configuration", "error", err)
		os.Exit(1)
	}

	slog.SetLogLoggerLevel(slog.LevelInfo)
	viperLevel := viper.GetString("log.level")
	if viperLevel == "debug" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	// Get the selected provider
	provider := viper.GetString("provider")

	baseDir := "grafana_tf"
	if len(os.Args) > 1 && os.Args[1] == "init" {
		slog.Info("Initializing basic Terraform template")
		err := createGrafanaTerraformStructure(baseDir)
		if err != nil {
			slog.Error("Failed to create Terraform structure", "error", err)
			os.Exit(1)
		}
		return
	}

	var groups []Group

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

	teamFile := "main.tf"
	// if inited with `init` command, the file will be created in grafana_tf directory
	if _, err := os.Stat(baseDir); err == nil {
		teamFile = baseDir + "/modules/teams/main.tf"
	}

	if _, err := os.Stat(teamFile); err == nil {
		slog.Warn("Terraform file already exists, overwriting", "file", teamFile)
	}

	// Generate Terraform file
	err = generateTerraformFile(teamFile, groups)
	if err != nil {
		slog.Error("Failed to generate Terraform file", "error", err)
		os.Exit(1)
	}

	slog.Info("Terraform file generated successfully", "file", teamFile)
}

func setupConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		slog.Warn("failed to read config file. Only environment variables will be used", "error", err)
	}

	// Set defaults
	viper.SetDefault("provider", "azure")
	return nil
}
