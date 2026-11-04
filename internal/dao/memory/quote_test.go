package memory

import (
	"reflect"
	"testing"

	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

func ciconiaQuote(charID string, chapter string) dto.CiconiaQuote {
	return dto.CiconiaQuote{
		CicroniaQuote: scriptdto.CicroniaQuote{
			BaseQuote: scriptdto.BaseQuote{
				CharacterID: charID,
				Character:   charID,
				Episode:     1,
			},
			Chapter: chapter,
		},
	}
}

func uminekoQuote(charID string, episode int) dto.UminekoQuote {
	return dto.UminekoQuote{
		UminekoQuote: scriptdto.UminekoQuote{
			BaseQuote: scriptdto.BaseQuote{
				CharacterID: charID,
				Episode:     episode,
			},
		},
	}
}

func englishCiconia(quotes ...dto.CiconiaQuote) *quoteDAO[dto.CiconiaQuote] {
	return NewQuotes(game.Ciconia{}, map[language.Language][]dto.CiconiaQuote{language.English: quotes})
}

func tallies(t *testing.T, d *quoteDAO[dto.CiconiaQuote]) params.SpeakerTallies {
	t.Helper()

	got, err := d.SpeakerTallies(language.English)
	if err != nil {
		t.Fatalf("SpeakerTallies: %v", err)
	}
	return got
}

func TestSpeakerTallies_CountsLinesPerSegment(t *testing.T) {
	got := tallies(t, englishCiconia(
		ciconiaQuote("miyao", "01"),
		ciconiaQuote("miyao", "01"),
		ciconiaQuote("jayden", "01"),
		ciconiaQuote("miyao", "02"),
	))

	want := map[string]map[string]int{
		"miyao":  {"01": 2, "02": 1},
		"jayden": {"01": 1},
	}
	if !reflect.DeepEqual(want, got.CharacterSegments) {
		t.Errorf("character segments: want %v, got %v", want, got.CharacterSegments)
	}
}

func TestSpeakerTallies_NarratorExcluded(t *testing.T) {
	got := tallies(t, englishCiconia(
		ciconiaQuote(game.NarratorID, "01"),
		ciconiaQuote("miyao", "01"),
	))

	if _, ok := got.CharacterSegments[game.NarratorID]; ok {
		t.Error("narrator should not be tallied")
	}
}

func TestSpeakerTallies_NarratorBreaksInteraction(t *testing.T) {
	got := tallies(t, englishCiconia(
		ciconiaQuote("miyao", "01"),
		ciconiaQuote(game.NarratorID, "01"),
		ciconiaQuote("jayden", "01"),
	))

	if len(got.SegmentInteractions) != 0 {
		t.Errorf("expected no interactions across a narrator line, got %v", got.SegmentInteractions)
	}
}

func TestSpeakerTallies_InteractionsBreakAcrossSegments(t *testing.T) {
	got := tallies(t, englishCiconia(
		ciconiaQuote("miyao", "01"),
		ciconiaQuote("jayden", "02"),
	))

	if len(got.SegmentInteractions) != 0 {
		t.Errorf("expected no interactions across chapters, got %v", got.SegmentInteractions)
	}
}

func TestSpeakerTallies_InteractionsKeyedBySegment(t *testing.T) {
	quotes := map[language.Language][]dto.UminekoQuote{
		language.English: {
			uminekoQuote("10", 1),
			uminekoQuote("27", 1),
			uminekoQuote("10", 1),
			uminekoQuote("10", 2),
			uminekoQuote("27", 2),
		},
	}

	got, err := NewQuotes(game.Umineko{}, quotes).SpeakerTallies(language.English)
	if err != nil {
		t.Fatalf("SpeakerTallies: %v", err)
	}

	want := map[string]map[string]int{
		"1": {"10|27": 2},
		"2": {"10|27": 1},
	}
	if !reflect.DeepEqual(want, got.SegmentInteractions) {
		t.Errorf("segment interactions: want %v, got %v", want, got.SegmentInteractions)
	}
}

func TestSpeakers_FirstSeenName(t *testing.T) {
	first := ciconiaQuote("misc", "01")
	first.Character = "First Name"
	second := ciconiaQuote("misc", "01")
	second.Character = "Second Name"

	speakers, err := englishCiconia(first, second).Speakers(language.English)
	if err != nil {
		t.Fatalf("Speakers: %v", err)
	}

	if speakers["misc"] != "First Name" {
		t.Errorf("expected first-seen name, got %q", speakers["misc"])
	}
}

func TestTruthPerEpisode(t *testing.T) {
	red := uminekoQuote(game.NarratorID, 1)
	red.HasRedTruth = true
	blue := uminekoQuote("10", 1)
	blue.HasBlueTruth = true
	gold := uminekoQuote("27", 2)
	gold.HasGoldTruth = true
	purple := uminekoQuote("27", 2)
	purple.HasPurpleTruth = true

	d := &uminekoDAO{quoteDAO: NewQuotes(game.Umineko{}, map[language.Language][]dto.UminekoQuote{
		language.English: {red, blue, gold, purple, uminekoQuote("10", 3)},
	})}

	got, err := d.TruthPerEpisode(language.English)
	if err != nil {
		t.Fatalf("TruthPerEpisode: %v", err)
	}

	want := map[int]params.TruthCounts{
		1: {Red: 1, Blue: 1},
		2: {Gold: 1, Purple: 1},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("truth per episode: want %v, got %v", want, got)
	}
}
