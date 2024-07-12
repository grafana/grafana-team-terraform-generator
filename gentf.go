package main

import (
	"fmt"
	"os"
	"strings"
)

func generateTerraformFile(filepath string, groups []Group) error {
	file, err := os.Create(filepath)
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
