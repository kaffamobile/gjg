#!/usr/bin/env bash
set -euo pipefail

ROOT=$(cd "$(dirname "$0")" && pwd)
SRC_DIR="$ROOT/src"
BUILD_DIR="$ROOT/.build"
CLASSES_DIR="$BUILD_DIR/classes"
JAR_PATH="$ROOT/myapp.jar"
MAIN_CLASS="hello.Main"

rm -rf "$BUILD_DIR"
mkdir -p "$CLASSES_DIR"

echo "[build-jar] Compiling Java sources"
javac -d "$CLASSES_DIR" $(find "$SRC_DIR" -name "*.java")

echo "[build-jar] Writing manifest"
mkdir -p "$BUILD_DIR/META-INF"
MANIFEST="$BUILD_DIR/META-INF/MANIFEST.MF"
cat > "$MANIFEST" <<EOF
Manifest-Version: 1.0
Main-Class: $MAIN_CLASS

EOF

echo "[build-jar] Creating JAR at $JAR_PATH"
jar cfm "$JAR_PATH" "$MANIFEST" -C "$CLASSES_DIR" .

echo "[build-jar] Done"

