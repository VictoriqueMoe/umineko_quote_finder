# Umineko Quote Search

A quote search engine for Umineko no Naku Koro ni. Search through thousands of lines of dialogue from the visual novel.

## Contents

- [Features](#features)
- [Quick Start](#quick-start)
  - [Voice Audio (Optional)](#voice-audio-optional)
  - [Expected zip structure](#expected-zip-structure)
- [API Endpoints](#api-endpoints)
  - [Query Parameters](#query-parameters)
  - [Response Format](#response-format)
- [Build](#build)
  - [Cross-compile](#cross-compile)
- [Docker](#docker)
- [Data](#data)
- [Architecture: The Lexar Package](#architecture-the-lexar-package)
  - [Pipeline Overview](#pipeline-overview)
  - [Package Structure](#package-structure)
  - [Key Design Decisions](#key-design-decisions)
- [Script Tag Parsing](#script-tag-parsing)
  - [Tags with HTML rendering](#tags-with-html-rendering)
  - [Preset colour reference](#preset-colour-reference)
  - [Tags stripped to content](#tags-stripped-to-content)
  - [Special character tags](#special-character-tags)
  - [Other cleanup](#other-cleanup)
- [Contributors](#contributors)

## Features

- Search through all dialogue
- Filter by character and episode
- Random quote generator
- Scene context viewer, see surrounding dialogue for any quote
- English/Japanese language toggle
- Inline audio playback for voiced lines
- Umineko-themed web interface

## Quick Start

```bash
# Build the frontend
cd frontend
npm install
npm run build
cd ..

# Build and run the Go server
go build -o umineko_quote .
./umineko_quote
```

Open http://127.0.0.1:3000

### Development

For frontend development with hot reload, run the Vite dev server alongside the Go backend:

```bash
# Terminal 1: Go backend
go run main.go

# Terminal 2: Vite dev server (proxies /api to :3000)
cd frontend
npm run dev
```

Open http://localhost:5173

### Voice Audio (Optional)

Audio playback requires a zip of the voice files. Create a `.env` file in the project root:

```env
# URL to download the zip
VOICE_ZIP_URL=https://example.com/voice.zip

# Or a local path to the zip
VOICE_ZIP_URL=C:\path\to\voice.zip
```

Then run the setup script:

**Linux / macOS:**
```bash
./setup_audio.sh
```

**Windows (PowerShell):**
```powershell
.\setup_audio.ps1
```

The script will detect whether `VOICE_ZIP_URL` is a local file or a URL and handle it accordingly. If the audio directory already exists, it skips extraction.

The app works without audio files, quotes will display normally but without playback controls.

### Expected zip structure

The zip must contain a `voice/` directory at its root with character ID subdirectories:

```
voice.zip
└── voice/
    ├── 00/
    │   ├── 00100001.ogg
    │   └── ...
    ├── 01/
    └── ...
```

## API Endpoints

| Endpoint                             | Description                            |
|--------------------------------------|----------------------------------------|
| `GET /api/v1/search`                 | Search quotes                          |
| `GET /api/v1/random`                 | Get random quote                       |
| `GET /api/v1/character/:id`          | Get quotes by character ID             |
| `GET /api/v1/context/:audioId`       | Get surrounding dialogue for a quote   |
| `GET /api/v1/characters`             | List all character IDs and names       |
| `GET /api/v1/audio/:charId/:audioId` | Stream audio file for a voice line     |
| `GET /api/v1/health`                 | Health check                           |

### Query Parameters

| Parameter   | Endpoints                          | Description                                        |
|-------------|------------------------------------|----------------------------------------------------|
| `q`         | search                             | Search query (required)                            |
| `lang`      | search, random, character, context | Language: `en` (default) or `ja`                   |
| `character` | search, random                     | Filter by character ID                             |
| `episode`   | search, random, character          | Filter by episode (1-8)                            |
| `lines`     | context                            | Number of lines before/after (default: 5, max: 20) |
| `limit`     | search, character                  | Results per page (default: 30)                     |
| `offset`    | search, character                  | Pagination offset                                  |

### Response Format

```json
{
  "results": [
    {
      "quote": {
        "text": "Without love, it cannot be seen.",
        "textHtml": "Without love, it cannot be seen.",
        "characterId": "27",
        "character": "Beatrice",
        "audioId": "10700001",
        "episode": 1,
        "contentType": ""
      },
      "score": 95
    }
  ]
}
```

The `contentType` field distinguishes content sections: `""` for main episodes, `"tea"` for tea parties, `"ura"` for ???? chapters, and `"omake"` for omakes (bonus content).

## Build

The frontend must be built before the Go binary, as the Go binary embeds the `static/` directory.

```bash
cd frontend && npm ci && npm run build && cd ..
```

### Windows
```powershell
go build -o umineko_quote.exe .
```

### Linux
```bash
go build -o umineko_quote .
```

### Cross-compile

```powershell
# Mac ARM (M1/M2/M3)
$env:GOOS="darwin"; $env:GOARCH="arm64"; go build -o umineko_quote_mac .; $env:GOOS=""; $env:GOARCH=""

# Mac Intel
$env:GOOS="darwin"; $env:GOARCH="amd64"; go build -o umineko_quote_mac_intel .; $env:GOOS=""; $env:GOARCH=""

# Linux x64
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -o umineko_quote_linux .; $env:GOOS=""; $env:GOARCH=""
```

## Docker

Requires a `.env` file with `VOICE_ZIP_URL` set (URL only for Docker builds).

```bash
docker compose up -d --build
```

## Data

Quote data is parsed from Umineko no Naku Koro ni script files:

```
internal/quote/data/
├── english.txt
├── japanese.txt
└── audio/          (extracted via setup script or Docker build)
    ├── 00/
    ├── 01/
    ├── ...
    └── 99/
```

Text files are embedded at compile time. Audio files are read from disk at runtime and are organized by character ID subdirectory.

## Architecture: The Lexar Package

The `internal/lexar` package handles parsing Umineko script files and extracting quotes. It follows a pipeline architecture that separates concerns.

### Pipeline Overview

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                              Source Text                                    │
│  d [lv 0*"27"*"10100001"]`"{p:1:Without love, it cannot be seen.}"`[\]     │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                           LEXER (lexer.go)                                  │
│  Tokenises input into a stream of typed tokens                              │
│  • TokenCommand: "d"                                                        │
│  • TokenInlineCommand: "lv 0*\"27\"*\"10100001\""                           │
│  • TokenBacktick: "`"                                                       │
│  • TokenFormatTag: "p:1:Without love, it cannot be seen."                   │
│  • etc.                                                                     │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                          PARSER (parser.go)                                 │
│  Builds Abstract Syntax Tree from tokens                                    │
│                                                                             │
│  Script                                                                     │
│   └── Lines[]                                                               │
│        ├── EpisodeMarkerLine { Episode: 1, Type: "episode" }                │
│        ├── PresetDefineLine { ID: 1, Colour: "#FF0000" }                    │
│        ├── LabelLine { Name: "ep1_scene1" }                                 │
│        └── DialogueLine                                                     │
│             ├── Command: "d"                                                │
│             └── Content[]                                                   │
│                  ├── VoiceCommand { CharacterID: "27", AudioID: "..." }     │
│                  └── FormatTag { Name: "p", Param: "1", Content: [...] }    │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                       EXTRACTOR (extractor.go)                              │
│  Walks AST, extracts quotes with metadata                                   │
│                                                                             │
│  ExtractedQuote {                                                           │
│      Content:     []DialogueElement  ◄── Raw AST, not yet transformed       │
│      CharacterID: "27"                                                      │
│      AudioID:     "10100001"                                                │
│      Episode:     1                                                         │
│      Truth:       { HasRed: true, HasBlue: false }                          │
│  }                                                                          │
└─────────────────────────────────────────────────────────────────────────────┘
                                      │
                                      ▼
┌─────────────────────────────────────────────────────────────────────────────┐
│                   TRANSFORMER FACTORY (transformer/)                        │
│  Converts raw AST to output format on demand                                │
│                                                                             │
│  factory.MustGet(FormatPlainText) ──► "Without love, it cannot be seen."    │
│  factory.MustGet(FormatHTML)      ──► "<span class=\"red-truth\">...</span>"│
│  factory.MustGet(FormatJSON)      ──► (add your own transformer)            │
└─────────────────────────────────────────────────────────────────────────────┘
```

### Package Structure

```
internal/lexar/
├── ast/                    # Abstract Syntax Tree types
│   └── ast.go              # Token, Line, DialogueElement types
├── transformer/            # Output format transformers
│   ├── transformer.go      # Transformer interface
│   ├── factory.go          # Factory for obtaining transformers
│   ├── preset.go           # Preset colour/class context
│   ├── plaintext.go        # Plain text output
│   └── html.go             # HTML output with styling
├── lexer.go                # Tokeniser
├── parser.go               # AST builder
├── extractor.go            # Quote extraction
└── truth.go                # Red/blue truth detection
```

### Key Design Decisions

**AST stores raw content**, the extractor outputs `ExtractedQuote` with raw `[]DialogueElement`, not pre-transformed strings. This allows transformation to happen on-demand via the factory.

**Factory pattern for transformers**, adding a new output format (e.g., JSON, Markdown) requires:
1. Implement the `Transformer` interface
2. Register it in the factory

No changes are needed to the extractor or parser.

**Preset context**, colour presets (`{p:1:text}`) are defined in script headers via `preset_define`. The `PresetContext` collects these definitions and provides semantic class lookups (preset 1 → "red-truth", preset 2 → "blue-truth") and dynamic colour lookups for other presets.

**Truth detection**, red and blue truth are detected by walking the AST looking for preset tags with semantic classes. This is stored as `TruthFlags` with `HasRed` and `HasBlue` booleans, allowing quotes with mixed truth (both red and blue) to appear in both filters.

## Script Tag Parsing

The source text files use [ONScripter-RU](https://github.com/umineko-project/onscripter-ru) dialogue formatting. The parser strips or converts these tags for display. Tags are processed in a loop to handle nesting (e.g. `{nobr:{m:-5:——}—}`).

### Tags with HTML rendering

| Script tag                          | Plain text     | HTML                                                    |
|-------------------------------------|----------------|---------------------------------------------------------|
| `{n}`                               | space          | `<br>`                                                  |
| `{i:text}` / `{italic:text}`        | text           | `<em>text</em>`                                         |
| `{c:HEX:text}` / `{color:HEX:text}` | text           | `<span style="color:#HEX">text</span>`                  |
| `{p:1:text}` (red truth)            | text           | `<span class="red-truth">text</span>`                   |
| `{p:2:text}` (blue truth)           | text           | `<span class="blue-truth">text</span>`                  |
| `{p:41:text}` (gold text)           | text           | `<span style="color:#FFAA00">text</span>`               |
| `{p:42:text}` (purple text)         | text           | `<span style="color:#AA71FF">text</span>`               |
| `{ruby:reading:text}`               | text (reading) | `<ruby>text<rp>(</rp><rt>reading</rt><rp>)</rp></ruby>` |

### Preset colour reference

The `{p:N:text}` tag applies a style preset defined in the script header via `preset_define`. The format is `preset_define number,font,size,colour,...`. Only presets that appear in actual dialogue lines are rendered with colour; the rest are stripped to plain text.

**Game presets** (used in dialogue):

| Preset | Colour    | Usage         | Rendering                           |
|--------|-----------|---------------|-------------------------------------|
| 0      | `#FFFFFF` | Japanese font | Stripped (white on dark is default) |
| 1      | `#FF0000` | Red truth     | `<span class="red-truth">`          |
| 2      | `#39C6FF` | Blue truth    | `<span class="blue-truth">`         |
| 7      | `#C0FFFF` | Chapter/Hint  | Not used in dialogue                |
| 41     | `#FFAA00` | Gold text     | `<span style="color:#FFAA00">`      |
| 42     | `#AA71FF` | Purple text   | `<span style="color:#AA71FF">`      |

**Menu/UI presets** (not rendered, stripped to plain text if they appear):

| Preset | Usage                          |
|--------|--------------------------------|
| 3      | Menu character text            |
| 4      | Menu JP text                   |
| 5      | Menu tips/notes text           |
| 6      | Music box BGM titles           |
| 8–9    | Menu titles and buttons        |
| 10     | Menu first setting line        |
| 11–12  | Menu buttons                   |
| 13     | Menu tips/notes titles         |
| 14–16  | Menu jump titles/lines         |
| 18     | Trophy description             |
| 20–25  | Credits                        |
| 30–31  | Load/Save                      |
| 32     | EP8 menu murder                |

### Tags stripped to content

These tags control visual styling in the game engine (font size, spacing, line breaking, gradients, etc.) that doesn't apply in a web context. The tag is removed and the inner text is kept.

| Script tag                                 | Result      |
|--------------------------------------------|-------------|
| `{f:N:text}` / `{font:N:text}`             | text        |
| `{p:N:text}` (other presets)               | text        |
| `{bold:text}` / `{b:text}`                 | text        |
| `{bolditalic:text}` / `{x:text}`           | text        |
| `{underline:text}` / `{u:text}`            | text        |
| `{gradient:N:text}` / `{g:N:text}`         | text        |
| `{nobreak:text}` / `{nobr:text}`           | text        |
| `{fit:text}` / `{j:text}`                  | text        |
| `{center:text}` / `{ac:text}`              | text        |
| `{fontsize:N:text}` / `{d:N:text}`         | text        |
| `{fontsizepercent:N:text}` / `{e:N:text}`  | text        |
| `{characterspacing:N:text}` / `{m:N:text}` | text        |
| `{border:N:text}` / `{o:N:text}`           | text        |
| `{shadow:X,Y:text}` / `{s:X,Y:text}`       | text        |
| `{shadowcolor:HEX:text}` / `{v:HEX:text}`  | text        |
| `{bordercolor:HEX:text}` / `{r:HEX:text}`  | text        |
| `{width:text}` / `{w:text}`                | text        |
| `{loghint:hint:text}` / `{l:hint:text}`    | text        |
| `{a:param:text}` (alignment)               | text        |
| `{n:N:text}` (conditional, default shown)  | text        |
| `{y:N:text}` (conditional, not default)    | *(removed)* |
| Any other `{Tag:...:text}`                 | text        |

### Special character tags

These are replaced before other processing.

| Tag                  | Replacement                            |
|----------------------|----------------------------------------|
| `{0}`                | *(zero-width space, removed)*          |
| `{-}`                | *(soft hyphen, removed)*               |
| `{qt}`               | `"`                                    |
| `{ob}` / `{eb}`      | *(removed, stray braces are stripped)* |
| `{os}` / `{es}`      | `[` / `]`                              |
| `{t}` / `{parallel}` | *(parallel display marker, removed)*   |

### Other cleanup

- Backticks (`` ` ``), inline commands (`[@]`, `[\]`, `[|]`), and voice metadata (`[lv ...]`) are stripped
- `{Comment:...}` translator notes are stripped entirely
- Any remaining `{` or `}` are stripped after all tag processing (catches stray braces from tags that span across backtick segments, e.g. `{p:1:` red truth split across voice lines)

## Contributors

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/HannahBanana1312">
        <img src="https://avatars.githubusercontent.com/u/36461227?v=4" width="100px;" alt="Hannah"/><br />
        <sub><b>Hannah</b></sub>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/nakedmcse">
        <img src="https://avatars.githubusercontent.com/u/133156975?v=4" width="100px;" alt="Walker Boh"/><br />
        <sub><b>Walker Boh</b></sub>
      </a>
    </td>
  </tr>
</table>
