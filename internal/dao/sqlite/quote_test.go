package sqlite

import (
	"database/sql"
	"path/filepath"
	"reflect"
	"testing"

	"umineko_quote/internal/dao"
	"umineko_quote/internal/dao/memory"
	"umineko_quote/internal/db"
	"umineko_quote/internal/dto"
	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"
)

type engines[Q any] struct {
	memory dao.QuoteDAO[Q]
	sqlite dao.QuoteDAO[Q]
}

const testDataHash = "test-hash-v1"

func testQuotes() map[language.Language][]dto.UminekoQuote {
	mk := func(idx int, text string, charID string, audioID string, episode int) dto.UminekoQuote {
		return dto.UminekoQuote{
			Text:        text,
			CharacterID: charID,
			AudioID:     audioID,
			Episode:     episode,
			Index:       idx,
		}
	}

	en := []dto.UminekoQuote{
		mk(0, "The wind blew softly.", game.NarratorID, "", 1),
		mk(1, "Without love, it cannot be seen.", "10", "30101088", 1),
		mk(2, "Without love, the truth cannot be found.", "27", "30101089, 30101090", 1),
		mk(3, "Silence.", game.NarratorID, "", 1),
		mk(4, "I'll deny your existence, witch!", "10", "30101091", 2),
		mk(5, "Checkmate. You are already mine.", "27", "30101092", 2),
	}
	en[1].HasRedTruth = true
	en[2].HasBlueTruth = true
	en[5].HasGoldTruth = true

	ja := []dto.UminekoQuote{
		mk(0, "風が吹いた。", game.NarratorID, "", 1),
	}

	return map[language.Language][]dto.UminekoQuote{
		language.English:  en,
		language.Japanese: ja,
	}
}

func ciconiaTestQuotes() map[language.Language][]dto.CiconiaQuote {
	mk := func(idx int, charID string, charName string, chapter string) dto.CiconiaQuote {
		return dto.CiconiaQuote{
			Text:        "text",
			CharacterID: charID,
			Character:   charName,
			Episode:     1,
			Chapter:     chapter,
			Index:       idx,
		}
	}

	en := []dto.CiconiaQuote{
		mk(0, game.NarratorID, "Narrator", "01"),
		mk(1, "mitake", "Mitake Miyao", "01"),
		mk(2, "jayden", "Jayden", "01"),
		mk(3, "mitake", "Second Name", "01"),
		mk(4, "mitake", "Mitake Miyao", "02"),
		mk(5, "jayden", "Jayden", "02"),
	}

	ja := make([]dto.CiconiaQuote, len(en))
	for i, q := range en {
		q.Text = "テキスト"
		ja[i] = q
	}

	return map[language.Language][]dto.CiconiaQuote{
		language.English:  en,
		language.Japanese: ja,
	}
}

func openTestDB(t *testing.T, path string) *sql.DB {
	t.Helper()

	conn, err := db.Open(path)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}

	if err := db.Migrate(conn); err != nil {
		t.Fatalf("db.Migrate: %v", err)
	}

	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("close database: %v", err)
		}
	})

	return conn
}

func newSeededDAO[Q any](t *testing.T, g game.Game[Q], quotes map[language.Language][]Q) *quoteDAO[Q] {
	t.Helper()

	conn := openTestDB(t, filepath.Join(t.TempDir(), "test.db"))

	s := &quoteDAO[Q]{db: conn, game: g}

	needsSeed, err := s.NeedsSeed(testDataHash)
	if err != nil {
		t.Fatalf("NeedsSeed: %v", err)
	}
	if !needsSeed {
		t.Fatal("expected fresh database to need seeding")
	}

	if err := s.Seed(quotes, testDataHash); err != nil {
		t.Fatalf("Seed: %v", err)
	}

	needsSeed, err = s.NeedsSeed(testDataHash)
	if err != nil {
		t.Fatalf("NeedsSeed after seed: %v", err)
	}
	if needsSeed {
		t.Fatal("expected seeded database with matching hash to not need seeding")
	}

	return s
}

