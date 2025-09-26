//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var Default = Build

func getVersion() string {
	if v := os.Getenv("VERSION"); v != "" {
		return v
	}

	cmd := exec.Command("git", "describe", "--exact-match", "--tags", "HEAD")
	if out, err := cmd.Output(); err == nil {
		tag := strings.TrimSpace(string(out))
		if strings.HasPrefix(tag, "v") {
			return tag
		}
	}

	return "dev"
}

func getGitVersion() string {
	cmd := exec.Command("git", "describe", "--exact-match", "--tags", "HEAD")
	if out, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(out))
	}

	cmd = exec.Command("git", "describe", "--tags", "--always", "--dirty")
	if out, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(out))
	}

	cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	if out, err := cmd.Output(); err == nil {
		return "git-" + strings.TrimSpace(string(out))
	}

	return ""
}

func Build() error {
	version := getVersion()
	fmt.Printf("Building GJG Launcher v%s for all Windows architectures...\n\n", version)

	targets := []struct {
		goarch string
		output string
		desc   string
	}{
		{"amd64", "gjg-launcher-windows-amd64.exe", "Windows 64-bit"},
		{"386", "gjg-launcher-windows-386.exe", "Windows 32-bit"},
		{"arm64", "gjg-launcher-windows-arm64.exe", "Windows ARM64"},
	}

	if err := os.MkdirAll("bin", 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	for _, target := range targets {
		fmt.Printf("Building %s...\n", target.desc)

		ldflags := fmt.Sprintf("-H windowsgui -X main.version=%s", version)

		cmd := exec.Command("go", "build",
			"-ldflags", ldflags,
			"-trimpath",
			"-o", "./bin/"+target.output,
			"./cmd/launcher/main.go")

		cmd.Env = append(os.Environ(),
			"GOOS=windows",
			"GOARCH="+target.goarch,
			"CGO_ENABLED=0",
		)

		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to build %s: %w\n%s", target.desc, err, output)
		}

		if info, err := os.Stat("bin/" + target.output); err == nil {
			sizeMB := float64(info.Size()) / 1024 / 1024
			fmt.Printf("✓ %s (%.2f MB)\n", target.output, sizeMB)
		} else {
			fmt.Printf("✓ %s\n", target.output)
		}
	}

	fmt.Printf("\n✅ Build completed successfully!\n")
	fmt.Printf("📁 Binaries available in: bin/\n")
	fmt.Printf("📌 Version: %s\n", version)

	return nil
}

func Test() error {
	fmt.Println("Running tests...")

	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("tests failed: %w", err)
	}

	fmt.Println("✓ All tests passed!")
	return nil
}

func Fmt() error {
	fmt.Println("Formatting code...")

	cmd := exec.Command("go", "fmt", "./...")
	output, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf("formatting failed: %w", err)
	}

	if len(output) > 0 {
		fmt.Printf("Formatted files:\n%s", output)
	} else {
		fmt.Println("✓ All files already formatted")
	}

	return nil
}

func Vet() error {
	fmt.Println("Running go vet...")

	cmd := exec.Command("go", "vet", "./...")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("vet failed: %w\n%s", err, output)
	}

	fmt.Println("✓ No issues found")
	return nil
}

func Check() error {
	fmt.Println("Running all checks...")
	fmt.Println("=" + strings.Repeat("=", 40))

	// Format
	if err := Fmt(); err != nil {
		return err
	}
	fmt.Println()

	// Vet
	if err := Vet(); err != nil {
		return err
	}
	fmt.Println()

	// Test
	if err := Test(); err != nil {
		return err
	}

	fmt.Println("=" + strings.Repeat("=", 40))
	fmt.Println("✅ All checks passed!")
	return nil
}

func Clean() error {
	fmt.Println("Cleaning build artifacts...")

	if err := os.RemoveAll("bin"); err != nil {
		return fmt.Errorf("failed to remove bin: %w", err)
	}

	patterns := []string{"*.out", "*.test", "*.exe"}
	for _, pattern := range patterns {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("rm -f %s", pattern))
		cmd.Run() // Ignora erros se não existir
	}

	fmt.Println("✓ Cleaned successfully")
	return nil
}

func Version() error {
	v := getVersion()
	fmt.Printf("Current version: %s\n", v)

	// Mostra informações extras se disponível
	if cmd := exec.Command("git", "status", "--porcelain"); cmd != nil {
		if out, _ := cmd.Output(); len(out) > 0 {
			fmt.Println("⚠ Working directory has uncommitted changes")
		}
	}

	return nil
}
