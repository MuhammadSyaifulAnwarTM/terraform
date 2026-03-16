// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: BUSL-1.1

package migrate

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseMigration(t *testing.T) {
	input := []byte(`{
		"name": "v3to4/rename_s3_bucket_object",
		"description": "Rename aws_s3_bucket_object to aws_s3_object",
		"match": {
			"block_type": "resource",
			"label": "aws_s3_bucket_object"
		},
		"actions": [
			{"action": "rename_resource", "to": "aws_s3_object"}
		]
	}`)

	got, err := ParseMigration(input)
	if err != nil {
		t.Fatal(err)
	}

	want := &Migration{
		Name:        "v3to4/rename_s3_bucket_object",
		Description: "Rename aws_s3_bucket_object to aws_s3_object",
		Match: Match{
			BlockType: "resource",
			Label:     "aws_s3_bucket_object",
		},
		Actions: []Action{
			{Action: "rename_resource", To: "aws_s3_object"},
		},
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("ParseMigration mismatch (-want +got):\n%s", diff)
	}
}

func TestParseMigration_allActions(t *testing.T) {
	input := []byte(`{
		"name": "v4to5/multi_action",
		"description": "Test all action types",
		"match": {"block_type": "resource", "label": "aws_instance"},
		"actions": [
			{"action": "rename_attribute", "from": "ami", "to": "image_id"},
			{"action": "remove_attribute", "name": "vpc"},
			{"action": "rename_resource", "to": "aws_ec2_instance"},
			{"action": "add_comment", "text": "FIXME: check manually"},
			{"action": "set_attribute_value", "name": "engine", "value": "mysql"},
			{"action": "add_attribute", "name": "engine", "value": "aurora"},
			{"action": "replace_value", "name": "enabled", "old_value": "true", "new_value": "\"Enabled\""}
		]
	}`)

	got, err := ParseMigration(input)
	if err != nil {
		t.Fatal(err)
	}

	if len(got.Actions) != 7 {
		t.Fatalf("expected 7 actions, got %d", len(got.Actions))
	}
	if got.Actions[0].Action != "rename_attribute" || got.Actions[0].From != "ami" || got.Actions[0].To != "image_id" {
		t.Errorf("action 0: %+v", got.Actions[0])
	}
	if got.Actions[6].Action != "replace_value" || got.Actions[6].OldValue != "true" || got.Actions[6].NewValue != `"Enabled"` {
		t.Errorf("action 6: %+v", got.Actions[6])
	}
}

func TestParseMigration_invalidJSON(t *testing.T) {
	_, err := ParseMigration([]byte(`not json`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestParseMigration_missingName(t *testing.T) {
	input := []byte(`{
		"match": {"block_type": "resource", "label": "test"},
		"actions": [{"action": "remove_attribute", "name": "foo"}]
	}`)
	_, err := ParseMigration(input)
	if err == nil {
		t.Error("expected error for missing name")
	}
}

func TestDiscoverMigrations(t *testing.T) {
	dir := t.TempDir()

	v3to4 := filepath.Join(dir, "v3to4")
	os.MkdirAll(v3to4, 0755)

	writeJSON(t, filepath.Join(v3to4, "rename_s3.json"), &Migration{
		Name:        "v3to4/rename_s3",
		Description: "Rename S3",
		Match:       Match{BlockType: "resource", Label: "aws_s3_bucket_object"},
		Actions:     []Action{{Action: "rename_resource", To: "aws_s3_object"}},
	})
	writeJSON(t, filepath.Join(v3to4, "remove_classic.json"), &Migration{
		Name:        "v3to4/remove_classic",
		Description: "Remove EC2-Classic",
		Match:       Match{BlockType: "resource", Label: "aws_instance"},
		Actions:     []Action{{Action: "remove_attribute", Name: "vpc_classic_link_id"}},
	})

	// Non-JSON file should be ignored
	os.WriteFile(filepath.Join(v3to4, "readme.txt"), []byte("ignore"), 0644)

	migrations, err := DiscoverMigrations(dir)
	if err != nil {
		t.Fatal(err)
	}

	if len(migrations) != 2 {
		t.Fatalf("expected 2 migrations, got %d", len(migrations))
	}

	names := make([]string, len(migrations))
	for i, m := range migrations {
		names[i] = m.Name
	}
	if names[0] != "v3to4/remove_classic" || names[1] != "v3to4/rename_s3" {
		t.Errorf("unexpected names: %v", names)
	}
}

func TestFilterMigrations(t *testing.T) {
	all := []*Migration{
		{Name: "v3to4/rename_s3"},
		{Name: "v3to4/remove_classic"},
		{Name: "v4to5/rename_elasticache"},
	}

	got := FilterMigrations(all, "v3to4/*")
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}

	got = FilterMigrations(all, "v4to5/*")
	if len(got) != 1 {
		t.Fatalf("expected 1, got %d", len(got))
	}

	got = FilterMigrations(all, "")
	if len(got) != 3 {
		t.Fatalf("expected 3 (no filter), got %d", len(got))
	}
}

func writeJSON(t *testing.T, path string, m *Migration) {
	t.Helper()
	data, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}
}
