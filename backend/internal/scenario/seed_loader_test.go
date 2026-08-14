package scenario_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/marlendd/anti-scam-trainer/internal/scenario"
	"github.com/marlendd/anti-scam-trainer/internal/testfixture"
	"github.com/stretchr/testify/require"
)

const (
	testSeedID        = scenario.ScenarioID("11111111-1111-1111-1111-111111111111")
	testSeedLogicalID = scenario.LogicalScenarioID("22222222-2222-2222-2222-222222222222")
)

func TestLoadSeedFiles(t *testing.T) {
	t.Parallel()

	t.Run("loads files in name order", func(t *testing.T) {
		t.Parallel()

		directory := t.TempDir()
		second := validSeedFile()
		second.ID = "33333333-3333-3333-3333-333333333333"
		second.Version = 2
		secondActive := false
		second.IsActive = &secondActive
		writeSeedFile(t, directory, "20-second.json", second)
		writeSeedFile(t, directory, "10-first.json", validSeedFile())

		seeds, err := scenario.LoadSeedFiles(directory)

		require.NoError(t, err)
		require.Len(t, seeds, 2)
		require.Equal(t, testSeedID, seeds[0].ID)
		require.Equal(t, second.ID, seeds[1].ID)
	})

	t.Run("rejects a directory without JSON files", func(t *testing.T) {
		t.Parallel()

		seeds, err := scenario.LoadSeedFiles(t.TempDir())

		require.ErrorContains(t, err, "no scenario seed files found")
		require.Nil(t, seeds)
	})

	t.Run("rejects malformed JSON", func(t *testing.T) {
		t.Parallel()

		directory := t.TempDir()
		require.NoError(t, os.WriteFile(filepath.Join(directory, "seed.json"), []byte("{"), 0o600))

		_, err := scenario.LoadSeedFiles(directory)

		require.ErrorContains(t, err, "decode scenario seed")
	})

	t.Run("rejects unknown JSON fields", func(t *testing.T) {
		t.Parallel()

		directory := t.TempDir()
		data := marshalSeedFile(t, validSeedFile())
		data[len(data)-1] = ','
		data = append(data, []byte(`"unexpected":true}`)...)
		require.NoError(t, os.WriteFile(filepath.Join(directory, "seed.json"), data, 0o600))

		_, err := scenario.LoadSeedFiles(directory)

		require.ErrorContains(t, err, "unknown field")
	})

	t.Run("rejects multiple JSON values", func(t *testing.T) {
		t.Parallel()

		directory := t.TempDir()
		data := append(marshalSeedFile(t, validSeedFile()), []byte("\n{}")...)
		require.NoError(t, os.WriteFile(filepath.Join(directory, "seed.json"), data, 0o600))

		_, err := scenario.LoadSeedFiles(directory)

		require.ErrorContains(t, err, "multiple JSON values")
	})
}