func newEngines[Q any](t *testing.T, g game.Game[Q], quotes map[language.Language][]Q) engines[Q] {
	t.Helper()

	return engines[Q]{
		memory: memory.NewQuotes(g, quotes),
		sqlite: newSeededDAO(t, g, quotes),
	}
}

func indexes(quotes []dto.UminekoQuote) []int {
	out := make([]int, len(quotes))
	for i, q := range quotes {
		out[i] = q.Index
	}
	return out
}

func TestFindQuotesParity(t *testing.T) {
	en := language.English
	e := newEngines(t, game.Umineko{}, testQuotes())

	cases := []struct {
		name   string
		filter params.QuoteFilter
		limit  int
		offset int
	}{
		{"all", params.QuoteFilter{Lang: en}, 50, 0},
		{"contains", params.QuoteFilter{Lang: en, Text: "without love"}, 50, 0},
		{"contains partial word", params.QuoteFilter{Lang: en, Text: "wit"}, 50, 0},
		{"no matches", params.QuoteFilter{Lang: en, Text: "xyzzy"}, 50, 0},
		{"exact word match", params.QuoteFilter{Lang: en, Text: "love", Exact: true}, 50, 0},
		{"exact partial word no match", params.QuoteFilter{Lang: en, Text: "wit", Exact: true}, 50, 0},
		{"exact word before punctuation", params.QuoteFilter{Lang: en, Text: "witch", Exact: true}, 50, 0},
		{"exact phrase", params.QuoteFilter{Lang: en, Text: "without love", Exact: true}, 50, 0},
		{"character", params.QuoteFilter{Lang: en, Text: "without", CharacterID: "10"}, 50, 0},
		{"unknown character", params.QuoteFilter{Lang: en, CharacterID: "99"}, 50, 0},
		{"unknown character with text", params.QuoteFilter{Lang: en, Text: "without", CharacterID: "99"}, 50, 0},
		{"unknown episode", params.QuoteFilter{Lang: en, Episode: 7}, 50, 0},
		{"unknown character excluding narrator", params.QuoteFilter{Lang: en, CharacterID: "99", ExcludeNarrator: true}, 50, 0},
		{"episode", params.QuoteFilter{Lang: en, Text: "you", Episode: 2}, 50, 0},
		{"interaction", params.QuoteFilter{Lang: en, InteractionA: "10", InteractionB: "27"}, 50, 0},
		{"interaction same character", params.QuoteFilter{Lang: en, InteractionA: "10", InteractionB: "10"}, 50, 0},
		{"truth", params.QuoteFilter{Lang: en, RedTruth: true}, 50, 0},
		{"segment", params.QuoteFilter{Lang: en, Segment: "2"}, 50, 0},
		{"exclude narrator", params.QuoteFilter{Lang: en, ExcludeNarrator: true}, 50, 0},
		{"unknown language", params.QuoteFilter{Lang: language.Russian}, 50, 0},
		{"first page", params.QuoteFilter{Lang: en}, 2, 0},
		{"second page", params.QuoteFilter{Lang: en}, 2, 2},
		{"partial last page", params.QuoteFilter{Lang: en}, 4, 4},
		{"offset beyond total", params.QuoteFilter{Lang: en, CharacterID: "10"}, 10, 100},
	}

	for _, c := range cases {
		mem, memTotal, err := e.memory.FindQuotes(c.filter, c.limit, c.offset)
		if err != nil {
			t.Fatalf("%s: memory FindQuotes: %v", c.name, err)
		}

		got, total, err := e.sqlite.FindQuotes(c.filter, c.limit, c.offset)
		if err != nil {
			t.Fatalf("%s: sqlite FindQuotes: %v", c.name, err)
		}

		if !reflect.DeepEqual(indexes(mem), indexes(got)) || memTotal != total {
			t.Errorf("%s: memory %v (total %d), sqlite %v (total %d)", c.name, indexes(mem), memTotal, indexes(got), total)
		}
	}
}

