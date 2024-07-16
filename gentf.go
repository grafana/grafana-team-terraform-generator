package main

import (
	"fmt"
	"strings"
)

func generateMainTerraformFile(groups []Group) string {
	// Write the locals block for teams
	maintf := `
locals {
  teams = {
`
	// Generate the teams map
	for _, group := range groups {
		resourceName := strings.ReplaceAll(strings.ToLower(group.Name), " ", "_")
		maintf += fmt.Sprintf("    %s = {\n      name = \"%s\"\n      group_id = \"%s\"\n      folder_name = \"%s Folder\"\n    }\n",
			resourceName, group.Name, group.Identifier, group.Name)
	}
	// Close the locals block
	maintf += `  }
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
`
	return maintf
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
