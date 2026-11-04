package memory

import (
	"math/rand/v2"
	"strings"

	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
)

type quoteDAO[Q any] struct {
	quotes  map[language.Language][]Q
	indexer *indexer
	game    game.Game[Q]
}

const concurrentMatchThreshold = 5000

func NewQuotes[Q any](g game.Game[Q], quotes map[language.Language][]Q) *quoteDAO[Q] {
	fields := make(map[language.Language]quoteFields, len(quotes))
	for lang, qs := range quotes {
		n := len(qs)
		f := quoteFields{
			Texts:        make([]string, n),
			CharacterIDs: make([]string, n),
			Episodes:     make([]int, n),
			AudioIDs:     make([]string, n),
		}
		for i := range qs {
			bq := g.Base(&qs[i])
			f.Texts[i] = bq.Text
			f.CharacterIDs[i] = bq.CharacterID
			f.Episodes[i] = bq.Episode
			f.AudioIDs[i] = bq.AudioID
		}
		fields[lang] = f
	}

	return &quoteDAO[Q]{
		quotes:  quotes,
		indexer: newIndexer(fields),
		game:    g,
	}
}

func (d *quoteDAO[Q]) FindQuotes(f params.QuoteFilter, limit int, offset int) ([]Q, int, error) {
	quotes := d.quotes[f.Lang]
	matches := d.matchingIndices(f)
	total := len(matches)

	start := min(offset, total)
	end := min(start+limit, total)

	page := make([]Q, end-start)
	for i, idx := range matches[start:end] {
		page[i] = quotes[idx]
	}

	return page, total, nil
}

func (d *quoteDAO[Q]) RandomQuote(f params.QuoteFilter) (*Q, error) {
	matches := d.matchingIndices(f)
	if len(matches) == 0 {
		return nil, nil
	}

	q := d.quotes[f.Lang][matches[rand.IntN(len(matches))]]
	return &q, nil
}

func (d *quoteDAO[Q]) QuoteIndexByAudioID(lang language.Language, audioID string) (int, bool, error) {
	idx, ok := d.indexer.QuoteIndex(lang, audioID)
	return idx, ok, nil
}

func (d *quoteDAO[Q]) QuoteByIndex(lang language.Language, index int) (*Q, error) {
	quotes := d.quotes[lang]
	if index < 0 || index >= len(quotes) {
		return nil, nil
	}

	q := quotes[index]
	return &q, nil
}

func (d *quoteDAO[Q]) QuotesBetween(lang language.Language, from int, to int) ([]Q, error) {
	quotes := d.quotes[lang]

	from = max(from, 0)
	to = min(to, len(quotes)-1)
	if from > to {
		return []Q{}, nil
	}

	between := make([]Q, to-from+1)
	copy(between, quotes[from:to+1])

	return between, nil
}

func (d *quoteDAO[Q]) NearestVoiced(lang language.Language, index int, previous bool) (string, error) {
	quotes := d.quotes[lang]

	step := 1
	if previous {
		step = -1
	}

	for i := index + step; i >= 0 && i < len(quotes); i += step {
		if audioID := d.game.Base(&quotes[i]).AudioID; audioID != "" {
			return audioID, nil
		}
	}

	return "", nil
}

func (d *quoteDAO[Q]) LanguageCounts() (map[language.Language]int, error) {
	counts := make(map[language.Language]int, len(d.quotes))
	for lang, quotes := range d.quotes {
		counts[lang] = len(quotes)
	}
	return counts, nil
}

func (d *quoteDAO[Q]) Speakers(lang language.Language) (map[string]string, error) {
	speakers := map[string]string{}

	for _, q := range d.quotes[lang] {
		bq := d.game.Base(&q)
		if _, exists := speakers[bq.CharacterID]; !exists {
			speakers[bq.CharacterID] = bq.Character
		}
	}

	return speakers, nil
}