func TestFindQuotes_Paging(t *testing.T) {
	e := newEngines(t, game.Umineko{}, testQuotes())

	for name, d := range map[string]dao.QuoteDAO[dto.UminekoQuote]{"memory": e.memory, "sqlite": e.sqlite} {
		page, total, err := d.FindQuotes(params.QuoteFilter{Lang: language.English}, 2, 2)
		if err != nil {
			t.Fatalf("%s: FindQuotes: %v", name, err)
		}
		if got := indexes(page); !reflect.DeepEqual(got, []int{2, 3}) || total != 6 {
			t.Errorf("%s: second page: got %v total %d, want [2 3] total 6", name, got, total)
		}

		page, total, err = d.FindQuotes(params.QuoteFilter{Lang: language.English}, 10, 100)
		if err != nil {
			t.Fatalf("%s: FindQuotes: %v", name, err)
		}
		if page == nil || len(page) != 0 || total != 6 {
			t.Errorf("%s: offset beyond total: got %v total %d, want empty non-nil and total 6", name, page, total)
		}
	}
}

func TestFindQuotes_Interaction(t *testing.T) {
	found, _, err := newSeededDAO(t, game.Umineko{}, testQuotes()).FindQuotes(params.QuoteFilter{Lang: language.English, InteractionA: "27", InteractionB: "10"}, 50, 0)
	if err != nil {
		t.Fatalf("FindQuotes: %v", err)
	}

	if got := indexes(found); !reflect.DeepEqual(got, []int{1, 2, 4, 5}) {
		t.Errorf("interaction lines: got %v, want [1 2 4 5]", got)
	}
}

func TestRandomQuote(t *testing.T) {
	e := newEngines(t, game.Umineko{}, testQuotes())

	for name, d := range map[string]dao.QuoteDAO[dto.UminekoQuote]{"memory": e.memory, "sqlite": e.sqlite} {
		q, err := d.RandomQuote(params.QuoteFilter{Lang: language.English, CharacterID: "27"})
		if err != nil {
			t.Fatalf("%s: RandomQuote: %v", name, err)
		}
		if q == nil || q.CharacterID != "27" {
			t.Errorf("%s: random with character: got %+v", name, q)
		}

		q, err = d.RandomQuote(params.QuoteFilter{Lang: language.English, ExcludeNarrator: true})
		if err != nil {
			t.Fatalf("%s: RandomQuote: %v", name, err)
		}
		if q == nil || q.CharacterID == game.NarratorID {
			t.Errorf("%s: random excluding narrator: got %+v", name, q)
		}

		q, err = d.RandomQuote(params.QuoteFilter{Lang: language.English, BlueTruth: true})
		if err != nil {
			t.Fatalf("%s: RandomQuote: %v", name, err)
		}
		if q == nil || !q.HasBlueTruth {
			t.Errorf("%s: random with truth: got %+v", name, q)
		}

		q, err = d.RandomQuote(params.QuoteFilter{Lang: language.English, CharacterID: "nonexistent"})
		if err != nil {
			t.Fatalf("%s: RandomQuote: %v", name, err)
		}
		if q != nil {
			t.Errorf("%s: random with unknown character should return nil, got %+v", name, q)
		}
	}
}

func TestQuoteIndexByAudioID(t *testing.T) {
	s := newSeededDAO(t, game.Umineko{}, testQuotes())

	idx, ok, err := s.QuoteIndexByAudioID(language.English, "30101090")
	if err != nil {
		t.Fatalf("QuoteIndexByAudioID: %v", err)
	}
	if !ok || idx != 2 {
		t.Errorf("second audio id: got %d %v, want 2 true", idx, ok)
	}

	_, ok, err = s.QuoteIndexByAudioID(language.English, "nope")
	if err != nil {
		t.Fatalf("QuoteIndexByAudioID: %v", err)
	}
	if ok {
		t.Error("unknown audio id should not be found")
	}
}

