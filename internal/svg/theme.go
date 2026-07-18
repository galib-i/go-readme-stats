package svg

const DefaultTheme = "dark"

type Theme struct {
	Background    string
	Border        string
	Text          string
	SecondaryText string
}

// themes contains predefined colour schemes for SVG generation.
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

// GetTheme returns the theme configuration for the given name.
// Falls back to DefaultTheme if the requested theme doesn't exist.
func GetTheme(name string) Theme {
	if theme, exists := themes[name]; exists {
		return theme
	}

	return themes[DefaultTheme]
}
