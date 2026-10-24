package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
	"time"

	"go-readme-stats/internal/stats"
	"go-readme-stats/internal/svg"

	"github.com/goccy/go-yaml"
)

const (
	configPath = "stats.yml"

	// runTimeout bounds the whole run so a stalled API cannot hang the workflow.
	runTimeout = 5 * time.Minute
)

type cardConfig struct {
	Output   string   `yaml:"output"`
	Theme    string   `yaml:"theme"`
	Header   string   `yaml:"header"`
	Mode     string   `yaml:"mode"`
	MaxLangs int      `yaml:"max_langs"`
	Ignore   []string `yaml:"ignore"`
}

type config struct {
	Cards []cardConfig `yaml:"cards"`
}

// flagOverrides holds values given on the command line. set records which flags
// were actually passed, so an empty value can still override stats.yml.
type flagOverrides struct {
	set      map[string]bool
	output   string
	theme    string
	header   string
	mode     string
	maxLangs int
	ignore   []string
}

func parseFlags() flagOverrides {
	output := flag.String("output", "", "output SVG file name")
	theme := flag.String("theme", "", "theme name: dark, soft-dark or light")
	header := flag.String("header", "", "card header text")
	mode := flag.String("mode", "", `percentage calculation method: "bytes" or "geometric"`)
	maxLangs := flag.Int("max_langs", 0, "maximum number of languages shown")
	ignore := flag.String("ignore", "", "comma-separated list of languages to ignore")

	flag.Parse()

	set := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { set[f.Name] = true })

	var ignored []string
	for _, name := range strings.Split(*ignore, ",") {
		if name = strings.TrimSpace(name); name != "" {
			ignored = append(ignored, name)
		}
	}

	return flagOverrides{
		set:      set,
		output:   *output,
		theme:    *theme,
		header:   *header,
		mode:     *mode,
		maxLangs: *maxLangs,
		ignore:   ignored,
	}
}

// applyFlagOverrides lets flags win over the file. Priority is: flags, then
// stats.yml, then defaults.
func applyFlagOverrides(cfg config, flags flagOverrides) config {
	for i := range cfg.Cards {
		card := &cfg.Cards[i]

		if flags.set["output"] {
			card.Output = flags.output
		}
		if flags.set["theme"] {
			card.Theme = flags.theme
		}
		if flags.set["header"] {
			card.Header = flags.header
		}
		if flags.set["mode"] {
			card.Mode = flags.mode
		}
		if flags.set["max_langs"] {
			card.MaxLangs = flags.maxLangs
		}
		if flags.set["ignore"] {
			card.Ignore = flags.ignore
		}
	}

	return cfg
}

// loadConfig reads stats.yml, filling in defaults for anything it leaves out.
func loadConfig(path string) (config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return config{}, fmt.Errorf("failed to read %s: %w", path, err)
	}

	var cfg config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return config{}, fmt.Errorf("failed to parse %s: %w", path, err)
	}

	for i := range cfg.Cards {
		card := &cfg.Cards[i]
		if card.Output == "" {
			card.Output = "languages.svg"
		}
		if card.Theme == "" {
			card.Theme = svg.DefaultTheme
		}
		if card.Header == "" {
			card.Header = "Languages"
		}
		if card.Mode == "" {
			card.Mode = "bytes"
		}
		if card.MaxLangs == 0 {
			card.MaxLangs = 6
		}
	}

	return cfg, nil
}

// warnUnknownValues flags settings that would otherwise fall back to a default
// without a word.
func warnUnknownValues(card cardConfig) {
	switch card.Mode {
	case "bytes", "geometric":
	default:
		log.Printf("Warning: %s: unknown mode %q, falling back to \"bytes\"", card.Output, card.Mode)
	}

	if !svg.ThemeExists(card.Theme) {
		log.Printf("Warning: %s: unknown theme %q, falling back to %q", card.Output, card.Theme, svg.DefaultTheme)
	}
}

// statsKey is the set of settings that decide a fetch, so cards asking for the
// same numbers can share one call.
type statsKey struct {
	mode     string
	maxLangs int
	ignore   string
}

func newStatsKey(card cardConfig) statsKey {
	ignored := slices.Clone(card.Ignore)
	slices.Sort(ignored)

	return statsKey{
		mode:     card.Mode,
		maxLangs: card.MaxLangs,
		ignore:   strings.Join(ignored, ","),
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	flags := parseFlags()

	cfg, err := loadConfig(configPath)
	if err != nil {
		return err
	}

	cfg = applyFlagOverrides(cfg, flags)

	ctx, cancel := context.WithTimeout(context.Background(), runTimeout)
	defer cancel()

	cache := make(map[statsKey][]stats.Lang)

	for _, card := range cfg.Cards {
		warnUnknownValues(card)

		key := newStatsKey(card)

		languages, fetched := cache[key]
		if !fetched {
			languages, err = stats.FetchStats(ctx, card.Ignore, card.Mode, card.MaxLangs)
			if err != nil {
				return fmt.Errorf("failed to fetch stats for %s: %w", card.Output, err)
			}

			cache[key] = languages
		}

		svgContent, err := svg.Generate(card.Theme, card.Header, languages)
		if err != nil {
			return fmt.Errorf("failed to generate %s: %w", card.Output, err)
		}

		if err = os.WriteFile(card.Output, []byte(svgContent), 0o644); err != nil {
			return fmt.Errorf("failed to save %s: %w", card.Output, err)
		}

		log.Printf("Generated %s", card.Output)
	}

	return nil
}
