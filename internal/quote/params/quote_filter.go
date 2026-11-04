package params

import "umineko_quote/internal/quote/language"

type (
	QuoteScope struct {
		Segment     string
		RedTruth    bool
		BlueTruth   bool
		GoldTruth   bool
		PurpleTruth bool
	}

	QuoteFilter struct {
		QuoteScope
		Lang            language.Language
		Text            string
		Exact           bool
		CharacterID     string
		Episode         int
		InteractionA    string
		InteractionB    string
		ExcludeNarrator bool
	}
)