func TestLoadSeedFilesRejectsInvalidMetadata(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		mutate  func(*scenario.SeedFile)
		message string
	}{
		{
			name: "invalid scenario ID",
			mutate: func(seed *scenario.SeedFile) {
				seed.ID = "not-a-uuid"
			},
			message: "invalid id UUID",
		},
		{
			name: "invalid logical ID",
			mutate: func(seed *scenario.SeedFile) {
				seed.LogicalID = "not-a-uuid"
			},
			message: "invalid logical_id UUID",
		},
		{
			name: "non-hex UUID",
			mutate: func(seed *scenario.SeedFile) {
				seed.ID = "zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz"
			},
			message: "invalid id UUID",
		},
		{
			name: "non-positive version",
			mutate: func(seed *scenario.SeedFile) {
				seed.Version = 0
			},
			message: "version must be positive",
		},
		{
			name: "unknown role",
			mutate: func(seed *scenario.SeedFile) {
				seed.Role = "moderator"
			},
			message: "unknown role",
		},
		{
			name: "blank title",
			mutate: func(seed *scenario.SeedFile) {
				seed.Title = " \t"
			},
			message: "title is required",
		},
		{
			name: "blank description",
			mutate: func(seed *scenario.SeedFile) {
				seed.Description = "\n"
			},
			message: "description is required",
		},
		{
			name: "missing active flag",
			mutate: func(seed *scenario.SeedFile) {
				seed.IsActive = nil
			},
			message: "is_active is required",
		},
		{
			name: "blank product title",
			mutate: func(seed *scenario.SeedFile) {
				seed.Content.Product.Title = " "
			},
			message: "product title is required",
		},
		{
			name: "non-positive product price",
			mutate: func(seed *scenario.SeedFile) {
				seed.Content.Product.Price = 0
			},
			message: "product price must be positive",
		},
		{
			name: "invalid graph",
			mutate: func(seed *scenario.SeedFile) {
				seed.Content.StartNodeID = ""
			},
			message: "invalid scenario graph",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			directory := t.TempDir()
			seed := validSeedFile()
			tc.mutate(&seed)
			writeSeedFile(t, directory, "seed.json", seed)

			_, err := scenario.LoadSeedFiles(directory)

			require.ErrorContains(t, err, tc.message)
		})
	}
}

func TestLoadSeedFilesRejectsDuplicateDefinitions(t *testing.T) {
	t.Parallel()

	t.Run("duplicate ID", func(t *testing.T) {
		t.Parallel()

		directory := t.TempDir()
		first := validSeedFile()
		second := validSeedFile()
		second.LogicalID = "33333333-3333-3333-3333-333333333333"
		writeSeedFile(t, directory, "first.json", first)
		writeSeedFile(t, directory, "second.json", second)

		_, err := scenario.LoadSeedFiles(directory)

		require.ErrorContains(t, err, "duplicate scenario seed id")
	})

	t.Run("duplicate logical version", func(t *testing.T) {
		t.Parallel()

		directory := t.TempDir()
		first := validSeedFile()
		second := validSeedFile()
		second.ID = "33333333-3333-3333-3333-333333333333"
		writeSeedFile(t, directory, "first.json", first)
		writeSeedFile(t, directory, "second.json", second)

		_, err := scenario.LoadSeedFiles(directory)

		require.ErrorContains(t, err, "duplicate scenario seed version")
	})

	t.Run("multiple active versions", func(t *testing.T) {
		t.Parallel()

		directory := t.TempDir()
		first := validSeedFile()
		second := validSeedFile()
		second.ID = "33333333-3333-3333-3333-333333333333"
		second.Version = 2
		writeSeedFile(t, directory, "first.json", first)
		writeSeedFile(t, directory, "second.json", second)

		_, err := scenario.LoadSeedFiles(directory)

		require.ErrorContains(t, err, "multiple active seed versions")
	})
}

func validSeedFile() scenario.SeedFile {
	fixture := testfixture.ValidScenario()
	isActive := true

	return scenario.SeedFile{
		ID:               testSeedID,
		LogicalID:        testSeedLogicalID,
		Version:          1,
		Role:             fixture.Role,
		Title:            fixture.Title,
		Description:      fixture.Description,
		IsActive:         &isActive,
		RewardFragmentID: fixture.RewardFragmentID,
		Content: scenario.Content{
			Product:             fixture.Product,
			StartNodeID:         fixture.StartNodeID,
			SuccessfulEndingIDs: fixture.SuccessfulEndingIDs,
			Nodes:               fixture.Nodes,
			Endings:             fixture.Endings,
		},
	}
}

func writeSeedFile(t *testing.T, directory, name string, seed scenario.SeedFile) {
	t.Helper()
	require.NoError(t, os.WriteFile(filepath.Join(directory, name), marshalSeedFile(t, seed), 0o600))
}

func marshalSeedFile(t *testing.T, seed scenario.SeedFile) []byte {
	t.Helper()
	data, err := json.Marshal(seed)
	require.NoError(t, err)
	return data
}
