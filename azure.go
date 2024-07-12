package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraph "github.com/microsoftgraph/msgraph-sdk-go"
	msgroups "github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/spf13/viper"
)

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

	pageNumber := 0
	for {
		slog.Debug("Fetching groups page", "page", pageNumber)
		result, err := fetchGroupsPage(client, nextLink, requestHeaders)
		if err != nil {
			return nil, err
		}

		groups = append(groups, result.Groups...)

		if result.NextLink == nil {
			break // No more pages
		}
		nextLink = result.NextLink
		pageNumber++
	}

	slog.Info("Fetched all groups", "count", len(groups))

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
