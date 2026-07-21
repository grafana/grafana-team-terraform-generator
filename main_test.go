package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestSetupConfigReadsPrefixedEnvironmentVariables(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	t.Setenv("GTF_PROVIDER", "azure")
	t.Setenv("GTF_LOG_LEVEL", "debug")
	t.Setenv("GTF_AZURE_CLIENT_ID", "client-id")
	t.Setenv("GTF_AZURE_CLIENT_SECRET", "client-secret")
	t.Setenv("GTF_AZURE_TENANT_ID", "tenant-id")

	workingDirectory := t.TempDir()
	oldWorkingDirectory, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(workingDirectory); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWorkingDirectory) })

	if err := setupConfig(); err != nil {
		t.Fatalf("setupConfig() error = %v", err)
	}

	checks := map[string]string{
		"provider":            "azure",
		"log.level":           "debug",
		"azure.client_id":     "client-id",
		"azure.client_secret": "client-secret",
		"azure.tenant_id":     "tenant-id",
	}
	for key, want := range checks {
		if got := viper.GetString(key); got != want {
			t.Errorf("viper.GetString(%q) = %q, want %q", key, got, want)
		}
	}

	if _, err := os.Stat(filepath.Join(workingDirectory, "config.yaml")); !os.IsNotExist(err) {
		t.Fatalf("test unexpectedly created config.yaml: %v", err)
	}
}

func TestFetchGroupsRejectsUnsupportedProvider(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)
	viper.Set("provider", "unsupported")
	cachedGroups = nil

	_, err := fetchGroups()
	if err == nil {
		t.Fatal("fetchGroups() returned nil error for unsupported provider")
	}
}
