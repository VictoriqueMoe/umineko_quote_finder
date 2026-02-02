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
- [Script Tag Parsing](#script-tag-parsing)
  - [Tags with HTML rendering](#tags-with-html-rendering)
  - [Preset colour reference](#preset-colour-reference)
  - [Tags stripped to content](#tags-stripped-to-content)
  - [Special character tags](#special-character-tags)
  - [Other cleanup](#other-cleanup)

## Features

- Fuzzy search through all dialogue
- Filter by character and episode
- Random quote generator
- Scene context viewer — see surrounding dialogue for any quote
- English/Japanese language toggle
- Inline audio playback for voiced lines
- Umineko-themed web interface

## Quick Start

```bash
go build -o umineko_quote.exe .
./umineko_quote.exe
```

Open http://127.0.0.1:3000

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

The app works without audio files — quotes will display normally but without playback controls.

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
| `GET /api/v1/search`                 | Fuzzy search quotes                    |
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

**Menu/UI presets** (not rendered — stripped to plain text if they appear):

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

| Tag                  | Replacement                             |
|----------------------|-----------------------------------------|
| `{0}`                | *(zero-width space, removed)*           |
| `{-}`                | *(soft hyphen, removed)*                |
| `{qt}`               | `"`                                     |
| `{ob}` / `{eb}`      | *(removed — stray braces are stripped)* |
| `{os}` / `{es}`      | `[` / `]`                               |
| `{t}` / `{parallel}` | *(parallel display marker, removed)*    |

### Other cleanup

- Backticks (`` ` ``), inline commands (`[@]`, `[\]`, `[|]`), and voice metadata (`[lv ...]`) are stripped
- `{Comment:...}` translator notes are stripped entirely
- Any remaining `{` or `}` are stripped after all tag processing (catches stray braces from tags that span across backtick segments, e.g. `{p:1:` red truth split across voice lines)
