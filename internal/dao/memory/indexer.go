package memory

import (
	"strings"
	"sync"

	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/language"
)

type (
	quoteFields struct {
		Texts        []string
		CharacterIDs []string
		Episodes     []int
		AudioIDs     []string
	}

	indexer struct {
		quoteLowerTexts  map[language.Language][]string
		characterIndex   map[language.Language]map[string][]int
		episodeIndex     map[language.Language]map[int][]int
		nonNarratorIndex map[language.Language][]int
		audioIndex       map[language.Language]map[string]int
		episodeArrays    map[language.Language][]int
	}

	langIndexResult struct {
		lang           language.Language
		lowerTexts     []string
		charIdx        map[string][]int
		epIdx          map[int][]int
		nonNarratorIdx []int
		audioIdx       map[string]int
		episodes       []int
	}
)

func newIndexer(fields map[language.Language]quoteFields) *indexer {
	results := make(chan langIndexResult, len(fields))
	var wg sync.WaitGroup

	for lang, f := range fields {
		wg.Go(func() {
			n := len(f.Texts)
			lowerTexts := make([]string, n)
			charIdx := make(map[string][]int)
			epIdx := make(map[int][]int)
			audioIdx := make(map[string]int)
			var nonNarratorIdx []int

			for i := range n {
				lowerTexts[i] = strings.ToLower(f.Texts[i])
				charIdx[f.CharacterIDs[i]] = append(charIdx[f.CharacterIDs[i]], i)
				if f.Episodes[i] > 0 {
					epIdx[f.Episodes[i]] = append(epIdx[f.Episodes[i]], i)
				}
				if f.CharacterIDs[i] != game.NarratorID {
					nonNarratorIdx = append(nonNarratorIdx, i)
				}
				if f.AudioIDs[i] != "" {
					for id := range strings.SplitSeq(f.AudioIDs[i], ", ") {
						audioIdx[id] = i
					}
				}
			}

			results <- langIndexResult{
				lang:           lang,
				lowerTexts:     lowerTexts,
				charIdx:        charIdx,
				epIdx:          epIdx,
				nonNarratorIdx: nonNarratorIdx,
				audioIdx:       audioIdx,
				episodes:       f.Episodes,
			}
		})
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	idx := &indexer{
		quoteLowerTexts:  make(map[language.Language][]string),
		characterIndex:   make(map[language.Language]map[string][]int),
		episodeIndex:     make(map[language.Language]map[int][]int),
		nonNarratorIndex: make(map[language.Language][]int),
		audioIndex:       make(map[language.Language]map[string]int),
		episodeArrays:    make(map[language.Language][]int),
	}

	for r := range results {
		idx.quoteLowerTexts[r.lang] = r.lowerTexts
		idx.characterIndex[r.lang] = r.charIdx
		idx.episodeIndex[r.lang] = r.epIdx
		idx.nonNarratorIndex[r.lang] = r.nonNarratorIdx
		idx.audioIndex[r.lang] = r.audioIdx
		idx.episodeArrays[r.lang] = r.episodes
	}

	return idx
}

func (idx *indexer) LowerTexts(lang language.Language) []string {
	return idx.quoteLowerTexts[lang]
}

func (idx *indexer) NonNarratorIndices(lang language.Language) []int {
	return idx.nonNarratorIndex[lang]
}

func (idx *indexer) QuoteIndex(lang language.Language, audioID string) (int, bool) {
	langAudioIdx := idx.audioIndex[lang]
	if langAudioIdx == nil {
		return 0, false
	}
	i, ok := langAudioIdx[audioID]
	return i, ok
}

func (idx *indexer) FilteredIndices(lang language.Language, characterID string, episode int) []int {
	hasChar := characterID != ""
	hasEp := episode > 0

	if !hasChar && !hasEp {
		return nil
	}

	if hasChar && !hasEp {
		return idx.characterIndex[lang][characterID]
	}

	if !hasChar {
		return idx.episodeIndex[lang][episode]
	}

	episodes := idx.episodeArrays[lang]

	var result []int
	for _, i := range idx.characterIndex[lang][characterID] {
		if episodes[i] == episode {
			result = append(result, i)
		}
	}
	return result
}
