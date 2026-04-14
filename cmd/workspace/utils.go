package workspace

import "ticktask/persistence"

// defaultWorkspace is used when no workspace is selected or exists.
const defaultWorkspace = "default"

// GetSelectedWorkspace returns the currently active workspace name.
// Selection priority:
//  1. Explicitly selected workspace (from SaveSelectedWorkspace)
//  2. First workspace in the list
//  3. "default" if no workspaces exist
//
// This function is used by all task commands to determine which bucket to use.
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
