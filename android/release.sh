#!/bin/bash

set -e

if [ -z "$JAVA_HOME" ]; then
    echo "Error: JAVA_HOME is not set. Set it to your JDK installation path."
    echo "Example: export JAVA_HOME=C:/Users/you/.jdks/temurin-24.0.1"
    exit 1
fi

export JAVA_OPTS="--enable-native-access=ALL-UNNAMED --add-opens=java.base/java.lang=ALL-UNNAMED --add-opens=java.base/sun.misc=ALL-UNNAMED"

cd "$(dirname "$0")"

kotlin -J--enable-native-access=ALL-UNNAMED -J--add-opens=java.base/java.lang=ALL-UNNAMED -J--add-opens=java.base/sun.misc=ALL-UNNAMED release.main.kts

echo ""
echo "[2/3] Running Gradle clean..."
./gradlew clean
echo "✓ Clean completed"

echo ""
echo "[3/3] Running Gradle assembleRelease..."
./gradlew assembleRelease
echo "✓ APK build completed"

VERSION=$(grep -oP 'versionName\s*=\s*"\K[^"]+' app/build.gradle.kts)
echo ""
echo "╔════════════════════════════════════════╗"
echo "║        BUILD SUCCESSFUL! ✓             ║"
echo "╚════════════════════════════════════════╝"
echo "APK: app/build/outputs/apk/release/umineko-quotes-v${VERSION}.apk"
