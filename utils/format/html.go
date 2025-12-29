package format

import (
	"html"
	"regexp"
	"strings"
)

func SanitizeHTML(htmlContent string) string {
	htmlContent = removeDangerousTags(htmlContent)
	htmlContent = removeEventHandlers(htmlContent)
	htmlContent = removeJavascriptProtocol(htmlContent)
	htmlContent = sanitizeStyles(htmlContent)
	return htmlContent
}

func removeDangerousTags(html string) string {
	dangerousTags := []string{
		"script", "iframe", "object", "embed", "applet",
		"meta", "link", "base", "form", "input", "button",
	}

	for _, tag := range dangerousTags {
		regex := regexp.MustCompile(`(?i)<` + tag + `[^>]*>[\s\S]*?</` + tag + `>`)
		html = regex.ReplaceAllString(html, "")
		regex = regexp.MustCompile(`(?i)<` + tag + `[^>]*>`)
		html = regex.ReplaceAllString(html, "")
	}

	return html
}

func removeEventHandlers(html string) string {
	eventHandlers := regexp.MustCompile(`(?i)\s*on\w+\s*=\s*["'][^"']*["']`)
	return eventHandlers.ReplaceAllString(html, "")
}

func removeJavascriptProtocol(html string) string {
	jsProtocol := regexp.MustCompile(`(?i)javascript:`)
	return jsProtocol.ReplaceAllString(html, "")
}

func sanitizeStyles(html string) string {
	dangerousStyles := []string{"behavior", "expression", "binding", "import", "moz-binding"}

	for _, style := range dangerousStyles {
		regex := regexp.MustCompile(`(?i)` + style + `\s*:\s*[^;]+;?`)
		html = regex.ReplaceAllString(html, "")
	}

	return html
}

func GenerateSnippet(bodyText, bodyHTML string) string {
	text := bodyText
	if text == "" && bodyHTML != "" {
		text = StripHTML(bodyHTML)
	}

	text = strings.TrimSpace(text)
	if len(text) > 150 {
		text = text[:150] + "..."
	}

	return text
}

func StripHTML(html string) string {
	text := html

	styleRegex := regexp.MustCompile(`(?i)<style[^>]*>[\s\S]*?</style>`)
	text = styleRegex.ReplaceAllString(text, "")

	scriptRegex := regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`)
	text = scriptRegex.ReplaceAllString(text, "")

	headRegex := regexp.MustCompile(`(?i)<head[^>]*>[\s\S]*?</head>`)
	text = headRegex.ReplaceAllString(text, "")

	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n\n")
	text = strings.ReplaceAll(text, "</div>", "\n")
	text = strings.ReplaceAll(text, "</tr>", "\n")
	text = strings.ReplaceAll(text, "</h1>", "\n")
	text = strings.ReplaceAll(text, "</h2>", "\n")
	text = strings.ReplaceAll(text, "</h3>", "\n")
	text = strings.ReplaceAll(text, "</li>", "\n")

	inTag := false
	var result strings.Builder
	for _, char := range text {
		if char == '<' {
			inTag = true
			continue
		}
		if char == '>' {
			inTag = false
			continue
		}
		if !inTag {
			result.WriteRune(char)
		}
	}

	cleanText := result.String()

	lines := strings.Split(cleanText, "\n")
	var cleanLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.TrimSpace(strings.Join(cleanLines, " "))
}

func HTMLToPlainText(htmlContent string) string {
	text := htmlContent

	text = regexp.MustCompile(`(?i)<style[^>]*>[\s\S]*?</style>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`(?i)<script[^>]*>[\s\S]*?</script>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`(?i)<head[^>]*>[\s\S]*?</head>`).ReplaceAllString(text, "")
	text = regexp.MustCompile(`(?i)<title[^>]*>[\s\S]*?</title>`).ReplaceAllString(text, "")

	text = regexp.MustCompile(`(?i)<br\s*/?>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)</p>`).ReplaceAllString(text, "\n\n")
	text = regexp.MustCompile(`(?i)</div>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)</tr>`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`(?i)</h[1-6]>`).ReplaceAllString(text, "\n\n")
	text = regexp.MustCompile(`(?i)</li>`).ReplaceAllString(text, "\n")

	text = regexp.MustCompile(`<[^>]+>`).ReplaceAllString(text, "")

	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&#x27;", "'")

	text = regexp.MustCompile(`\n\s*\n\s*\n+`).ReplaceAllString(text, "\n\n")
	text = regexp.MustCompile(`[ \t]+`).ReplaceAllString(text, " ")

	lines := strings.Split(text, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	return strings.TrimSpace(strings.Join(cleanLines, "\n"))
}

func DecodeHTML(text string) string {
	return html.UnescapeString(text)
}
