package workspace

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	removeStateDirMode  = 0o750
	removeStateFileMode = 0o600
)

type RemoveState struct {
	WorkspaceID         string   `json:"workspace_id"`
	StartedAt           string   `json:"started_at"`
	UpdatedAt           string   `json:"updated_at"`
	ReposTotal          int      `json:"repos_total"`
	RemovedAliases      []string `json:"removed_aliases"`
	PendingAliases      []string `json:"pending_aliases"`
	WorkspaceDirPresent bool     `json:"workspace_dir_present"`
	LastError           string   `json:"last_error,omitempty"`
}

func LoadRemoveState(rootDir, workspaceID string) (RemoveState, bool, error) {
	path, err := removeStatePath(rootDir, workspaceID)
	if err != nil {
		return RemoveState{}, false, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return RemoveState{}, false, nil
		}
		return RemoveState{}, false, fmt.Errorf("read remove state: %w", err)
	}
	var state RemoveState
	if err := json.Unmarshal(data, &state); err != nil {
		return RemoveState{}, false, fmt.Errorf("parse remove state: %w", err)
	}
	return normalizeRemoveState(state), true, nil
}

func SaveRemoveState(rootDir string, state RemoveState) error {
	path, err := removeStatePath(rootDir, state.WorkspaceID)
	if err != nil {
		return err
	}
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, removeStateDirMode); err != nil {
		return fmt.Errorf("create remove state dir: %w", err)
	}
	state = normalizeRemoveState(state)
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal remove state: %w", err)
	}
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".*.tmp")
	if err != nil {
		return fmt.Errorf("create temp remove state: %w", err)
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return fmt.Errorf("write temp remove state: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("close temp remove state: %w", err)
	}
	if err := os.Chmod(tmpName, removeStateFileMode); err != nil {
		return fmt.Errorf("chmod temp remove state: %w", err)
	}
	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("rename remove state: %w", err)
	}
	return nil
}

func DeleteRemoveState(rootDir, workspaceID string) error {
	path, err := removeStatePath(rootDir, workspaceID)
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("delete remove state: %w", err)
	}
	return nil
}

func removeStatePath(rootDir, workspaceID string) (string, error) {
	if strings.TrimSpace(rootDir) == "" {
		return "", fmt.Errorf("root directory is required")
	}
	if strings.TrimSpace(workspaceID) == "" {
		return "", fmt.Errorf("workspace id is required")
	}
	if strings.ContainsAny(workspaceID, `/\`) {
		return "", fmt.Errorf("invalid workspace id: must not contain path separators")
	}
	return filepath.Join(rootDir, MetadataDirName, "state", "operations", "workspace-remove", workspaceID+".json"), nil
}

func NewRemoveState(workspaceID string, aliases []string, workspaceDirPresent bool) RemoveState {
	now := removeStateTimestamp()
	return normalizeRemoveState(RemoveState{
		WorkspaceID:         workspaceID,
		StartedAt:           now,
		UpdatedAt:           now,
		ReposTotal:          len(aliases),
		RemovedAliases:      nil,
		PendingAliases:      append([]string(nil), aliases...),
		WorkspaceDirPresent: workspaceDirPresent,
	})
}

func (s RemoveState) Summary() string {
	return fmt.Sprintf("removed=%v pending=%v workspace_dir_present=%t", s.RemovedAliases, s.PendingAliases, s.WorkspaceDirPresent)
}

func normalizeRemoveState(state RemoveState) RemoveState {
	state.WorkspaceID = strings.TrimSpace(state.WorkspaceID)
	state.StartedAt = strings.TrimSpace(state.StartedAt)
	state.UpdatedAt = strings.TrimSpace(state.UpdatedAt)
	state.LastError = strings.TrimSpace(state.LastError)
	if state.RemovedAliases == nil {
		state.RemovedAliases = []string{}
	}
	if state.PendingAliases == nil {
		state.PendingAliases = []string{}
	}
	if state.ReposTotal == 0 {
		state.ReposTotal = len(state.RemovedAliases) + len(state.PendingAliases)
	}
	return state
}

func removeStateTimestamp() string {
	return time.Now().UTC().Format(time.RFC3339)
}
