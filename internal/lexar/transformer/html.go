package transformer

import (
	"html"
	"strings"

	"umineko_quote/internal/lexar/ast"
)

// HtmlTransformer converts dialogue elements to HTML.
type HtmlTransformer struct {
	presets *PresetContext
}

// NewHtmlTransformer creates a new HtmlTransformer with the given preset context.
func NewHtmlTransformer(presets *PresetContext) *HtmlTransformer {
	return &HtmlTransformer{presets: presets}
}

// Transform converts dialogue elements to HTML.
func (t *HtmlTransformer) Transform(elements []ast.DialogueElement) string {
	var sb strings.Builder
	t.collect(&sb, elements)
	text := sb.String()

	text = strings.TrimSpace(text)
	text = strings.Trim(text, "`\"")
	text = strings.TrimPrefix(text, "&#34;")
	text = strings.TrimSuffix(text, "&#34;")
	text = strings.ReplaceAll(text, "{", "")
	text = strings.ReplaceAll(text, "}", "")
	text = strings.TrimSpace(text)

	return text
}

func (t *HtmlTransformer) collect(sb *strings.Builder, elements []ast.DialogueElement) {
	for _, elem := range elements {
		switch el := elem.(type) {
		case *ast.PlainText:
			sb.WriteString(html.EscapeString(el.Text))

		case *ast.FormatTag:
			t.writeFormatTag(sb, el)

		case *ast.SpecialChar:
			switch el.Name {
			case "n":
				sb.WriteString("<br>")
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

func (t *HtmlTransformer) writeFormatTag(sb *strings.Builder, tag *ast.FormatTag) {
	switch tag.Name {
	case "y":
		return

	case "n":
		t.collect(sb, tag.Content)

	case "i", "italic":
		sb.WriteString("<em>")
		t.collect(sb, tag.Content)
		sb.WriteString("</em>")

	case "c", "color", "colour":
		sb.WriteString(`<span style="color:#`)
		sb.WriteString(tag.Param)
		sb.WriteString(`">`)
		t.collect(sb, tag.Content)
		sb.WriteString("</span>")

	case "ruby", "h":
		sb.WriteString("<ruby>")
		t.collect(sb, tag.Content)
		sb.WriteString("<rp>(</rp><rt>")
		sb.WriteString(html.EscapeString(tag.Param))
		sb.WriteString("</rt><rp>)</rp></ruby>")

	case "p", "preset":
		if class := t.presets.GetSemanticClass(tag.Param); class != "" {
			sb.WriteString(`<span class="`)
			sb.WriteString(class)
			sb.WriteString(`">`)
			t.collect(sb, tag.Content)
			sb.WriteString("</span>")
			return
		}

		if colour := t.presets.GetDynamicColour(tag.Param); colour != "" {
			sb.WriteString(`<span style="color:`)
			sb.WriteString(colour)
			sb.WriteString(`">`)
			t.collect(sb, tag.Content)
			sb.WriteString("</span>")
			return
		}

		t.collect(sb, tag.Content)

	default:
		t.collect(sb, tag.Content)
	}
}