func TestQuoteByIndex(t *testing.T) {
	s := newSeededDAO(t, game.Umineko{}, testQuotes())

	q, err := s.QuoteByIndex(language.English, 2)
	if err != nil {
		t.Fatalf("QuoteByIndex: %v", err)
	}
	if q == nil || q.Index != 2 || !q.HasBlueTruth {
		t.Errorf("QuoteByIndex round trip: got %+v", q)
	}

	q, err = s.QuoteByIndex(language.English, 100)
	if err != nil {
		t.Fatalf("QuoteByIndex: %v", err)
	}
	if q != nil {
		t.Error("QuoteByIndex out of range should return nil")
	}
}

func TestQuotesBetweenParity(t *testing.T) {
	e := newEngines(t, game.Umineko{}, testQuotes())

	mem, err := e.memory.QuotesBetween(language.English, -2, 2)
	if err != nil {
		t.Fatalf("memory QuotesBetween: %v", err)
	}

	got, err := e.sqlite.QuotesBetween(language.English, -2, 2)
	if err != nil {
		t.Fatalf("sqlite QuotesBetween: %v", err)
	}

	if !reflect.DeepEqual(indexes(mem), indexes(got)) {
		t.Errorf("between: memory %v, sqlite %v", indexes(mem), indexes(got))
	}
	if !reflect.DeepEqual(indexes(got), []int{0, 1, 2}) {
		t.Errorf("between: got %v, want [0 1 2]", indexes(got))
	}
}

func TestNearestVoiced(t *testing.T) {
	s := newSeededDAO(t, game.Umineko{}, testQuotes())

	next, err := s.NearestVoiced(language.English, 1, false)
	if err != nil {
		t.Fatalf("NearestVoiced: %v", err)
	}
	if next != "30101089, 30101090" {
		t.Errorf("next: got %q", next)
	}

	prev, err := s.NearestVoiced(language.English, 4, true)
	if err != nil {
		t.Fatalf("NearestVoiced: %v", err)
	}
	if prev != "30101089, 30101090" {
		t.Errorf("prev: got %q", prev)
	}

	none, err := s.NearestVoiced(language.English, 5, false)
	if err != nil {
		t.Fatalf("NearestVoiced: %v", err)
	}
	if none != "" {
		t.Errorf("past the end: got %q, want empty", none)
	}
}

func TestLanguageCounts(t *testing.T) {
	counts, err := newSeededDAO(t, game.Umineko{}, testQuotes()).LanguageCounts()
	if err != nil {
		t.Fatalf("LanguageCounts: %v", err)
	}

	if counts[language.English] != 6 || counts[language.Japanese] != 1 {
		t.Errorf("counts: got %v", counts)
	}
}

