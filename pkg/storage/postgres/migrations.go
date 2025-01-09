package postgres

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func (s *Storage) applyMigrations(dir string) error {
	const op = "storage.postgres.ApplyMigrations"

	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf(`'%s: failed to read migrations directory %s: %w'`, op, dir, err)
	}

	var files []string
	for _, entry := range entries {
		if entry.Type().IsRegular() && strings.HasSuffix(entry.Name(), ".sql") {
			files = append(files, fmt.Sprintf("%s/%s", dir, entry.Name()))
		}
	}

	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf(`'%s: failed to read migration file %s: %w'`, op, file, err)
		}

		_, err = s.db.Exec(string(content))
		if err != nil {
			return fmt.Errorf(`'%s: failed to execute migration %s: %w'`, op, file, err)
		}
	}

	return nil
}
