package transformer

import (
	"strings"

	"umineko_quote/internal/lexar/ast"
)

// PlainTextTransformer converts dialogue elements to plain text.
type PlainTextTransformer struct{}

// NewPlainTextTransformer creates a new PlainTextTransformer.
func NewPlainTextTransformer() *PlainTextTransformer {
	return &PlainTextTransformer{}
}

// Transform converts dialogue elements to plain text.
func (t *PlainTextTransformer) Transform(elements []ast.DialogueElement) string {
	var sb strings.Builder
	t.collect(&sb, elements)
	text := sb.String()

	text = strings.TrimSpace(text)
	text = strings.Trim(text, "`\"")
	text = strings.ReplaceAll(text, "{", "")
	text = strings.ReplaceAll(text, "}", "")
	text = strings.TrimSpace(text)

	return text
}

func (t *PlainTextTransformer) collect(sb *strings.Builder, elements []ast.DialogueElement) {
	for _, elem := range elements {
		switch el := elem.(type) {
		case *ast.PlainText:
			sb.WriteString(el.Text)

		case *ast.FormatTag:
			if el.Name == "y" {
				continue
			}
			if el.Name == "ruby" || el.Name == "h" {
				t.collect(sb, el.Content)
				if el.Param != "" {
					sb.WriteString(" (")
					sb.WriteString(el.Param)
					sb.WriteString(")")
				}
				continue
			}
			t.collect(sb, el.Content)

		case *ast.SpecialChar:
			switch el.Name {
			case "n":
				sb.WriteString(" ")
			case "qt":
				sb.WriteString(`"`)
			case "os":
				sb.WriteString("[")
			case "es":
				sb.WriteString("]")
			}
		}
	}
}
