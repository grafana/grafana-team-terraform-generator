package main

import (
	"strings"
	"testing"
)

func TestGenerateMainTerraformFileQuotesValuesAndKeys(t *testing.T) {
	terraform, err := generateMainTerraformFile([]Group{{
		Name:       `2024 Ops "Blue"`,
		Identifier: `group-${identifier}`,
	}})
	if err != nil {
		t.Fatalf("generateMainTerraformFile() error = %v", err)
	}

	for _, expected := range []string{
		`"2024_ops_blue" = {`,
		`name = "2024 Ops \"Blue\""`,
		`group_id = "group-${identifier}"`,
		`folder_name = "2024 Ops \"Blue\" Folder"`,
	} {
		if !strings.Contains(terraform, expected) {
			t.Errorf("generated Terraform does not contain %q:\n%s", expected, terraform)
		}
	}
}

func TestGenerateMainTerraformFileRejectsDuplicateKeys(t *testing.T) {
	generated, err := generateMainTerraformFile([]Group{
		{Name: "Test Team", Identifier: "one"},
		{Name: "test-team", Identifier: "two"},
	})
	if generated != "" {
		t.Errorf("duplicate-key generation returned output: %s", generated)
	}
	if err == nil || !strings.Contains(err.Error(), "duplicate Terraform key") {
		t.Fatalf("expected duplicate key error, got %v", err)
	}
}

func TestGenerateMainTerraformFileRejectsEmptyName(t *testing.T) {
	generated, err := generateMainTerraformFile([]Group{{Name: "   ", Identifier: "group-id"}})
	if generated != "" {
		t.Errorf("empty-name generation returned output: %s", generated)
	}
	if err == nil || !strings.Contains(err.Error(), "cannot be empty") {
		t.Fatalf("expected empty name error, got %v", err)
	}
}
