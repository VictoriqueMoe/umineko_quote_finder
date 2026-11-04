package sqlite

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"runtime/debug"
	"strings"
	"time"
	"unicode"

	"umineko_quote/internal/quote"
	"umineko_quote/internal/quote/game"
	"umineko_quote/internal/quote/language"
	"umineko_quote/internal/quote/params"

	scriptdto "github.com/VictoriqueMoe/umineko_script_parser/dto"
)

type quoteDAO[Q any] struct {
	db   *sql.DB
	game game.Game[Q]
}

const (
	seedVersion   = "4"
	parserModule  = "github.com/VictoriqueMoe/umineko_script_parser"
	selectColumns = `idx, text, text_html, character_id, character_name, audio_id, audio_char_map, audio_text_map, episode, segment, content_type, sound_effects, has_red_truth, has_blue_truth, has_gold_truth, has_purple_truth, text_jp, text_jp_html`
	insertColumns = `game, lang, ` + selectColumns + `, text_lower, text_words`
)

var (
	likeEscaper = strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	insertQuote = `INSERT INTO quotes (` + insertColumns + `) VALUES (?` + strings.Repeat(`, ?`, strings.Count(insertColumns, ",")) + `)`
)

func newQuotes[Q any](db *sql.DB, g game.Game[Q]) (*quoteDAO[Q], error) {
	d := &quoteDAO[Q]{
		db:   db,
		game: g,
	}

	files, err := g.DataFiles()
	if err != nil {
		return nil, fmt.Errorf("list %s data files: %w", g.Name(), err)
	}

	hash, err := hashFiles(quote.DataFS, files)
	if err != nil {
		return nil, fmt.Errorf("hash %s data files: %w", g.Name(), err)
	}

	needsSeed, err := d.NeedsSeed(hash)
	if err != nil {
		return nil, fmt.Errorf("check %s seed: %w", g.Name(), err)
	}

	if !needsSeed {
		log.Printf("[%s] sqlite data up to date", g.Name())
		return d, nil
	}

	seedStart := time.Now()

	quotes, err := g.Load()
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", g.Name(), err)
	}

	if err := d.Seed(quotes, hash); err != nil {
		return nil, fmt.Errorf("seed %s: %w", g.Name(), err)
	}

	log.Printf("[%s] seeded sqlite database in %v", g.Name(), time.Since(seedStart).Round(time.Millisecond))
	return d, nil
}

