-- +goose Up
CREATE TABLE quotes (
    game             TEXT    NOT NULL,
    lang             TEXT    NOT NULL,
    idx              INTEGER NOT NULL,
    text             TEXT    NOT NULL,
    text_html        TEXT    NOT NULL DEFAULT '',
    text_lower       TEXT    NOT NULL DEFAULT '',
    text_words       TEXT    NOT NULL DEFAULT '',
    character_id     TEXT    NOT NULL DEFAULT '',
    character_name   TEXT    NOT NULL DEFAULT '',
    audio_id         TEXT    NOT NULL DEFAULT '',
    audio_char_map   TEXT,
    audio_text_map   TEXT,
    episode          INTEGER NOT NULL DEFAULT 0,
    segment          TEXT    NOT NULL DEFAULT '',
    content_type     TEXT    NOT NULL DEFAULT '',
    sound_effects    TEXT,
    has_red_truth    INTEGER NOT NULL DEFAULT 0,
    has_blue_truth   INTEGER NOT NULL DEFAULT 0,
    has_gold_truth   INTEGER NOT NULL DEFAULT 0,
    has_purple_truth INTEGER NOT NULL DEFAULT 0,
    text_jp          TEXT    NOT NULL DEFAULT '',
    text_jp_html     TEXT    NOT NULL DEFAULT '',
    PRIMARY KEY (game, lang, idx)
);

CREATE INDEX idx_quotes_character ON quotes (game, lang, character_id);
CREATE INDEX idx_quotes_episode ON quotes (game, lang, episode);
CREATE INDEX idx_quotes_segment ON quotes (game, lang, segment);

CREATE TABLE quote_audio_ids (
    game      TEXT    NOT NULL,
    lang      TEXT    NOT NULL,
    audio_id  TEXT    NOT NULL,
    quote_idx INTEGER NOT NULL,
    PRIMARY KEY (game, lang, audio_id, quote_idx)
);

CREATE TABLE seed_meta (
    game      TEXT PRIMARY KEY,
    data_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE seed_meta;
DROP TABLE quote_audio_ids;
DROP TABLE quotes;
