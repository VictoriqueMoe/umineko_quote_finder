package lexar

import (
	"umineko_quote/internal/lexar/ast"
	"umineko_quote/internal/lexar/transformer"
)

// TruthType indicates the type of truth used in a quote.
type TruthType int

const (
	TruthNone TruthType = iota
	TruthRed
	TruthBlue
)

// String returns the string representation of a TruthType.
func (t TruthType) String() string {
	switch t {
	case TruthRed:
		return "red"
	case TruthBlue:
		return "blue"
	default:
		return ""
	}
}

// DetectTruthType examines dialogue elements and returns the truth type.
// Red truth takes precedence if both are present.
func DetectTruthType(elements []ast.DialogueElement, presets *transformer.PresetContext) TruthType {
	hasRed := false
	hasBlue := false

	detectInElements(elements, presets, &hasRed, &hasBlue)

	if hasRed {
		return TruthRed
	}
	if hasBlue {
		return TruthBlue
	}
	return TruthNone
}

func detectInElements(elements []ast.DialogueElement, presets *transformer.PresetContext, hasRed, hasBlue *bool) {
	for _, elem := range elements {
		if tag, ok := elem.(*ast.FormatTag); ok {
			if tag.Name == "p" || tag.Name == "preset" {
				class := presets.GetSemanticClass(tag.Param)
				if class == "red-truth" {
					*hasRed = true
				} else if class == "blue-truth" {
					*hasBlue = true
				}
			}
			detectInElements(tag.Content, presets, hasRed, hasBlue)
		}
	}
}