func hashFiles(fsys fs.FS, paths []string) (string, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "", errors.New("build info unavailable, cannot read the script parser version")
	}

	h := sha256.New()
	h.Write([]byte(seedVersion))

	for _, dep := range info.Deps {
		if dep.Path != parserModule {
			continue
		}
		if dep.Replace != nil {
			dep = dep.Replace
		}
		h.Write([]byte(dep.Path + "@" + dep.Version + dep.Sum))
	}

	for _, path := range paths {
		f, err := fsys.Open(path)
		if err != nil {
			return "", err
		}

		h.Write([]byte(path))
		_, copyErr := io.Copy(h, f)
		closeErr := f.Close()

		if copyErr != nil {
			return "", copyErr
		}
		if closeErr != nil {
			return "", closeErr
		}
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

func (s *quoteDAO[Q]) NeedsSeed(dataHash string) (bool, error) {
	var hash string
	err := s.db.QueryRow(`SELECT data_hash FROM seed_meta WHERE game = ?`, s.game.Name()).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return hash != dataHash, nil
}

func (s *quoteDAO[Q]) Seed(quotes map[language.Language][]Q, dataHash string) error {
	name := s.game.Name()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, table := range []string{`quotes`, `quote_audio_ids`} {
		if _, err := tx.Exec(`DELETE FROM `+table+` WHERE game = ?`, name); err != nil {
			return err
		}
	}

	quoteStmt, err := tx.Prepare(insertQuote)
	if err != nil {
		return err
	}
	defer quoteStmt.Close()

	audioStmt, err := tx.Prepare(`INSERT OR IGNORE INTO quote_audio_ids (game, lang, audio_id, quote_idx) VALUES (?,?,?,?)`)
	if err != nil {
		return err
	}
	defer audioStmt.Close()

	for lang, qs := range quotes {
		for _, q := range qs {
			bq := s.game.Base(&q)
			row := s.game.Row(&q)

			charMap, err := jsonOrNil(bq.AudioCharMap)
			if err != nil {
				return err
			}

			textMap, err := jsonOrNil(bq.AudioTextMap)
			if err != nil {
				return err
			}

			sfx, err := jsonOrNil(bq.SoundEffects)
			if err != nil {
				return err
			}

			_, err = quoteStmt.Exec(
				name, string(lang),
				row.Index, bq.Text, bq.TextHtml, bq.CharacterID, bq.Character, bq.AudioID, charMap, textMap,
				bq.Episode, row.Segment, bq.ContentType, sfx,
				row.HasRedTruth, row.HasBlueTruth, row.HasGoldTruth, row.HasPurpleTruth,
				row.TextJP, row.TextJPHtml,
				strings.ToLower(bq.Text), " "+wordText(bq.Text)+" ",
			)
			if err != nil {
				return err
			}

			if bq.AudioID != "" {
				for _, id := range strings.Split(bq.AudioID, ", ") {
					if _, err := audioStmt.Exec(name, string(lang), id, row.Index); err != nil {
						return err
					}
				}
			}
		}
	}

	_, err = tx.Exec(
		`INSERT INTO seed_meta (game, data_hash) VALUES (?, ?) ON CONFLICT(game) DO UPDATE SET data_hash = excluded.data_hash`,
		name, dataHash,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func jsonOrNil[T map[string]string | []scriptdto.SoundEffect](v T) (any, error) {
	if len(v) == 0 {
		return nil, nil
	}

	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func wordText(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return unicode.ToLower(r)
		}
		return ' '
	}, s)
}

func (s *quoteDAO[Q]) FindQuotes(f params.QuoteFilter, limit int, offset int) ([]Q, int, error) {
	clauses, args := s.filterClauses(f)
	where, args := s.where(f.Lang, clauses, args)

	rows, err := s.db.Query(`SELECT `+selectColumns+`, COUNT(*) OVER () FROM quotes`+where+` ORDER BY idx LIMIT ? OFFSET ?`, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	page := []Q{}
	total := 0
	for rows.Next() {
		q, err := s.scanQuote(rows, &total)
		if err != nil {
			return nil, 0, err
		}
		page = append(page, q)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if len(page) == 0 && offset > 0 {
		if err := s.db.QueryRow(`SELECT COUNT(*) FROM quotes`+where, args...).Scan(&total); err != nil {
			return nil, 0, err
		}
	}

	return page, total, nil
}

func (s *quoteDAO[Q]) RandomQuote(f params.QuoteFilter) (*Q, error) {
	clauses, args := s.filterClauses(f)
	where, args := s.where(f.Lang, clauses, args)

	return s.getQuote(`SELECT `+selectColumns+` FROM quotes`+where+` ORDER BY RANDOM() LIMIT 1`, args...)
}

func (s *quoteDAO[Q]) QuoteIndexByAudioID(lang language.Language, audioID string) (int, bool, error) {
	var idx int
	err := s.db.QueryRow(
		`SELECT quote_idx FROM quote_audio_ids WHERE game = ? AND lang = ? AND audio_id = ? ORDER BY quote_idx DESC LIMIT 1`,
		s.game.Name(), string(lang), audioID,
	).Scan(&idx)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return idx, true, nil
}

func (s *quoteDAO[Q]) QuoteByIndex(lang language.Language, index int) (*Q, error) {
	return s.getQuote(`SELECT `+selectColumns+` FROM quotes WHERE game = ? AND lang = ? AND idx = ?`, s.game.Name(), string(lang), index)
}

func (s *quoteDAO[Q]) QuotesBetween(lang language.Language, from int, to int) ([]Q, error) {
	rows, err := s.db.Query(
		`SELECT `+selectColumns+` FROM quotes WHERE game = ? AND lang = ? AND idx BETWEEN ? AND ? ORDER BY idx`,
		s.game.Name(), string(lang), from, to,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	between := []Q{}
	for rows.Next() {
		q, err := s.scanQuote(rows)
		if err != nil {
			return nil, err
		}
		between = append(between, q)
	}

	return between, rows.Err()
}

func (s *quoteDAO[Q]) NearestVoiced(lang language.Language, index int, previous bool) (string, error) {
	query := `SELECT audio_id FROM quotes WHERE game = ? AND lang = ? AND idx > ? AND audio_id != '' ORDER BY idx ASC LIMIT 1`
	if previous {
		query = `SELECT audio_id FROM quotes WHERE game = ? AND lang = ? AND idx < ? AND audio_id != '' ORDER BY idx DESC LIMIT 1`
	}

	var audioID string
	err := s.db.QueryRow(query, s.game.Name(), string(lang), index).Scan(&audioID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return audioID, nil
}

func (s *quoteDAO[Q]) LanguageCounts() (map[language.Language]int, error) {
	rows, err := s.db.Query(`SELECT lang, COUNT(*) FROM quotes WHERE game = ? GROUP BY lang`, s.game.Name())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[language.Language]int)
	for rows.Next() {
		var lang string
		var count int
		if err := rows.Scan(&lang, &count); err != nil {
			return nil, err
		}
		counts[language.Language(lang)] = count
	}

	return counts, rows.Err()
}

func (s *quoteDAO[Q]) Speakers(lang language.Language) (map[string]string, error) {
	rows, err := s.db.Query(`SELECT character_id, character_name, MIN(idx) FROM quotes WHERE game = ? AND lang = ? GROUP BY character_id`, s.game.Name(), string(lang))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	speakers := map[string]string{}
	for rows.Next() {
		var id string
		var name string
		if err := rows.Scan(&id, &name, new(int)); err != nil {
			return nil, err
		}
		speakers[id] = name
	}

	return speakers, rows.Err()
}

func (s *quoteDAO[Q]) SpeakerTallies(lang language.Language) (params.SpeakerTallies, error) {
	tallies := params.NewSpeakerTallies()

	if err := s.tallyLines(lang, tallies); err != nil {
		return tallies, err
	}

	if err := s.tallyInteractions(lang, tallies); err != nil {
		return tallies, err
	}

	return tallies, nil
}

func (s *quoteDAO[Q]) tallyLines(lang language.Language, tallies params.SpeakerTallies) error {
	rows, err := s.db.Query(
		`SELECT character_id, segment, COUNT(*) FROM quotes WHERE game = ? AND lang = ? AND character_id != ? GROUP BY character_id, segment`,
		s.game.Name(), string(lang), game.NarratorID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id string
		var segment string
		var count int
		if err := rows.Scan(&id, &segment, &count); err != nil {
			return err
		}
		tallies.AddLines(id, segment, count)
	}

	return rows.Err()
}

func (s *quoteDAO[Q]) tallyInteractions(lang language.Language, tallies params.SpeakerTallies) error {
	rows, err := s.db.Query(`SELECT a.segment, MIN(a.character_id, b.character_id) AS first, MAX(a.character_id, b.character_id) AS second, COUNT(*)
		FROM quotes a
		JOIN quotes b ON b.game = a.game AND b.lang = a.lang AND b.idx = a.idx + 1 AND b.segment = a.segment
		WHERE a.game = ? AND a.lang = ?
			AND a.character_id NOT IN (?, '') AND b.character_id NOT IN (?, '')
			AND a.character_id != b.character_id
		GROUP BY a.segment, first, second`,
		s.game.Name(), string(lang), game.NarratorID, game.NarratorID,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var segment string
		var first string
		var second string
		var count int
		if err := rows.Scan(&segment, &first, &second, &count); err != nil {
			return err
		}
		tallies.AddInteractions(segment, first, second, count)
	}

	return rows.Err()
}

func (s *quoteDAO[Q]) filterClauses(f params.QuoteFilter) ([]string, []any) {
	var clauses []string
	var args []any

	if f.Text != "" && f.Exact {
		clauses = append(clauses, `text_words LIKE ? ESCAPE '\'`)
		args = append(args, "% "+likeEscaper.Replace(wordText(f.Text))+" %")
	}

	if f.Text != "" && !f.Exact {
		clauses = append(clauses, `text_lower LIKE ? ESCAPE '\'`)
		args = append(args, "%"+likeEscaper.Replace(strings.ToLower(f.Text))+"%")
	}

	if f.CharacterID != "" {
		clauses = append(clauses, `character_id = ?`)
		args = append(args, f.CharacterID)
	}

	if f.Episode > 0 {
		clauses = append(clauses, `episode = ?`)
		args = append(args, f.Episode)
	}

	if f.Segment != "" {
		clauses = append(clauses, `segment = ?`)
		args = append(args, f.Segment)
	}

	if f.RedTruth {
		clauses = append(clauses, `has_red_truth = 1`)
	}

	if f.BlueTruth {
		clauses = append(clauses, `has_blue_truth = 1`)
	}

	if f.GoldTruth {
		clauses = append(clauses, `has_gold_truth = 1`)
	}

	if f.PurpleTruth {
		clauses = append(clauses, `has_purple_truth = 1`)
	}

	if f.InteractionA != "" && f.InteractionB != "" {
		clauses = append(clauses, `idx IN (SELECT a.idx
			FROM quotes a
			JOIN quotes b ON b.game = a.game AND b.lang = a.lang AND b.episode = a.episode AND b.idx IN (a.idx - 1, a.idx + 1)
			WHERE a.game = ? AND a.lang = ?
				AND a.character_id IN (?, ?) AND b.character_id IN (?, ?)
				AND a.character_id != b.character_id
				AND a.character_id NOT IN (?, '') AND b.character_id NOT IN (?, ''))`)
		args = append(args, s.game.Name(), string(f.Lang), f.InteractionA, f.InteractionB, f.InteractionA, f.InteractionB, game.NarratorID, game.NarratorID)
	}

	if f.ExcludeNarrator {
		clauses = append(clauses, `character_id != ?`)
		args = append(args, game.NarratorID)
	}

	return clauses, args
}

func (s *quoteDAO[Q]) where(lang language.Language, clauses []string, args []any) (string, []any) {
	where := ` WHERE game = ? AND lang = ?`
	if len(clauses) > 0 {
		where += ` AND ` + strings.Join(clauses, ` AND `)
	}

	return where, append([]any{s.game.Name(), string(lang)}, args...)
}

func (s *quoteDAO[Q]) getQuote(query string, args ...any) (*Q, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, rows.Err()
	}

	q, err := s.scanQuote(rows)
	if err != nil {
		return nil, err
	}

	return &q, nil
}

func (s *quoteDAO[Q]) scanQuote(rows *sql.Rows, extra ...any) (Q, error) {
	var (
		bq      scriptdto.BaseQuote
		row     game.QuoteRow
		charMap sql.NullString
		textMap sql.NullString
		sfx     sql.NullString
		empty   Q
	)

	dest := []any{
		&row.Index, &bq.Text, &bq.TextHtml, &bq.CharacterID, &bq.Character, &bq.AudioID,
		&charMap, &textMap, &bq.Episode, &row.Segment, &bq.ContentType, &sfx,
		&row.HasRedTruth, &row.HasBlueTruth, &row.HasGoldTruth, &row.HasPurpleTruth,
		&row.TextJP, &row.TextJPHtml,
	}

	if err := rows.Scan(append(dest, extra...)...); err != nil {
		return empty, err
	}

	if charMap.Valid {
		if err := json.Unmarshal([]byte(charMap.String), &bq.AudioCharMap); err != nil {
			return empty, err
		}
	}

	if textMap.Valid {
		if err := json.Unmarshal([]byte(textMap.String), &bq.AudioTextMap); err != nil {
			return empty, err
		}
	}

	if sfx.Valid {
		if err := json.Unmarshal([]byte(sfx.String), &bq.SoundEffects); err != nil {
			return empty, err
		}
	}

	return s.game.Build(bq, row), nil
}
