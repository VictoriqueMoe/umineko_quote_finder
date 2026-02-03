package quote

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type (
	Indexer interface {
		LowerTexts(lang string) []string
		FilteredIndices(lang string, characterID string, episode int) []int
		CharacterIndices(lang string, characterID string) []int
		NonNarratorIndices(lang string) []int
		AudioFilePath(characterId string, audioId string) string
		QuoteIndex(lang string, audioID string) (int, bool)
		HasAudio() bool
	}

	indexer struct {
		quoteLowerTexts  map[string][]string
		characterIndex   map[string]map[string][]int
		episodeIndex     map[string]map[int][]int
		nonNarratorIndex map[string][]int
		audioIndex       map[string]map[string]int
		quotes           map[string][]ParsedQuote
		audioDir         string
		hasAudio         bool
	}

	langIndexResult struct {
		lang           string
		lowerTexts     []string
		charIdx        map[string][]int
		epIdx          map[int][]int
		nonNarratorIdx []int
		audioIdx       map[string]int
	}
)

func NewIndexer(quotes map[string][]ParsedQuote, audioDir string) Indexer {
	results := make(chan langIndexResult, len(quotes))
	var wg sync.WaitGroup

	for lang, parsed := range quotes {
		wg.Go(func() {
			lowerTexts := make([]string, len(parsed))
			charIdx := make(map[string][]int)
			epIdx := make(map[int][]int)
			audioIdx := make(map[string]int)
			var nonNarratorIdx []int

			for i := 0; i < len(parsed); i++ {
				lowerTexts[i] = strings.ToLower(parsed[i].Text)
				charIdx[parsed[i].CharacterID] = append(charIdx[parsed[i].CharacterID], i)
				if parsed[i].Episode > 0 {
					epIdx[parsed[i].Episode] = append(epIdx[parsed[i].Episode], i)
				}
				if parsed[i].CharacterID != "narrator" {
					nonNarratorIdx = append(nonNarratorIdx, i)
				}
				if parsed[i].AudioID != "" {
					for _, id := range strings.Split(parsed[i].AudioID, ", ") {
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
		quoteLowerTexts:  make(map[string][]string),
		characterIndex:   make(map[string]map[string][]int),
		episodeIndex:     make(map[string]map[int][]int),
		nonNarratorIndex: make(map[string][]int),
		audioIndex:       make(map[string]map[string]int),
		quotes:           quotes,
		audioDir:         audioDir,
		hasAudio:         hasAudio,
	}

	for r := range results {
		idx.quoteLowerTexts[r.lang] = r.lowerTexts
		idx.characterIndex[r.lang] = r.charIdx
		idx.episodeIndex[r.lang] = r.epIdx
		idx.nonNarratorIndex[r.lang] = r.nonNarratorIdx
		idx.audioIndex[r.lang] = r.audioIdx
	}

	return idx
}

func (idx *indexer) LowerTexts(lang string) []string {
	return idx.quoteLowerTexts[lang]
}

func (idx *indexer) CharacterIndices(lang string, characterID string) []int {
	langCharIdx := idx.characterIndex[lang]
	if langCharIdx == nil {
		return nil
	}
	return langCharIdx[characterID]
}

func (idx *indexer) AudioFilePath(characterId string, audioId string) string {
	if idx.audioDir == "" {
		return ""
	}
	path := filepath.Join(idx.audioDir, characterId, audioId+".ogg")
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

func (idx *indexer) NonNarratorIndices(lang string) []int {
	return idx.nonNarratorIndex[lang]
}

func (idx *indexer) QuoteIndex(lang string, audioID string) (int, bool) {
	langAudioIdx := idx.audioIndex[lang]
	if langAudioIdx == nil {
		return 0, false
	}
	i, ok := langAudioIdx[audioID]
	return i, ok
}

func (idx *indexer) FilteredIndices(lang string, characterID string, episode int) []int {
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
	quotes := idx.quotes[lang]

	var result []int
	for _, i := range charIndices {
		if quotes[i].Episode == episode {
			result = append(result, i)
		}
	}
	return result
}

func (idx *indexer) HasAudio() bool {
	return idx.hasAudio
}
