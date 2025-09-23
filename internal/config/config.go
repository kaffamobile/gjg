// Package config provides configuration file parsing and management for the GJG launcher.
package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	JavaDir        string
	JavaExecutable string
	JarFile        string
	JVMArgs        string
	AppArgs        string
	Env            map[string]string
}

// LoadFromExe locates the .conf file next to the executable and loads it.
func LoadFromExe() (Config, string, error) {
	exe, err := os.Executable()
	if err != nil {
		return Config{}, "", fmt.Errorf("failed to get executable path: %w", err)
	}
	base := strings.TrimSuffix(exe, filepath.Ext(exe))
	confPath := base + ".conf"
	cfg, err := Load(confPath, filepath.Base(base)+".jar")
	if err != nil {
		return Config{}, "", err
	}
	return cfg, confPath, nil
}

// Load parses a .conf file at the given path. defaultJar is used when jar_file not specified.
func Load(path, defaultJar string) (Config, error) {
	// #nosec G304 - This is intentional file inclusion via variable, path is from executable name
	f, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("configuration file error: %w", err)
	}
	defer f.Close()

	cfg := Config{
		JavaExecutable: "java",
		JarFile:        defaultJar,
		Env:            map[string]string{},
	}

	scanner := bufio.NewScanner(f)
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.IndexRune(line, '=')
		if eq <= 0 { // key cannot be empty
			return Config{}, fmt.Errorf("invalid config line %d: %q", lineNo, line)
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])

		if strings.HasPrefix(key, "env_") {
			envKey := strings.TrimPrefix(key, "env_")
			if envKey == "" {
				return Config{}, fmt.Errorf("invalid env_ key at line %d", lineNo)
			}
			cfg.Env[envKey] = val
			continue
		}

		switch key {
		case "java_dir":
			cfg.JavaDir = val
		case "java_executable":
			cfg.JavaExecutable = val
		case "jar_file":
			cfg.JarFile = val
		case "jvm_args":
			cfg.JVMArgs = val
		case "app_args":
			cfg.AppArgs = val
		default:
			return Config{}, fmt.Errorf("unknown config key %q at line %d", key, lineNo)
		}
	}
	if err := scanner.Err(); err != nil {
		return Config{}, fmt.Errorf("error reading config: %w", err)
	}

	return cfg, nil
}

// MergeEnv overlays overrides onto base environment in KEY=VAL form.
func MergeEnv(base []string, overrides map[string]string) []string {
	out := make([]string, 0, len(base)+len(overrides))
	keys := map[string]int{}
	for i, kv := range base {
		eq := strings.IndexRune(kv, '=')
		if eq <= 0 {
			continue
		}
		k := kv[:eq]
		keys[strings.ToUpper(k)] = i
		out = append(out, kv)
	}
	for k, v := range overrides {
		key := strings.ToUpper(k)
		kv := k + "=" + v
		if idx, ok := keys[key]; ok {
			out[idx] = kv
		} else {
			out = append(out, kv)
		}
	}
	return out
}

// EnvSummary formats environment overrides as "K=V, K2=V2" for debug.
func EnvSummary(env map[string]string) string {
	if len(env) == 0 {
		return ""
	}
	parts := make([]string, 0, len(env))
	for k, v := range env {
		parts = append(parts, fmt.Sprintf("%s=%s", k, v))
	}
	// Stable-ish order not guaranteed; this is acceptable for debugging purposes.
	return strings.Join(parts, ", ")
}

// Helpers for explicit error kinds
var (
	ErrConfig = errors.New("config error")
)
