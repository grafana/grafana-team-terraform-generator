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

var cachedGroups []Group // Define a global variable to hold the groups as a cache

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

	baseDir := "grafana_tf"

	slog.Info("Generating Terraform files for teams")
	errCmd := createGrafanaTerraformStructure(baseDir)
	if errCmd != nil {
		slog.Error("Failed to execute command", "error", errCmd)
		os.Exit(1)
	}
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
