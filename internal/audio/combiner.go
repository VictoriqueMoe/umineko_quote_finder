package audio

import (
	"fmt"
	"os"
)

type FilePathResolver func(charId, audioId string) string

type AudioSegment struct {
	CharID  string
	AudioID string
}

type Combiner interface {
	CombineOgg(segments []AudioSegment, resolve FilePathResolver) ([]byte, error)
}

type combiner struct{}

func NewCombiner() (Combiner, error) {
	return &combiner{}, nil
}

func (c *combiner) CombineOgg(segments []AudioSegment, resolve FilePathResolver) ([]byte, error) {
	var allFilePages [][]oggPage

	for i := 0; i < len(segments); i++ {
		filePath := resolve(segments[i].CharID, segments[i].AudioID)
		if filePath == "" {
			return nil, fmt.Errorf("audio file not found: %s/%s", segments[i].CharID, segments[i].AudioID)
		}
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read audio file: %s", segments[i].AudioID)
		}
		pages, err := parseOggPages(data)
		if err != nil {
			return nil, fmt.Errorf("failed to parse OGG file %s: %v", segments[i].AudioID, err)
		}
		allFilePages = append(allFilePages, pages)
	}

	if len(allFilePages) == 0 {
		return nil, fmt.Errorf("no audio files to combine")
	}

	serialNumber := allFilePages[0][0].serialNumber
	var result []byte
	var sequenceNum uint32
	var granuleOffset int64

	for fileIdx := 0; fileIdx < len(allFilePages); fileIdx++ {
		pages := allFilePages[fileIdx]
		isFirst := fileIdx == 0
		isLast := fileIdx == len(allFilePages)-1

		var fileLastGranule int64
		for i := len(pages) - 1; i >= 0; i-- {
			if pages[i].granulePos > 0 {
				fileLastGranule = pages[i].granulePos
				break
			}
		}

		headersDone := isFirst
		firstIncluded := false

		for pi := 0; pi < len(pages); pi++ {
			page := pages[pi]

			if !headersDone {
				if page.headerType&0x02 != 0 || page.granulePos <= 0 {
					continue
				}
				headersDone = true
			}

			modified := oggPage{
				headerType:     page.headerType,
				granulePos:     page.granulePos,
				serialNumber:   serialNumber,
				sequenceNumber: sequenceNum,
				segmentTable:   page.segmentTable,
				data:           page.data,
			}

			if !isFirst {
				modified.headerType &^= 0x02
				if !firstIncluded {
					modified.headerType &^= 0x01
					firstIncluded = true
				}
			}

			if !isFirst && page.granulePos > 0 {
				modified.granulePos = page.granulePos + granuleOffset
			}

			if !isLast {
				modified.headerType &^= 0x04
			}

			result = append(result, modified.serialize()...)
			sequenceNum++
		}

		granuleOffset += fileLastGranule
	}

	return result, nil
}
