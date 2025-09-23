#!/usr/bin/env bash
# Development testing script - quick iteration for developers
set -euo pipefail

ROOT=$(cd "$(dirname "$0")/.." && pwd)
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

echo "[dev-test] Development testing script"
echo "[dev-test] Working in temporary directory: $TEMP_DIR"

cd "$ROOT"

# Quick build and basic test
echo "[dev-test] Building launcher..."
make build-windows >/dev/null 2>&1

echo "[dev-test] Building test artifacts..."
make build-stubs >/dev/null 2>&1
make build-testapp >/dev/null 2>&1

# Copy to temp dir for testing
cp bin/gjg-launcher.exe "$TEMP_DIR/test-launcher.exe"
cp -r testdata/jre "$TEMP_DIR/"
cp testdata/java/testapp.jar "$TEMP_DIR/"

cd "$TEMP_DIR"

echo "[dev-test] Testing basic functionality..."

# Test 1: Dry run
echo "[dev-test] Test 1: Dry run"
cat > test-launcher.conf <<EOF
java_dir=./jre
java_executable=java
jar_file=testapp.jar
jvm_args=-Xmx128m -Dtest.mode=dev
app_args=--test-arg value
env_DEV_TEST=true
EOF

echo "  Running dry run..."
if ! "$ROOT/scripts/wine-run.sh" ./test-launcher.exe --gjg-dry-run | grep -q "java"; then
    echo "  ❌ Dry run test failed"
    exit 1
fi
echo "  ✅ Dry run test passed"

# Test 2: Debug output
echo "[dev-test] Test 2: Debug output"
echo "  Running debug..."
if ! "$ROOT/scripts/wine-run.sh" ./test-launcher.exe --gjg-debug 2>/dev/null | grep -q "Configuration loaded"; then
    echo "  ❌ Debug test failed"
    exit 1
fi
echo "  ✅ Debug test passed"

# Test 3: Execution
echo "[dev-test] Test 3: Execution"
echo "  Running execution..."
OUTPUT=$("$ROOT/scripts/wine-run.sh" ./test-launcher.exe --dev-mode 2>/dev/null | tail -n 10)
if ! echo "$OUTPUT" | grep -q "GJG Test Application"; then
    echo "  ❌ Execution test failed"
    echo "  Output: $OUTPUT"
    exit 1
fi
echo "  ✅ Execution test passed"

# Test 4: Error handling
echo "[dev-test] Test 4: Error handling"
echo "  Testing missing JAR..."
cat > test-launcher.conf <<EOF
java_dir=./jre
jar_file=nonexistent.jar
EOF

if "$ROOT/scripts/wine-run.sh" ./test-launcher.exe 2>/dev/null; then
    echo "  ❌ Error handling test failed (should have failed)"
    exit 1
fi
echo "  ✅ Error handling test passed"

echo "[dev-test] All tests passed! ✅"
echo "[dev-test] Quick development iteration complete"