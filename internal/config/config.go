package config

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

type Config struct {
	JavaExecutableAbsolutePath string
	JarFileAbsolutePath        string
	JVMArgs                    string
	AppArgs                    string
	Env                        []string
}

func Load() (*Config, string, error) {
	exe, err := os.Executable()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get executable path: %w", err)
	}

	exeBase := strings.TrimSuffix(filepath.Base(exe), filepath.Ext(exe))
	exeDir := filepath.Dir(exe)
	searchPaths := []string{
		filepath.Join(exeDir, exeBase+".gjg.conf"),
		filepath.Join(".", exeBase+".gjg.conf"),
		filepath.Join(exeDir, "example.gjg.conf"),
		filepath.Join(".", "example.gjg.conf"),
	}

	var confFilePath string
	for _, path := range searchPaths {
		if _, err := os.Stat(path); err == nil {
			confFilePath = path
			break
		}
	}
	if confFilePath == "" {
		return nil, "", fmt.Errorf("configuration file not found. Searched for: %v", searchPaths)
	}
	confFilePath, err = filepath.Abs(confFilePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get absolute path of configuration file: %w", err)
	}

	cfg, err := buildConfig(confFilePath)
	if err != nil {
		return nil, "", err
	}

	return cfg, confFilePath, nil
}

func buildConfig(configFilePath string) (*Config, error) {
	f, err := os.Open(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("configuration file error: %w", err)
	}
	defer f.Close()

	cfg := &Config{
		Env: os.Environ(),
	}

	var javaDir string
	var jarFile string
	envOverrides := make(map[string]string)

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
		if eq <= 0 {
			return nil, fmt.Errorf("invalid config line %d: %q", lineNo, line)
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])

		switch {
		case strings.HasPrefix(key, "env_"):
			envKey := strings.TrimPrefix(key, "env_")
			if envKey == "" {
				return nil, fmt.Errorf("invalid env_ key at line %d", lineNo)
			}
			envOverrides[envKey] = val
		case key == "java_dir":
			javaDir = val
		case key == "jar_file":
			jarFile = val
		case key == "jvm_args":
			cfg.JVMArgs = val
		case key == "app_args":
			cfg.AppArgs = val
		default:
			return nil, fmt.Errorf("unknown config key %q at line %d", key, lineNo)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading config: %w", err)
	}

	if jarFile == "" {
		exeBase := strings.TrimSuffix(filepath.Base(configFilePath), ".gjg.conf")
		jarFile = exeBase + ".jar"
	}

	configDir := filepath.Dir(configFilePath)
	javaPath, err := resolveJava(javaDir, configDir)
	if err != nil {
		return nil, fmt.Errorf("java resolution failed: %w", err)
	}
	cfg.JavaExecutableAbsolutePath = javaPath

	jarPath, err := resolveJar(jarFile, configDir)
	if err != nil {
		return nil, fmt.Errorf("jar resolution failed: %w", err)
	}
	cfg.JarFileAbsolutePath = jarPath
	cfg.Env = mergeEnv(cfg.Env, envOverrides)

	return cfg, nil
}

func mergeEnv(base []string, overrides map[string]string) []string {
	result := make([]string, 0, len(base)+len(overrides))
	seen := make(map[string]bool)

	for _, env := range base {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToUpper(parts[0])
		if newVal, ok := overrides[parts[0]]; ok {
			result = append(result, parts[0]+"="+newVal)
			seen[key] = true
		} else {
			result = append(result, env)
		}
	}
	for k, v := range overrides {
		if !seen[strings.ToUpper(k)] {
			result = append(result, k+"="+v)
		}
	}

	return result
}

func resolveJava(javaDir, configDir string) (string, error) {
	exeName := "javaw.exe"
	if runtime.GOOS != "windows" {
		exeName = "java"
	}

	if strings.TrimSpace(javaDir) == "" {
		if p, err := exec.LookPath(exeName); err == nil {
			return filepath.Abs(p)
		}
		return "", fmt.Errorf("java executable not found in PATH and not set on conf file")
	}

	base := javaDir
	if !filepath.IsAbs(base) {
		base = filepath.Join(configDir, base)
	}

	javaPath := filepath.Join(base, "bin", exeName)
	if _, err := os.Stat(javaPath); err != nil {
		return "", fmt.Errorf("java executable not found: %s", javaPath)
	}

	return filepath.Abs(javaPath)
}

func resolveJar(jar, configDir string) (string, error) {
	p := jar
	if !filepath.IsAbs(p) {
		p = filepath.Join(configDir, p)
	}
	info, err := os.Stat(p)
	if err != nil {
		return "", fmt.Errorf("jar file not found: %s", p)
	}
	if info.IsDir() {
		return "", fmt.Errorf("jar path is a directory: %s", p)
	}
	return filepath.Abs(p)
}
