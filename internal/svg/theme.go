package svg

const DefaultTheme = "dark"

type Theme struct {
	Background    string
	Border        string
	Text          string
	SecondaryText string
}

// themes are the built-in colour schemes.
var themes = map[string]Theme{
	"dark": {
		Background:    "#0D1117",
		Border:        "#2F353D",
		Text:          "#F0F6FC",
		SecondaryText: "#9198A1",
	},
	"soft-dark": {
		Background:    "#212830",
		Border:        "#353C44",
		Text:          "#D1D7E0",
		SecondaryText: "#9198A1",
	},
	"light": {
		Background:    "#FFFFFF",
		Border:        "#DFE4E9",
		Text:          "#1F2328",
		SecondaryText: "#59636E",
	},
}

// ThemeExists reports whether a theme with the given name is defined.
func ThemeExists(name string) bool {
	_, exists := themes[name]
	return exists
}

// GetTheme returns the theme for the given name, falling back to DefaultTheme.
func GetTheme(name string) Theme {
	if theme, exists := themes[name]; exists {
		return theme
	}

	return themes[DefaultTheme]
}
