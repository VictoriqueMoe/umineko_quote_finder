package params

type SpeakerTallies struct {
	CharacterSegments   map[string]map[string]int
	SegmentInteractions map[string]map[string]int
}

func NewSpeakerTallies() SpeakerTallies {
	return SpeakerTallies{
		CharacterSegments:   make(map[string]map[string]int),
		SegmentInteractions: make(map[string]map[string]int),
	}
}

func (t SpeakerTallies) AddLines(characterID string, segment string, count int) {
	if t.CharacterSegments[characterID] == nil {
		t.CharacterSegments[characterID] = make(map[string]int)
	}
	t.CharacterSegments[characterID][segment] += count
}

func (t SpeakerTallies) AddInteractions(segment string, charA string, charB string, count int) {
	if charA > charB {
		charA, charB = charB, charA
	}

	if t.SegmentInteractions[segment] == nil {
		t.SegmentInteractions[segment] = make(map[string]int)
	}
	t.SegmentInteractions[segment][charA+"|"+charB] += count
}

func (t SpeakerTallies) CharacterTotals() map[string]int {
	totals := make(map[string]int, len(t.CharacterSegments))
	for id, segments := range t.CharacterSegments {
		for _, count := range segments {
			totals[id] += count
		}
	}
	return totals
}

func (t SpeakerTallies) InteractionTotals() map[string]int {
	totals := make(map[string]int)
	for _, pairs := range t.SegmentInteractions {
		for key, count := range pairs {
			totals[key] += count
		}
	}
	return totals
}
