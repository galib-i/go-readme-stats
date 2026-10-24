package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeConfig(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "stats.yml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	return path
}

func TestLoadConfig_AppliesDefaults(t *testing.T) {
	cfg, err := loadConfig(writeConfig(t, "cards:\n  - output: card.svg\n  - theme: light\n"))
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}

	want := []cardConfig{
		{Output: "card.svg", Theme: "dark", Header: "Languages", Mode: "bytes", MaxLangs: 6},
		{Output: "languages.svg", Theme: "light", Header: "Languages", Mode: "bytes", MaxLangs: 6},
	}

	if !reflect.DeepEqual(cfg.Cards, want) {
		t.Errorf("loadConfig() = %+v, want %+v", cfg.Cards, want)
	}
}

func TestLoadConfig_KeepsExplicitValues(t *testing.T) {
	cfg, err := loadConfig(writeConfig(t, `cards:
  - output: languages-dark.svg
    theme: soft-dark
    header: Top Languages
    mode: geometric
    max_langs: 3
    ignore:
      - HTML
      - CSS
`))
	if err != nil {
		t.Fatalf("loadConfig() error = %v", err)
	}

	want := cardConfig{
		Output:   "languages-dark.svg",
		Theme:    "soft-dark",
		Header:   "Top Languages",
		Mode:     "geometric",
		MaxLangs: 3,
		Ignore:   []string{"HTML", "CSS"},
	}

	if !reflect.DeepEqual(cfg.Cards, []cardConfig{want}) {
		t.Errorf("loadConfig() = %+v, want %+v", cfg.Cards, want)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	if _, err := loadConfig(filepath.Join(t.TempDir(), "absent.yml")); err == nil {
		t.Error("loadConfig() expected an error for a missing file")
	}
}

// Guards the shipped stats.yml: it must parse, and the shared ignore list must
// really reach every card.
func TestShippedConfigParses(t *testing.T) {
	cfg, err := loadConfig(filepath.Join("..", "..", "stats.yml"))
	if err != nil {
		t.Fatalf("stats.yml does not parse: %v", err)
	}

	if len(cfg.Cards) < 2 {
		t.Fatalf("stats.yml has %d cards, want at least 2", len(cfg.Cards))
	}

	first := cfg.Cards[0]
	for _, card := range cfg.Cards[1:] {
		if !reflect.DeepEqual(card.Ignore, first.Ignore) {
			t.Errorf("%s ignore list = %v, want the same as %s (%v)", card.Output, card.Ignore, first.Output, first.Ignore)
		}
	}

	if len(first.Ignore) == 0 {
		t.Error("stats.yml cards have an empty ignore list")
	}
}

func TestApplyFlagOverrides(t *testing.T) {
	base := cardConfig{
		Output:   "from-yml.svg",
		Theme:    "light",
		Header:   "From YAML",
		Mode:     "bytes",
		MaxLangs: 6,
		Ignore:   []string{"HTML"},
	}

	tests := []struct {
		name  string
		flags flagOverrides
		want  cardConfig
	}{
		{
			name:  "no flags keeps the file",
			flags: flagOverrides{set: map[string]bool{}},
			want:  base,
		},
		{
			name: "flags win over the file",
			flags: flagOverrides{
				set:      map[string]bool{"output": true, "header": true, "mode": true, "max_langs": true, "ignore": true},
				output:   "from-flag.svg",
				header:   "From Flags",
				mode:     "geometric",
				maxLangs: 3,
				ignore:   []string{"HTML", "CSS"},
			},
			want: cardConfig{
				Output:   "from-flag.svg",
				Theme:    base.Theme, // not passed as a flag
				Header:   "From Flags",
				Mode:     "geometric",
				MaxLangs: 3,
				Ignore:   []string{"HTML", "CSS"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := applyFlagOverrides(config{Cards: []cardConfig{base}}, tt.flags).Cards[0]

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyFlagOverrides() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