func (d *quoteDAO[Q]) SpeakerTallies(lang language.Language) (params.SpeakerTallies, error) {
	tallies := params.NewSpeakerTallies()

	var prevCharID string
	var prevSegment string

	for _, q := range d.quotes[lang] {
		bq := d.game.Base(&q)
		segment := d.game.Row(&q).Segment

		if bq.CharacterID == game.NarratorID {
			prevCharID = ""
			continue
		}

		tallies.AddLines(bq.CharacterID, segment, 1)

		if prevCharID != "" && bq.CharacterID != "" && prevCharID != bq.CharacterID && prevSegment == segment {
			tallies.AddInteractions(segment, prevCharID, bq.CharacterID, 1)
		}

		prevCharID = bq.CharacterID
		prevSegment = segment
	}

	return tallies, nil
}

func (d *quoteDAO[Q]) matchingIndices(f params.QuoteFilter) []int {
	quotes := d.quotes[f.Lang]
	if len(quotes) == 0 {
		return nil
	}

	excludeNarrator := f.ExcludeNarrator
	var source []int

	switch {
	case f.CharacterID != "" || f.Episode > 0:
		source = d.indexer.FilteredIndices(f.Lang, f.CharacterID, f.Episode)
	case excludeNarrator:
		source = d.indexer.NonNarratorIndices(f.Lang)
		excludeNarrator = false
	default:
		source = make([]int, len(quotes))
		for i := range source {
			source[i] = i
		}
	}

	byPair := f.InteractionA != "" && f.InteractionB != ""
	byScope := f.QuoteScope != params.QuoteScope{}

	if f.Text == "" && !byScope && !excludeNarrator && !byPair {
		return source
	}

	keep := func(idx int) bool {
		q := &quotes[idx]
		if excludeNarrator && d.game.Base(q).CharacterID == game.NarratorID {
			return false
		}
		if byPair && !d.interacts(quotes, idx, f.InteractionA, f.InteractionB) {
			return false
		}
		return !byScope || matchesScope(d.game.Row(q), f.QuoteScope)
	}

	if f.Text == "" {
		matches := make([]int, 0, len(source))
		for _, idx := range source {
			if keep(idx) {
				matches = append(matches, idx)
			}
		}
		return matches
	}

	lowerTexts := d.indexer.LowerTexts(f.Lang)
	queryLower := strings.ToLower(f.Text)

	if len(source) > concurrentMatchThreshold {
		return concurrentMatch(source, lowerTexts, queryLower, f.Exact, keep)
	}

	var matches []int
	for _, idx := range source {
		if matchText(lowerTexts[idx], queryLower, f.Exact) && keep(idx) {
			matches = append(matches, idx)
		}
	}
	return matches
}

func (d *quoteDAO[Q]) interacts(quotes []Q, idx int, charA string, charB string) bool {
	speaker := d.game.Base(&quotes[idx])
	if !isPairSpeaker(speaker.CharacterID, charA, charB) {
		return false
	}

	for _, n := range []int{idx - 1, idx + 1} {
		if n < 0 || n >= len(quotes) {
			continue
		}
		neighbour := d.game.Base(&quotes[n])
		if neighbour.Episode == speaker.Episode && neighbour.CharacterID != speaker.CharacterID && isPairSpeaker(neighbour.CharacterID, charA, charB) {
			return true
		}
	}

	return false
}

func isPairSpeaker(characterID string, charA string, charB string) bool {
	return characterID != "" && characterID != game.NarratorID && (characterID == charA || characterID == charB)
}

func matchesScope(row game.QuoteRow, scope params.QuoteScope) bool {
	if scope.Segment != "" && row.Segment != scope.Segment {
		return false
	}
	return (!scope.RedTruth || row.HasRedTruth) &&
		(!scope.BlueTruth || row.HasBlueTruth) &&
		(!scope.GoldTruth || row.HasGoldTruth) &&
		(!scope.PurpleTruth || row.HasPurpleTruth)
}
