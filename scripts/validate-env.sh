#!/usr/bin/env bash
# Environment validation script
set -euo pipefail

echo "=== GJG Development Environment Validation ==="
echo

SUCCESS=true

# Check Go
echo "Checking Go..."
if command -v go >/dev/null 2>&1; then
    GO_VERSION=$(go version | grep -o 'go[0-9.]*' | head -1)
    echo "  ✅ Go found: $GO_VERSION"

    # Check Go version (need 1.20+)
    if [[ "$GO_VERSION" < "go1.20" ]]; then
        echo "  ⚠️  Warning: Go 1.20+ recommended, found $GO_VERSION"
    fi
else
    echo "  ❌ Go not found"
    SUCCESS=false
fi

echo

# Check Java
echo "Checking Java..."
if command -v java >/dev/null 2>&1; then
    JAVA_VERSION=$(java -version 2>&1 | head -1)
    echo "  ✅ Java runtime found: $JAVA_VERSION"
else
    echo "  ❌ Java runtime not found"
    SUCCESS=false
fi

if command -v javac >/dev/null 2>&1; then
    JAVAC_VERSION=$(javac -version 2>&1)
    echo "  ✅ Java compiler found: $JAVAC_VERSION"

    # Test Java 8 compatibility
    TEMP_DIR=$(mktemp -d)
    cat > "$TEMP_DIR/Test.java" <<EOF
public class Test {
    public static void main(String[] args) {
        System.out.println("Java 8 compatibility test");
    }
}
EOF

    if javac -source 8 -target 8 "$TEMP_DIR/Test.java" 2>/dev/null; then
        echo "  ✅ Java 8 compatibility confirmed"
    else
        echo "  ⚠️  Warning: Java 8 compatibility test failed"
    fi

    rm -rf "$TEMP_DIR"
else
    echo "  ❌ Java compiler (javac) not found"
    SUCCESS=false
fi

echo

# Check Wine (optional for Linux)
echo "Checking Wine..."
if command -v wine64 >/dev/null 2>&1; then
    WINE_VERSION=$(wine64 --version 2>/dev/null || echo "unknown")
    echo "  ✅ Wine found: $WINE_VERSION"

    # Check Wine configuration
    export WINEDEBUG=-all
    if wine64 cmd /c echo "Wine test" >/dev/null 2>&1; then
        echo "  ✅ Wine is functional"
    else
        echo "  ⚠️  Warning: Wine may need initialization (run 'winecfg' first)"
    fi
else
    echo "  ⚠️  Wine not found (Linux Wine testing will not work)"
fi

echo

# Check Make
echo "Checking Make..."
if command -v make >/dev/null 2>&1; then
    MAKE_VERSION=$(make --version | head -1)
    echo "  ✅ Make found: $MAKE_VERSION"
else
    echo "  ❌ Make not found"
    SUCCESS=false
fi

echo

# Check Git
echo "Checking Git..."
if command -v git >/dev/null 2>&1; then
    GIT_VERSION=$(git --version)
    echo "  ✅ Git found: $GIT_VERSION"
else
    echo "  ⚠️  Git not found (version info will not work)"
fi

echo

# Check optional tools
echo "Checking optional tools..."

if command -v golangci-lint >/dev/null 2>&1; then
    LINT_VERSION=$(golangci-lint version 2>/dev/null | head -1 || echo "unknown")
    echo "  ✅ golangci-lint found: $LINT_VERSION"
else
    echo "  ⚠️  golangci-lint not found (will use 'go vet' instead)"
    echo "     Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi

echo

# Summary
echo "=== Summary ==="
if [ "$SUCCESS" = true ]; then
    echo "✅ Environment validation successful!"
    echo "You can run: make dev-setup && make test"
else
    echo "❌ Environment validation failed!"
    echo "Please install missing dependencies before continuing."
    exit 1
fi

echo
echo "Quick start commands:"
echo "  make help          - Show available commands"
echo "  make dev-setup     - Verify and set up environment"
echo "  make test          - Run all tests"
echo "  make build         - Build launcher"
echo "  scripts/dev-test.sh - Quick development test"