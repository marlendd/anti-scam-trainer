package scenario

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// SeedFile описывает одну версию сценария в каталоге seeds.
// Указатель IsActive позволяет отличить false от отсутствующего обязательного поля.
type SeedFile struct {
	ID          ScenarioID        `json:"id"`
	LogicalID   LogicalScenarioID `json:"logical_id"`
	Version     int               `json:"version"`
	Role        Role              `json:"role"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	IsActive    *bool             `json:"is_active"`
	Content     Content           `json:"content"`
}

// LoadSeedFiles строго декодирует и проверяет все JSON-файлы каталога.
func LoadSeedFiles(directory string) ([]SeedFile, error) {
	paths, err := filepath.Glob(filepath.Join(directory, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("list scenario seeds: %w", err)
	}
	if len(paths) == 0 {
		return nil, fmt.Errorf("no scenario seed files found in %q", directory)
	}
	sort.Strings(paths)

	seeds := make([]SeedFile, 0, len(paths))
	seenIDs := make(map[ScenarioID]string, len(paths))
	type versionKey struct {
		logicalID LogicalScenarioID
		version   int
	}
	seenVersions := make(map[versionKey]string, len(paths))
	activeVersions := make(map[LogicalScenarioID]string, len(paths))

	for _, path := range paths {
		seed, err := decodeSeedFile(path)
		if err != nil {
			return nil, err
		}
		if err := validateSeedMetadata(seed); err != nil {
			return nil, fmt.Errorf("validate scenario seed %q: %w", path, err)
		}

		if previous, exists := seenIDs[seed.ID]; exists {
			return nil, fmt.Errorf("duplicate scenario seed id %q in %q and %q", seed.ID, previous, path)
		}
		seenIDs[seed.ID] = path

		key := versionKey{logicalID: seed.LogicalID, version: seed.Version}
		if previous, exists := seenVersions[key]; exists {
			return nil, fmt.Errorf(
				"duplicate scenario seed version (%q, %d) in %q and %q",
				seed.LogicalID,
				seed.Version,
				previous,
				path,
			)
		}
		seenVersions[key] = path

		if *seed.IsActive {
			if previous, exists := activeVersions[seed.LogicalID]; exists {
				return nil, fmt.Errorf(
					"multiple active seed versions for logical id %q in %q and %q",
					seed.LogicalID,
					previous,
					path,
				)
			}
			activeVersions[seed.LogicalID] = path
		}

		seeds = append(seeds, seed)
	}

	return seeds, nil
}

// ApplySeedFiles загружает все seeds одной транзакцией и не изменяет уже существующие версии.
func ApplySeedFiles(ctx context.Context, db *sql.DB, directory string) (int, error) {
	seeds, err := LoadSeedFiles(directory)
	if err != nil {
		return 0, err
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin scenario seed transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const query = `
		INSERT INTO scenario_versions (
			id,
			logical_id,
			version,
			role,
			title,
			description,
			is_active,
			content
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO NOTHING
	`

	for _, seed := range seeds {
		content, err := json.Marshal(seed.Content)
		if err != nil {
			return 0, fmt.Errorf("encode scenario seed %q content: %w", seed.ID, err)
		}

		if _, err := tx.ExecContext(
			ctx,
			query,
			seed.ID,
			seed.LogicalID,
			seed.Version,
			seed.Role,
			seed.Title,
			seed.Description,
			*seed.IsActive,
			content,
		); err != nil {
			return 0, fmt.Errorf("upsert scenario seed %q: %w", seed.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit scenario seeds: %w", err)
	}

	return len(seeds), nil
}

func decodeSeedFile(path string) (SeedFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return SeedFile{}, fmt.Errorf("open scenario seed %q: %w", path, err)
	}
	defer func() {
		_ = file.Close()
	}()

	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()

	var seed SeedFile
	if err := decoder.Decode(&seed); err != nil {
		return SeedFile{}, fmt.Errorf("decode scenario seed %q: %w", path, err)
	}

	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		if err == nil {
			err = errors.New("multiple JSON values")
		}
		return SeedFile{}, fmt.Errorf("decode scenario seed %q: %w", path, err)
	}

	return seed, nil
}

func validateSeedMetadata(seed SeedFile) error {
	if !isUUID(string(seed.ID)) {
		return fmt.Errorf("invalid id UUID %q", seed.ID)
	}
	if !isUUID(string(seed.LogicalID)) {
		return fmt.Errorf("invalid logical_id UUID %q", seed.LogicalID)
	}
	if seed.Version <= 0 {
		return fmt.Errorf("version must be positive: got %d", seed.Version)
	}
	if seed.Role != RoleBuyer && seed.Role != RoleSeller {
		return fmt.Errorf("unknown role %q", seed.Role)
	}
	if strings.TrimSpace(seed.Title) == "" {
		return errors.New("title is required")
	}
	if strings.TrimSpace(seed.Description) == "" {
		return errors.New("description is required")
	}
	if seed.IsActive == nil {
		return errors.New("is_active is required")
	}

	s := Scenario{
		ID:          seed.ID,
		LogicalID:   seed.LogicalID,
		Version:     seed.Version,
		Role:        seed.Role,
		Title:       seed.Title,
		Description: seed.Description,
		StartNodeID: seed.Content.StartNodeID,
		Nodes:       seed.Content.Nodes,
		Endings:     seed.Content.Endings,
	}
	if err := Validate(s); err != nil {
		return fmt.Errorf("invalid scenario graph: %w", err)
	}

	return nil
}

func isUUID(value string) bool {
	if len(value) != 36 || value[8] != '-' || value[13] != '-' || value[18] != '-' || value[23] != '-' {
		return false
	}

	compact := strings.ReplaceAll(value, "-", "")
	decoded := make([]byte, 16)
	_, err := hex.Decode(decoded, []byte(compact))
	return err == nil
}
