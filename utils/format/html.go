package format

import (
	"regexp"
	"strings"
)

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
