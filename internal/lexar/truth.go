package lexar

import (
	"umineko_quote/internal/lexar/ast"
	"umineko_quote/internal/lexar/transformer"
)

// TruthFlags indicates which truth types are present in a quote.
type TruthFlags struct {
	HasRed  bool
	HasBlue bool
}

// DetectTruth examines dialogue elements and returns which truth types are present.
func DetectTruth(elements []ast.DialogueElement, presets *transformer.PresetContext) TruthFlags {
	var flags TruthFlags
	detectInElements(elements, presets, &flags.HasRed, &flags.HasBlue)
	return flags
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