func TestSpeakerTalliesParity(t *testing.T) {
	umineko := newEngines(t, game.Umineko{}, testQuotes())
	ciconia := newEngines(t, game.Ciconia{}, ciconiaTestQuotes())

	want, err := umineko.memory.SpeakerTallies(language.English)
	if err != nil {
		t.Fatalf("memory SpeakerTallies: %v", err)
	}

	got, err := umineko.sqlite.SpeakerTallies(language.English)
	if err != nil {
		t.Fatalf("sqlite SpeakerTallies: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("umineko: sqlite tallies differ from memory\nwant: %+v\ngot:  %+v", want, got)
	}

	want, err = ciconia.memory.SpeakerTallies(language.English)
	if err != nil {
		t.Fatalf("memory SpeakerTallies: %v", err)
	}

	got, err = ciconia.sqlite.SpeakerTallies(language.English)
	if err != nil {
		t.Fatalf("sqlite SpeakerTallies: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("ciconia: sqlite tallies differ from memory\nwant: %+v\ngot:  %+v", want, got)
	}
	if got.SegmentInteractions["02"]["jayden|mitake"] != 1 {
		t.Errorf("chapter 02 interactions: got %+v", got.SegmentInteractions["02"])
	}
}

func TestSpeakersParity(t *testing.T) {
	e := newEngines(t, game.Ciconia{}, ciconiaTestQuotes())

	want, err := e.memory.Speakers(language.English)
	if err != nil {
		t.Fatalf("memory Speakers: %v", err)
	}

	got, err := e.sqlite.Speakers(language.English)
	if err != nil {
		t.Fatalf("sqlite Speakers: %v", err)
	}

	if !reflect.DeepEqual(want, got) {
		t.Errorf("sqlite speakers differ from memory\nwant: %+v\ngot:  %+v", want, got)
	}
	if got["mitake"] != "Mitake Miyao" {
		t.Errorf("expected first-seen character name, got %q", got["mitake"])
	}
}

func TestCiconiaRoundTrip(t *testing.T) {
	d := newSeededDAO(t, game.Ciconia{}, ciconiaTestQuotes())

	en, err := d.QuoteByIndex(language.English, 4)
	if err != nil {
		t.Fatalf("QuoteByIndex en: %v", err)
	}
	if en == nil || en.Chapter != "02" || en.Text != "text" {
		t.Errorf("english round trip: got %+v", en)
	}

	ja, err := d.QuoteByIndex(language.Japanese, 4)
	if err != nil {
		t.Fatalf("QuoteByIndex ja: %v", err)
	}
	if ja == nil || ja.Chapter != "02" || ja.Text != "テキスト" {
		t.Errorf("japanese round trip: got %+v", ja)
	}
}

func TestTruthPerEpisode(t *testing.T) {
	d := &uminekoDAO{quoteDAO: newSeededDAO(t, game.Umineko{}, testQuotes())}

	got, err := d.TruthPerEpisode(language.English)
	if err != nil {
		t.Fatalf("TruthPerEpisode: %v", err)
	}

	want := map[int]params.TruthCounts{
		1: {Red: 1, Blue: 1},
		2: {Gold: 1},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("truth per episode: want %v, got %v", want, got)
	}
}

func TestSkipsSeedWhenPopulated(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")

	conn := openTestDB(t, path)
	s := &quoteDAO[dto.UminekoQuote]{db: conn, game: game.Umineko{}}
	if err := s.Seed(testQuotes(), testDataHash); err != nil {
		t.Fatalf("Seed: %v", err)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("close database: %v", err)
	}

	conn = openTestDB(t, path)
	s = &quoteDAO[dto.UminekoQuote]{db: conn, game: game.Umineko{}}

	needsSeed, err := s.NeedsSeed(testDataHash)
	if err != nil {
		t.Fatalf("NeedsSeed: %v", err)
	}
	if needsSeed {
		t.Fatal("expected matching hash to skip seeding after reopen")
	}

	q, err := s.QuoteByIndex(language.English, 1)
	if err != nil {
		t.Fatalf("QuoteByIndex: %v", err)
	}
	if q == nil || q.Text != "Without love, it cannot be seen." {
		t.Errorf("quote after reopen: got %+v", q)
	}
}

func TestReseedsOnHashChange(t *testing.T) {
	s := newSeededDAO(t, game.Umineko{}, testQuotes())

	needsSeed, err := s.NeedsSeed("new-hash")
	if err != nil {
		t.Fatalf("NeedsSeed: %v", err)
	}
	if !needsSeed {
		t.Fatal("expected changed hash to trigger re-seed")
	}

	updated := testQuotes()
	updated[language.English][1].Text = "Without love, it cannot be seen. Really."

	if err := s.Seed(updated, "new-hash"); err != nil {
		t.Fatalf("re-Seed: %v", err)
	}

	q, err := s.QuoteByIndex(language.English, 1)
	if err != nil {
		t.Fatalf("QuoteByIndex: %v", err)
	}
	if q == nil || q.Text != "Without love, it cannot be seen. Really." {
		t.Errorf("quote after re-seed: got %+v", q)
	}

	needsSeed, err = s.NeedsSeed("new-hash")
	if err != nil {
		t.Fatalf("NeedsSeed after re-seed: %v", err)
	}
	if needsSeed {
		t.Fatal("expected new hash to be recorded after re-seed")
	}
}
