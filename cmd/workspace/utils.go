package workspace

import "ticktask/persistence"

const defaultWorkspace = "default"

func GetSelectedWorkspace() string {
	workspace := persistence.GetDB().GetSelectedWorkspace()
	if len(workspace) > 0 {
		return workspace
	}

	list := persistence.GetDB().GetWorkspaces()
	if len(list) > 0 {
		return list[0]
	}

	return defaultWorkspace
}
