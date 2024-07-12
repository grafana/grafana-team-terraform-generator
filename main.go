package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	msgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
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

	// Get the selected provider
	provider := viper.GetString("provider")

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

	// Generate Terraform file
	err = generateTerraformFile(groups)
	if err != nil {
		slog.Error("Failed to generate Terraform file", "error", err)
		os.Exit(1)
	}

	slog.Info("Terraform file generated successfully", "file", "grafana_teams.tf")
}

func setupConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Set defaults
	viper.SetDefault("provider", "azure")
	return nil
}

func getAzureGroups() ([]Group, error) {
	clientID := viper.GetString("azure.client_id")
	clientSecret := viper.GetString("azure.client_secret")
	tenantID := viper.GetString("azure.tenant_id")

	// Azure AD authentication
	cred, err := azidentity.NewClientSecretCredential(tenantID, clientID, clientSecret, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get Azure credentials: %v", err)
	}

	// Create a new Graph client
	client, err := msgraph.NewGraphServiceClientWithCredentials(cred, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Graph client: %v", err)
	}

	return fetchAllGroups(client)
}

// Original function refactored to use fetchGroupsPage
func fetchAllGroups(client *msgraph.GraphServiceClient) ([]Group, error) {
	var groups []Group
	var nextLink *string
	requestHeaders := abstractions.NewRequestHeaders()
	requestHeaders.Add("ConsistencyLevel", "eventual")

	for {
		result, err := fetchGroupsPage(client, nextLink, requestHeaders)
		if err != nil {
			return nil, err
		}

		groups = append(groups, result.Groups...)

		if result.NextLink == nil {
			break // No more pages
		}
		nextLink = result.NextLink
	}

	return groups, nil
}

// fetchGroupsPage fetches a single page of groups
func fetchGroupsPage(client *msgraph.GraphServiceClient, nextLink *string, requestHeaders *abstractions.RequestHeaders) (fetchGroupsResult, error) {
	var err error
	var result models.GroupCollectionResponseable

	if nextLink == nil {
		// First request
		result, err = client.Groups().Get(context.Background(), &msgroups.GroupsRequestBuilderGetRequestConfiguration{
			QueryParameters: &msgroups.GroupsRequestBuilderGetQueryParameters{
				Select: []string{"displayName", "id"},
				Top:    Int32(300), // Max value allowed
			},
			Headers: requestHeaders,
		})
	} else {
		// Subsequent requests using the next link
		result, err = client.Groups().WithUrl(*nextLink).Get(context.Background(), nil)
	}

	if err != nil {
		return fetchGroupsResult{}, fmt.Errorf("failed to fetch groups: %v", err)
	}

	var groups []Group
	for _, group := range result.GetValue() {
		groups = append(groups, Group{
			Name:       *group.GetDisplayName(),
			Identifier: *group.GetId(),
		})
	}

	return fetchGroupsResult{
		Groups:   groups,
		NextLink: result.GetOdataNextLink(),
	}, nil
}

// Helper function to convert int32 to *int32
func Int32(v int32) *int32 {
	return &v
}

func generateTerraformFile(groups []Group) error {
	file, err := os.Create("grafana_teams.tf")
	if err != nil {
		return err
	}
	defer file.Close()

	for _, group := range groups {
		resourceName := strings.ReplaceAll(strings.ToLower(group.Name), " ", "_")

		_, err = file.WriteString(fmt.Sprintf(`
resource "grafana_team" "%s" {
  name = "%s"
}

resource "grafana_team_external_group" "%s_group" {
  team_id = grafana_team.%s.id
  groups  = ["%s"]
}

`, resourceName, group.Name, resourceName, resourceName, group.Identifier))

		if err != nil {
			return err
		}
	}

	return nil
}
