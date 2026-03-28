package dto

import scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"

type (
	ScriptParsedQuote = scriptdto.ParsedQuote

	ParsedQuote struct {
		ScriptParsedQuote
		// Position index in the parsed script (stable across text changes)
		Index int `json:"index"`
	}
)
