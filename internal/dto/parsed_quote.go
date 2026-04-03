package dto

import scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"

type (
	UminekoQuote struct {
		scriptdto.UminekoQuote
		// Position index in the parsed script (stable across text changes)
		Index int `json:"index"`
	}

	HigurashiQuote struct {
		scriptdto.HigurashiQuote
		// Position index in the parsed script (stable across text changes)
		Index int `json:"index"`
	}

	CharactersResult struct {
		Characters map[string]string `json:"characters"`
		Additional map[string]string `json:"additional"`
	}
)
