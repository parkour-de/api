package description

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/microcosm-cc/bluemonday"
	"strings"
)

var p *parser.Parser
var r *html.Renderer

var b *bluemonday.Policy

// FixTitle takes a title and text as input, normalises the text's first line into a h1 heading and adds
// block-separating newlines if necessary. If the text doesn't start with a heading, a provided heading is added.
func FixTitle(title string, text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) > 0 {
		firstLine := strings.TrimSpace(lines[0])
		if !strings.HasPrefix(firstLine, "#") {
			return "# " + strings.TrimSpace(title) + "\n\n" + text
		} else {
			for strings.HasPrefix(firstLine, "#") {
				firstLine = strings.TrimPrefix(firstLine, "#")
			}
			firstLine = "# " + strings.TrimSpace(firstLine)
			lines[0] = firstLine
		}
	}
	if len(lines) > 1 && len(lines[1]) > 0 {
		lines[1] = "\n" + lines[1]
	} else {
		lines[0] += "\n\n"
	}
	return strings.Join(lines, "\n")
}

// GetTitle extracts the title from a Markdown text's first line, removing the "# " prefix and trimming spaces.
func GetTitle(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) > 0 {
		line := strings.TrimSpace(lines[0])
		if strings.HasPrefix(line, "# ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "# "))
		}
	}
	return ""
}

func Render(md []byte) string {
	// if p == nil {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p = parser.NewWithExtensions(extensions)
	// }

	// parser can not be recycled because p.tip is nil on second execution.
	doc := p.Parse(md)

	if r == nil {
		htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.NofollowLinks | html.NoreferrerLinks | html.NoopenerLinks | html.LazyLoadImages
		opts := html.RendererOptions{Flags: htmlFlags}
		r = html.NewRenderer(opts)
	}

	if b == nil {
		b = bluemonday.UGCPolicy().AddTargetBlankToFullyQualifiedLinks(true)
	}

	html := markdown.Render(doc, r)

	return b.Sanitize(string(html))
}
