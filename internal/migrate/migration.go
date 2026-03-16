// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package migrate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Migration represents a single JSON migration file.
type Migration struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Match       Match    `json:"match"`
	Actions     []Action `json:"actions"`
}

// Match specifies which blocks to target.
type Match struct {
	BlockType string `json:"block_type"`
	Label     string `json:"label"`
}

// Action describes a single mutation to apply to matched blocks.
type Action struct {
	Action   string `json:"action"`
	From     string `json:"from,omitempty"`
	To       string `json:"to,omitempty"`
	Name     string `json:"name,omitempty"`
	Text     string `json:"text,omitempty"`
	Value    string `json:"value,omitempty"`
	OldValue string `json:"old_value,omitempty"`
	NewValue string `json:"new_value,omitempty"`
}

// ParseMigration parses a JSON migration file.
func ParseMigration(data []byte) (*Migration, error) {
	var m Migration
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing migration JSON: %w", err)
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Migration) validate() error {
	if m.Name == "" {
		return fmt.Errorf("migration missing required field \"name\"")
	}
	if m.Match.BlockType == "" {
		return fmt.Errorf("migration %q: match missing required field \"block_type\"", m.Name)
	}
	if len(m.Actions) == 0 {
		return fmt.Errorf("migration %q: must have at least one action", m.Name)
	}
	for i, a := range m.Actions {
		if err := a.validate(); err != nil {
			return fmt.Errorf("migration %q action %d: %w", m.Name, i, err)
		}
	}
	return nil
}

var validActions = map[string]bool{
	"rename_attribute":    true,
	"remove_attribute":    true,
	"rename_resource":     true,
	"add_comment":         true,
	"set_attribute_value": true,
	"add_attribute":       true,
	"replace_value":       true,
}

func (a *Action) validate() error {
	if !validActions[a.Action] {
		return fmt.Errorf("unknown action %q", a.Action)
	}
	switch a.Action {
	case "rename_attribute":
		if a.From == "" || a.To == "" {
			return fmt.Errorf("rename_attribute requires \"from\" and \"to\"")
		}
	case "remove_attribute":
		if a.Name == "" {
			return fmt.Errorf("remove_attribute requires \"name\"")
		}
	case "rename_resource":
		if a.To == "" {
			return fmt.Errorf("rename_resource requires \"to\"")
		}
	case "add_comment":
		if a.Text == "" {
			return fmt.Errorf("add_comment requires \"text\"")
		}
	case "set_attribute_value":
		if a.Name == "" || a.Value == "" {
			return fmt.Errorf("set_attribute_value requires \"name\" and \"value\"")
		}
	case "add_attribute":
		if a.Name == "" || a.Value == "" {
			return fmt.Errorf("add_attribute requires \"name\" and \"value\"")
		}
	case "replace_value":
		if a.Name == "" || a.OldValue == "" || a.NewValue == "" {
			return fmt.Errorf("replace_value requires \"name\", \"old_value\", and \"new_value\"")
		}
	}
	return nil
}

// DiscoverMigrations recursively finds all *.json files under dir,
// parses them as migrations, and returns them sorted by name.
func DiscoverMigrations(dir string) ([]*Migration, error) {
	var migrations []*Migration

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}
		m, err := ParseMigration(data)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
		migrations = append(migrations, m)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})
	return migrations, nil
}

// FilterMigrations returns only migrations whose name matches the glob pattern.
// An empty pattern matches all migrations.
func FilterMigrations(migrations []*Migration, pattern string) []*Migration {
	if pattern == "" {
		return migrations
	}
	var result []*Migration
	for _, m := range migrations {
		matched, err := filepath.Match(pattern, m.Name)
		if err != nil {
			continue
		}
		if matched {
			result = append(result, m)
		}
	}
	return result
}
