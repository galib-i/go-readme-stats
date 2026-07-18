package svg

import "testing"

func TestGetTheme_NonExistingTheme(t *testing.T) {
	theme := GetTheme("nonexistent")
	expected := themes[DefaultTheme]

	if theme != expected {
		t.Errorf("GetTheme('nonexistent') = %+v, want default theme %+v", theme, expected)
	}
}

func TestGetTheme_EmptyString(t *testing.T) {
	theme := GetTheme("")
	expected := themes[DefaultTheme]

	if theme != expected {
		t.Errorf("GetTheme('') = %+v, want default theme %+v", theme, expected)
	}
}

func TestDefaultThemeExists(t *testing.T) {
	if _, exists := themes[DefaultTheme]; !exists {
		t.Errorf("Default theme '%s' doesn't exist in themes map", DefaultTheme)
	}
}

func TestAllPredefinedThemes(t *testing.T) {
	for name := range themes {
		t.Run(name, func(t *testing.T) {
			theme := themes[name]

			if theme.Background == "" {
				t.Errorf("'%s' has empty Background", name)
			}

			if theme.Border == "" {
				t.Errorf("'%s' has empty Border", name)
			}

			if theme.Text == "" {
				t.Errorf("'%s' has empty Text", name)
			}
			if theme.SecondaryText == "" {
				t.Errorf("'%s' has empty SecondaryText", name)
			}
		})
	}
}
