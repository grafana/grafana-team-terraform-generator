package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
)

func generateMainTerraformFile(groups []Group) (string, error) {
	var maintf strings.Builder
	maintf.WriteString(`
locals {
  teams = {
`)

	resourceNames := make(map[string]string, len(groups))
	for _, group := range groups {
		resourceName, err := terraformResourceName(group.Name)
		if err != nil {
			return "", err
		}
		if previousName, ok := resourceNames[resourceName]; ok {
			return "", fmt.Errorf("group names %q and %q normalize to the duplicate Terraform key %q", previousName, group.Name, resourceName)
		}
		resourceNames[resourceName] = group.Name

		if _, err := fmt.Fprintf(&maintf, "    %s = {\n      name = %s\n      group_id = %s\n      folder_name = %s\n    }\n",
			quoteHCL(resourceName), quoteHCL(group.Name), quoteHCL(group.Identifier), quoteHCL(group.Name+" Folder")); err != nil {
			return "", fmt.Errorf("failed to write Terraform team definition: %w", err)
		}
	}

	maintf.WriteString(`  }
}

module "teams" {
  source = "./modules/teams"
  teams  = local.teams
}

module "folders" {
  source    = "./modules/folders"
  teams     = local.teams
  team_ids  = module.teams.team_ids
}
`)
	return maintf.String(), nil
}

func terraformResourceName(name string) (string, error) {
	trimmedName := strings.TrimSpace(name)
	if trimmedName == "" {
		return "", fmt.Errorf("group name cannot be empty")
	}

	var builder strings.Builder
	lastWasUnderscore := false
	for _, r := range strings.ToLower(trimmedName) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			builder.WriteRune(r)
			lastWasUnderscore = false
		case r == '_':
			if !lastWasUnderscore {
				builder.WriteRune(r)
				lastWasUnderscore = true
			}
		default:
			if !lastWasUnderscore {
				builder.WriteByte('_')
				lastWasUnderscore = true
			}
		}
	}

	resourceName := strings.Trim(builder.String(), "_")
	if resourceName == "" {
		return "", fmt.Errorf("group name %q contains no usable Terraform key characters", name)
	}
	return resourceName, nil
}

func quoteHCL(value string) string {
	encoded, err := json.Marshal(value)
	if err != nil {
		return `""`
	}
	return string(encoded)
}

const TeamModuleMain = `variable "teams" {
  description = "Map of team configurations"
  type = map(object({
    name      = string
    group_id  = string
    folder_name = string
  }))
}

resource "grafana_team" "teams" {
  for_each = var.teams
  name     = each.value.name
}

resource "grafana_team_external_group" "team_groups" {
  for_each = var.teams
  team_id  = resource.grafana_team.teams[each.key].id
  groups = [each.value.group_id]
}

output "team_ids" {
  value = {
    for key, team in resource.grafana_team.teams : key => team.id
  }
  description = "Map of team names to their IDs"
}
`

const FolderModuleMain = `variable "teams" {
  description = "Map of team configurations"
  type = map(object({
    name        = string
    group_id    = string
    folder_name = string
  }))
}

variable "team_ids" {
  description = "Map of team names to their IDs"
  type        = map(string)
}

resource "grafana_folder" "folders" {
  for_each = var.teams
  title    = each.value.folder_name
}

resource "grafana_folder_permission" "folder_permissions" {
  for_each   = grafana_folder.folders
  folder_uid = each.value.uid
  permissions {
    team_id    = var.team_ids[each.key]
    permission = "Admin"
  }
}

output "folder_ids" {
  value = {
    for key, folder in grafana_folder.folders : key => folder.id
  }
  description = "Map of folder names to their IDs"
}
`
