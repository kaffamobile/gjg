#!/usr/bin/env bash
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/.." && pwd)

"$ROOT/scripts/wine-check.sh" >/dev/null

echo "[smoke] Building launcher (windows/amd64)"
mkdir -p "$ROOT/bin"
GOOS=windows GOARCH=amd64 go build -o "$ROOT/bin/gjg-launcher.exe" ./cmd/launcher

echo "[smoke] Building stub java.exe and javaw.exe"
mkdir -p "$ROOT/testdata/jre/bin"
GOOS=windows GOARCH=amd64 go build -o "$ROOT/testdata/jre/bin/java.exe" ./testdata/stubs
cp "$ROOT/testdata/jre/bin/java.exe" "$ROOT/testdata/jre/bin/javaw.exe"
echo "[smoke] Copying stubs into example tree"
mkdir -p "$ROOT/examples/wine/testdata/jre/bin"
cp "$ROOT/testdata/jre/bin/"*.exe "$ROOT/examples/wine/testdata/jre/bin/"

echo "[smoke] Preparing example app"
EX="$ROOT/examples/wine"
mkdir -p "$EX"
cat > "$EX/myapp.conf" <<'CONF'
java_dir=./testdata/jre
java_executable=javaw
jar_file=myapp.jar
jvm_args=-Xmx512m "-Djava.library.path=./libs"
app_args=--config myapp.properties
env_MY_HOME=/home/user/workspace
env_DEBUG_MODE=true
CONF

echo "[smoke] Building a valid myapp.jar"
bash "$EX/build-jar.sh"

echo "[smoke] Preparing launcher in example dir"
cp "$ROOT/bin/gjg-launcher.exe" "$EX/myapp.exe"

echo "[smoke] Running dry-run"
"$ROOT/scripts/wine-run.sh" "$EX/myapp.exe" --gjg-dry-run | sed -n '1,200p'

echo "[smoke] Running debug"
"$ROOT/scripts/wine-run.sh" "$EX/myapp.exe" --gjg-debug | sed -n '1,200p'

echo "[smoke] Launching stubbed java to verify execution"
pushd "$EX" >/dev/null
"$ROOT/scripts/wine-run.sh" "$EX/myapp.exe" --foo bar | sed -n '1,200p'
popd >/dev/null

echo "[smoke] Preparing second app (console)"
cat > "$EX/myapp-console.conf" <<'CONF'
java_dir=./testdata/jre
java_executable=java
jar_file=myapp.jar
jvm_args=-Xmx256m
app_args=--mode console
env_MY_HOME=/home/user/workspace
env_DEBUG_MODE=false
CONF

cp "$ROOT/bin/gjg-launcher.exe" "$EX/myapp-console.exe"

echo "[smoke] Running dry-run (console)"
"$ROOT/scripts/wine-run.sh" "$EX/myapp-console.exe" --gjg-dry-run | sed -n '1,200p'

echo "[smoke] Running debug (console)"
"$ROOT/scripts/wine-run.sh" "$EX/myapp-console.exe" --gjg-debug | sed -n '1,200p'

echo "[smoke] Launching stubbed java (console)"
pushd "$EX" >/dev/null
"$ROOT/scripts/wine-run.sh" "$EX/myapp-console.exe" --foo bar | sed -n '1,200p'
popd >/dev/null

echo "[smoke] Done"
