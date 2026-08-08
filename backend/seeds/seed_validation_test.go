package seeds_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/stretchr/testify/require"
)

func TestScenarioSeedsAreValid(t *testing.T) {
	seeds, err := scenario.LoadSeedFiles(".")

	require.NoError(t, err)
	require.Len(t, seeds, 4)

	roleCounts := map[scenario.Role]int{}
	for _, seed := range seeds {
		roleCounts[seed.Role]++
	}
	require.Equal(t, 2, roleCounts[scenario.RoleBuyer])
	require.Equal(t, 2, roleCounts[scenario.RoleSeller])
}

func TestLoadSeedFilesRejectsInvalidFiles(t *testing.T) {
	validSeed := `{
		"id":"11111111-1111-4111-8111-111111111111",
		"logical_id":"22222222-2222-4222-8222-222222222222",
		"version":1,
		"role":"buyer",
		"title":"Test",
		"description":"Test seed",
		"is_active":true,
		"content":{
			"start_node_id":"start",
			"nodes":[
				{"id":"start","author":"seller","text":"Start","choices":[
					{"id":"s1","text":"1","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"next_node_id":"middle"},
					{"id":"s2","text":"2","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"next_node_id":"middle"},
					{"id":"s3","text":"3","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"next_node_id":"middle"},
					{"id":"s4","text":"4","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"next_node_id":"middle"}
				]},
				{"id":"middle","author":"seller","text":"Middle","choices":[
					{"id":"m1","text":"1","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"next_node_id":"final"},
					{"id":"m2","text":"2","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"next_node_id":"final"},
					{"id":"m3","text":"3","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"next_node_id":"final"},
					{"id":"m4","text":"4","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"next_node_id":"final"}
				]},
				{"id":"final","author":"seller","text":"Final","choices":[
					{"id":"f1","text":"1","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"ending_id":"safe"},
					{"id":"f2","text":"2","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"ending_id":"safe"},
					{"id":"f3","text":"3","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"ending_id":"risk"},
					{"id":"f4","text":"4","consequence":"c","explanation":"e","weight":1,"score":100,"risk_categories":[],"ending_id":"risk"}
				]}
			],
			"endings":[
				{"id":"safe","header":"Safe","result":"Safe result"},
				{"id":"risk","header":"Risk","result":"Risk result"}
			]
		}
	}`

	tests := []struct {
		name    string
		content string
		wantErr string
	}{
		{
			name:    "unknown field",
			content: replaceOnce(t, validSeed, `"is_active":true`, `"is_active":true,"is_actve":true`),
			wantErr: "unknown field",
		},
		{
			name:    "missing is active",
			content: replaceOnce(t, validSeed, `"is_active":true,`, ""),
			wantErr: "is_active is required",
		},
		{
			name:    "invalid id UUID",
			content: replaceOnce(t, validSeed, "11111111-1111-4111-8111-111111111111", "not-a-uuid"),
			wantErr: "invalid id UUID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			writeSeed(t, directory, "seed.json", test.content)

			_, err := scenario.LoadSeedFiles(directory)

			require.ErrorContains(t, err, test.wantErr)
		})
	}
}

func TestLoadSeedFilesRejectsDuplicateIDsAndVersions(t *testing.T) {
	first := readProjectSeed(t)

	tests := []struct {
		name       string
		secondSeed string
		wantErr    string
	}{
		{
			name:       "duplicate id",
			secondSeed: replaceOnce(t, first, "5d473394-7f64-4e92-90ac-3159ed53f2e7", "33333333-3333-4333-8333-333333333333"),
			wantErr:    "duplicate scenario seed id",
		},
		{
			name:       "duplicate logical id and version",
			secondSeed: replaceOnce(t, first, "45d4cc8c-f604-4a7c-b8c5-f2464717b71f", "33333333-3333-4333-8333-333333333333"),
			wantErr:    "duplicate scenario seed version",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			directory := t.TempDir()
			writeSeed(t, directory, "first.json", first)
			writeSeed(t, directory, "second.json", test.secondSeed)

			_, err := scenario.LoadSeedFiles(directory)

			require.ErrorContains(t, err, test.wantErr)
		})
	}
}

func readProjectSeed(t *testing.T) string {
	t.Helper()
	raw, err := os.ReadFile("001_buyer_gpu_fake_support.json")
	require.NoError(t, err)
	return string(raw)
}

func writeSeed(t *testing.T, directory string, name string, content string) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(directory, name), []byte(content), 0o600))
}

func replaceOnce(t *testing.T, value string, old string, replacement string) string {
	t.Helper()
	result := strings.Replace(value, old, replacement, 1)
	require.NotEqual(t, value, result, "test fixture replacement must match")
	return result
}
