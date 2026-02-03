package transformer

import "umineko_quote/internal/lexar/ast"

// Transformer converts dialogue elements into a string representation.
type Transformer interface {
	Transform(elements []ast.DialogueElement) string
}
