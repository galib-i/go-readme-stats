package handler

import (
	_ "embed"
	"go-readme-stats/app/stats"
	"go-readme-stats/app/svg"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	FetchStats  = stats.FetchStats
	GenerateSVG = svg.Generate
)

//go:embed ignored_languages.json
var ignoredLanguages []byte

func Favicon(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func GetLanguageStats(c *gin.Context) {
	theme := c.DefaultQuery("theme", svg.DefaultTheme)
	header := c.DefaultQuery("header", "Languages")
	mode := c.DefaultQuery("mode", "bytes")

	languages, err := FetchStats(c.Request.Context(), ignoredLanguages, mode)
	if err != nil {
		log.Printf("Error: Failed to fetch stats for request %s: %v", c.Request.URL.String(), err)
		c.String(http.StatusInternalServerError, "Error fetching stats")
		return
	}

	svgContent, err := GenerateSVG(theme, header, languages)
	if err != nil {
		log.Printf("Error: Failed to generate SVG for request %s: %v", c.Request.URL.String(), err)
		c.String(http.StatusInternalServerError, "Error generating SVG")
		return
	}

	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svgContent)
}
