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
	opt := "team"
	if len(os.Args) > 1 {
		opt = os.Args[1]
	}

	var errCmd error
	switch opt {
	case "team":
		slog.Info("Generating Terraform file for teams")
		errCmd = genTeamResources(baseDir)
	case "team-outputs":
		slog.Info("Generating Terraform file for teams")
		errCmd = genTeamResourcesOutput(baseDir)
	case "init":
		slog.Info("Initializing basic Terraform template")
		errCmd = createGrafanaTerraformStructure(baseDir)
	}

	if errCmd != nil {
		slog.Error("Failed to execute command", "error", errCmd)
		os.Exit(1)
	}
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

	return groups, nil
}

func genTeamResourcesOutput(baseDir string) error {
	teamOutputsFile := "teams-output.tf"
	// if inited with `init` command, the file will be created in grafana_tf directory
	if _, err := os.Stat(baseDir); err == nil {
		teamOutputsFile = baseDir + "/modules/teams/outputs.tf"
	}

	if _, err := os.Stat(teamOutputsFile); err == nil {
		slog.Warn("Terraform file already exists, overwriting", "file", teamOutputsFile)
	}

	groups, err := fetchGroups()
	if err != nil {
		return err
	}

	// Generate Terraform file
	if err = generateOutputsFile(teamOutputsFile, groups); err != nil {
		slog.Error("Failed to generate Terraform file", "error", err)
		os.Exit(1)
	}

	slog.Info("Terraform file generated successfully", "file", teamOutputsFile)

	return nil
}

func genTeamResources(baseDir string) error {
	teamFile := "teams.tf"
	// if inited with `init` command, the file will be created in grafana_tf directory
	if _, err := os.Stat(baseDir); err == nil {
		teamFile = baseDir + "/modules/teams/main.tf"
	}

	if _, err := os.Stat(teamFile); err == nil {
		slog.Warn("Terraform file already exists, overwriting", "file", teamFile)
	}

	groups, err := fetchGroups()
	if err != nil {
		return err
	}

	// Generate Terraform file
	if err = generateTerraformFile(teamFile, groups); err != nil {
		slog.Error("Failed to generate Terraform file", "error", err)
		os.Exit(1)
	}

	slog.Info("Terraform file generated successfully", "file", teamFile)

	return nil
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
