package rules_test

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/TiesDO/gdlint/rules"
)

func loadFixture(name string) ([]byte, error) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, fmt.Errorf("failed to extract filename from caller stack")
	}

	rootDir := filepath.Dir(filepath.Dir(filename))

	path := filepath.Join(rootDir, "fixtures", "scripts", name)
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func createDefaultRunnerFromFixture(fixture_path string) (rules.RuleRunner, error) {
	source, err := loadFixture(fixture_path)

	if err != nil {
		return rules.RuleRunner{}, err
	}

	return rules.NewRuleRunner(&rules.DefaultRuleRegistry, source), nil
}
