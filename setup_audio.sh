#!/bin/sh
set -e

AUDIO_DIR="internal/quote/data/audio"

if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

ZIP_SOURCE="${VOICE_ZIP_URL:?VOICE_ZIP_URL is not set. Create a .env file with VOICE_ZIP_URL=<url or path>}"

if [ -d "$AUDIO_DIR" ]; then
    echo "Audio directory already exists at $AUDIO_DIR, skipping download."
    exit 0
fi

if [ -f "$ZIP_SOURCE" ]; then
    echo "Extracting from local file: $ZIP_SOURCE"
    mkdir -p internal/quote/data
    unzip -qo "$ZIP_SOURCE" -d /tmp/voice
    mv /tmp/voice/voice "$AUDIO_DIR"
    rm -rf /tmp/voice
else
    echo "Downloading voice files..."
    curl -fSL -o /tmp/voice.zip "$ZIP_SOURCE"
    echo "Extracting..."
    mkdir -p internal/quote/data
    unzip -qo /tmp/voice.zip -d /tmp/voice
    mv /tmp/voice/voice "$AUDIO_DIR"
    rm -rf /tmp/voice.zip /tmp/voice
fi

echo "Done. Audio files extracted to $AUDIO_DIR"

SE_DIR="internal/quote/data/se"
SE_SOURCE="${SE_ZIP_URL:-}"

if [ -z "$SE_SOURCE" ]; then
    echo "SE_ZIP_URL is not set, skipping SE download."
    exit 0
fi

if [ -d "$SE_DIR" ]; then
    echo "SE directory already exists at $SE_DIR, skipping download."
    exit 0
fi

if [ -f "$SE_SOURCE" ]; then
    echo "Extracting SE from local file: $SE_SOURCE"
    mkdir -p internal/quote/data
    unzip -qo "$SE_SOURCE" -d /tmp/se
    mv /tmp/se/se "$SE_DIR"
    rm -rf /tmp/se
else
    echo "Downloading SE files..."
    curl -fSL -o /tmp/se.zip "$SE_SOURCE"
    echo "Extracting..."
    mkdir -p internal/quote/data
    unzip -qo /tmp/se.zip -d /tmp/se
    mv /tmp/se/se "$SE_DIR"
    rm -rf /tmp/se.zip /tmp/se
fi

echo "Done. SE files extracted to $SE_DIR"
