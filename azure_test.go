package main

import (
	"testing"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func TestGroupsFromResponseSkipsIncompleteGroups(t *testing.T) {
	validName := "Platform Engineering"
	validID := "group-id"
	missingIDName := "Missing ID"

	valid := models.NewGroup()
	valid.SetDisplayName(&validName)
	valid.SetId(&validID)

	missingID := models.NewGroup()
	missingID.SetDisplayName(&missingIDName)

	groups := groupsFromResponse([]models.Groupable{valid, missingID, nil})
	if len(groups) != 1 {
		t.Fatalf("groupsFromResponse() returned %d groups, want 1", len(groups))
	}
	if groups[0].Name != validName || groups[0].Identifier != validID {
		t.Fatalf("groupsFromResponse() returned %+v, want name %q and ID %q", groups[0], validName, validID)
	}
}
