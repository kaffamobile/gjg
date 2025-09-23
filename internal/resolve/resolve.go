package resolve

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
)

// ResolveJava returns the path to the java executable per rules.
func ResolveJava(javaDir string, javaExe string, exeDir string) (string, error) {
    exeName := javaExe
    if !strings.HasSuffix(strings.ToLower(exeName), ".exe") {
        exeName += ".exe"
    }
    if strings.TrimSpace(javaDir) == "" {
        // PATH lookup
        p, err := exec.LookPath(exeName)
        if err != nil {
            return "", fmt.Errorf("java executable not found in PATH: %s", exeName)
        }
        return p, nil
    }
    // resolve relative to exeDir
    base := javaDir
    if !filepath.IsAbs(base) {
        base = filepath.Join(exeDir, base)
    }
    p := filepath.Join(base, "bin", exeName)
    if _, err := os.Stat(p); err != nil {
        return "", fmt.Errorf("java executable not found: %s", p)
    }
    return p, nil
}

// ResolveJar resolves jar file path relative to exeDir if needed and verifies existence.
func ResolveJar(jar string, exeDir string) (string, error) {
    p := jar
    if !filepath.IsAbs(p) {
        p = filepath.Join(exeDir, p)
    }
    info, err := os.Stat(p)
    if err != nil {
        return "", fmt.Errorf("jar file not found: %s", p)
    }
    if info.IsDir() {
        return "", fmt.Errorf("jar path is a directory: %s", p)
    }
    return p, nil
}

