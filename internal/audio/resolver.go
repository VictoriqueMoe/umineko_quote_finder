package audio

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type Resolver struct {
	dir      string
	hasAudio bool
}

var subtitleAudioIDSuffix = regexp.MustCompile(`_s\d+$`)

func NewResolver(audioDir string) (Resolver, error) {
	if audioDir == "" {
		return Resolver{}, nil
	}

	entries, err := os.ReadDir(audioDir)
	if errors.Is(err, fs.ErrNotExist) {
		return Resolver{dir: audioDir}, nil
	}
	if err != nil {
		return Resolver{}, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		files, err := os.ReadDir(filepath.Join(audioDir, entry.Name()))
		if err != nil {
			return Resolver{}, err
		}
		if len(files) > 0 {
			return Resolver{dir: audioDir, hasAudio: true}, nil
		}
	}

	return Resolver{dir: audioDir}, nil
}

func (r Resolver) AudioFilePath(characterId string, audioId string) string {
	if r.dir == "" {
		return ""
	}
	path := filepath.Join(r.dir, characterId, audioId+".ogg")
	if _, err := os.Stat(path); err != nil {
		base := subtitleAudioIDSuffix.ReplaceAllString(audioId, "")
		if base != audioId {
			path = filepath.Join(r.dir, characterId, base+".ogg")
			if _, err := os.Stat(path); err != nil {
				return ""
			}
			return path
		}
		return ""
	}
	return path
}

func (r Resolver) HasAudio() bool {
	return r.hasAudio
}
