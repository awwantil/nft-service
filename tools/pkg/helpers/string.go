package helpers

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/microcosm-cc/bluemonday"
)

// Truncate truncates the given content string to fit within the specified maximum length.
func Truncate(content string, maxLength int) string {
	if utf8.RuneCountInString(content) <= maxLength {
		return content
	}

	truncatedText := string([]rune(content)[:maxLength])

	lastSpaceIndex := strings.LastIndexByte(truncatedText, ' ')

	if lastSpaceIndex != -1 && lastSpaceIndex != maxLength-1 {
		truncatedText = truncatedText[:lastSpaceIndex]
	}

	return truncatedText
}

// Sanitize removes unwanted elements from the input `content` string,
func Sanitize(content string) string {
	p := bluemonday.StrictPolicy()
	rg := regexp.MustCompile(`\s{2,}\n?`)
	content = p.Sanitize(content)
	content = rg.ReplaceAllString(content, "\n")
	return strings.TrimSpace(content)
}
