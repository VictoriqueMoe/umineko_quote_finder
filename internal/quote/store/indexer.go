package store

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"umineko_quote/internal/quote/language"
)

var subtitleAudioIDSuffix = regexp.MustCompile(`_s\d+$`)

type (
	QuoteFields struct {
		Texts        []string
		CharacterIDs []string
		Episodes     []int
		AudioIDs     []string
	}

	Indexer interface {
		LowerTexts(lang language.Language) []string
		FilteredIndices(lang language.Language, characterID string, episode int) []int
		InteractionIndices(lang language.Language, interactionA string, interactionB string) []int
		CharacterIndices(lang language.Language, characterID string) []int
		NonNarratorIndices(lang language.Language) []int
		AudioFilePath(characterId string, audioId string) string
		QuoteIndex(lang language.Language, audioID string) (int, bool)
		HasAudio() bool
	}

	indexer struct {
		quoteLowerTexts  map[language.Language][]string
		characterIndex   map[language.Language]map[string][]int
		episodeIndex     map[language.Language]map[int][]int
		interactionIndex map[language.Language]map[string][]int
		nonNarratorIndex map[language.Language][]int
		audioIndex       map[language.Language]map[string]int
		episodeArrays    map[language.Language][]int
		audioDir         string
		hasAudio         bool
	}

	langIndexResult struct {
		lang           language.Language
		lowerTexts     []string
		charIdx        map[string][]int
		epIdx          map[int][]int
		interactionIdx map[string][]int
		nonNarratorIdx []int
		audioIdx       map[string]int
		episodes       []int
	}
)

func NewIndexer(fields map[language.Language]QuoteFields, audioDir string) Indexer {
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

			for i := 0; i < n; i++ {
				lowerTexts[i] = strings.ToLower(f.Texts[i])
				charIdx[f.CharacterIDs[i]] = append(charIdx[f.CharacterIDs[i]], i)
				if f.Episodes[i] > 0 {
					epIdx[f.Episodes[i]] = append(epIdx[f.Episodes[i]], i)
				}
				if f.CharacterIDs[i] != "narrator" {
					nonNarratorIdx = append(nonNarratorIdx, i)
				}
				if f.AudioIDs[i] != "" {
					for _, id := range strings.Split(f.AudioIDs[i], ", ") {
						audioIdx[id] = i
					}
				}
			}

			interactionIdx := buildInteractionIndex(f.CharacterIDs, f.Episodes)

			results <- langIndexResult{
				lang:           lang,
				lowerTexts:     lowerTexts,
				charIdx:        charIdx,
				epIdx:          epIdx,
				interactionIdx: interactionIdx,
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

	hasAudio := false
	if audioDir != "" {
		if entries, err := os.ReadDir(audioDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					subPath := filepath.Join(audioDir, entry.Name())
					if files, err := os.ReadDir(subPath); err == nil && len(files) > 0 {
						hasAudio = true
						break
					}
				}
			}
		}
	}

	idx := &indexer{
		quoteLowerTexts:  make(map[language.Language][]string),
		characterIndex:   make(map[language.Language]map[string][]int),
		episodeIndex:     make(map[language.Language]map[int][]int),
		interactionIndex: make(map[language.Language]map[string][]int),
		nonNarratorIndex: make(map[language.Language][]int),
		audioIndex:       make(map[language.Language]map[string]int),
		episodeArrays:    make(map[language.Language][]int),
		audioDir:         audioDir,
		hasAudio:         hasAudio,
	}

	for r := range results {
		idx.quoteLowerTexts[r.lang] = r.lowerTexts
		idx.characterIndex[r.lang] = r.charIdx
		idx.episodeIndex[r.lang] = r.epIdx
		idx.interactionIndex[r.lang] = r.interactionIdx
		idx.nonNarratorIndex[r.lang] = r.nonNarratorIdx
		idx.audioIndex[r.lang] = r.audioIdx
		idx.episodeArrays[r.lang] = r.episodes
	}

	return idx
}

func (idx *indexer) LowerTexts(lang language.Language) []string {
	return idx.quoteLowerTexts[lang]
}

func (idx *indexer) CharacterIndices(lang language.Language, characterID string) []int {
	langCharIdx := idx.characterIndex[lang]
	if langCharIdx == nil {
		return nil
	}
	return langCharIdx[characterID]
}

func (idx *indexer) InteractionIndices(lang language.Language, interactionA string, interactionB string) []int {
	interactionIdx := idx.interactionIndex[lang]
	if interactionIdx == nil {
		return []int{}
	}
	key, ok := InteractionPairKey(interactionA, interactionB)
	if !ok {
		return []int{}
	}
	indices, exists := interactionIdx[key]
	if !exists {
		return []int{}
	}
	return indices
}

func (idx *indexer) AudioFilePath(characterId string, audioId string) string {
	if idx.audioDir == "" {
		return ""
	}
	path := filepath.Join(idx.audioDir, characterId, audioId+".ogg")
	if _, err := os.Stat(path); err != nil {
		base := subtitleAudioIDSuffix.ReplaceAllString(audioId, "")
		if base != audioId {
			path = filepath.Join(idx.audioDir, characterId, base+".ogg")
			if _, err := os.Stat(path); err != nil {
				return ""
			}
			return path
		}
		return ""
	}
	return path
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
		langCharIdx := idx.characterIndex[lang]
		if langCharIdx == nil {
			return []int{}
		}
		return langCharIdx[characterID]
	}

	if !hasChar && hasEp {
		langEpIdx := idx.episodeIndex[lang]
		if langEpIdx == nil {
			return []int{}
		}
		return langEpIdx[episode]
	}

	langCharIdx := idx.characterIndex[lang]
	if langCharIdx == nil {
		return []int{}
	}
	charIndices := langCharIdx[characterID]
	episodes := idx.episodeArrays[lang]

	var result []int
	for _, i := range charIndices {
		if episodes[i] == episode {
			result = append(result, i)
		}
	}
	return result
}

func (idx *indexer) HasAudio() bool {
	return idx.hasAudio
}
