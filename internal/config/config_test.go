package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *Config
		wantErr  bool
	}{
		{
			name: "basic config",
			input: `java_dir=./jre
java_executable=javaw
jar_file=app.jar
jvm_args=-Xmx512m
app_args=--verbose
env_TEST_VAR=test_value`,
			expected: &Config{
				JavaDir:        "./jre",
				JavaExecutable: "javaw",
				JarFile:        "app.jar",
				JVMArgs:        "-Xmx512m",
				AppArgs:        "--verbose",
				Env: map[string]string{
					"TEST_VAR": "test_value",
				},
			},
			wantErr: false,
		},
		{
			name: "minimal config with defaults",
			input: `jvm_args=-Xmx256m`,
			expected: &Config{
				JavaDir:        "",
				JavaExecutable: "java",
				JarFile:        "test.jar", // will be set by test
				JVMArgs:        "-Xmx256m",
				AppArgs:        "",
				Env:            map[string]string{},
			},
			wantErr: false,
		},
		{
			name: "config with comments and empty lines",
			input: `# This is a comment
java_dir=./jre

# Another comment
java_executable=java
jar_file=myapp.jar

env_HOME=/home/user
env_DEBUG=true`,
			expected: &Config{
				JavaDir:        "./jre",
				JavaExecutable: "java",
				JarFile:        "myapp.jar",
				JVMArgs:        "",
				AppArgs:        "",
				Env: map[string]string{
					"HOME":  "/home/user",
					"DEBUG": "true",
				},
			},
			wantErr: false,
		},
		{
			name: "config with quoted arguments",
			input: `jvm_args=-Xmx512m "-Djava.library.path=./libs with spaces"
app_args=--config "my config.properties" --verbose`,
			expected: &Config{
				JavaDir:        "",
				JavaExecutable: "java",
				JarFile:        "test.jar", // will be set by test
				JVMArgs:        `-Xmx512m "-Djava.library.path=./libs with spaces"`,
				AppArgs:        `--config "my config.properties" --verbose`,
				Env:            map[string]string{},
			},
			wantErr: false,
		},
		{
			name:     "invalid config - no equals sign",
			input:    `java_dir./jre`,
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "empty config",
			input:    "",
			expected: &Config{
				JavaDir:        "",
				JavaExecutable: "java",
				JarFile:        "test.jar", // will be set by test
				JVMArgs:        "",
				AppArgs:        "",
				Env:            map[string]string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary config file
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "test.conf")
			err := os.WriteFile(configPath, []byte(tt.input), 0644)
			if err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			got, err := Load(configPath, "test.jar")
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Set default jar file based on config name if not set
			if tt.expected.JarFile == "test.jar" {
				tt.expected.JarFile = "test.jar"
			}

			if !reflect.DeepEqual(got, *tt.expected) {
				t.Errorf("Load() = %+v, want %+v", got, *tt.expected)
			}
		})
	}
}

func TestMergeEnv(t *testing.T) {
	baseEnv := []string{
		"PATH=/usr/bin",
		"HOME=/home/user",
		"EXISTING=original",
	}

	configEnv := map[string]string{
		"NEW_VAR":  "new_value",
		"EXISTING": "overridden",
		"EMPTY":    "",
	}

	result := MergeEnv(baseEnv, configEnv)

	// Convert result back to map for easier testing
	resultMap := make(map[string]string)
	for _, env := range result {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) == 2 {
			resultMap[parts[0]] = parts[1]
		}
	}

	// Check that original env vars are preserved
	if resultMap["PATH"] != "/usr/bin" {
		t.Errorf("Expected PATH=/usr/bin, got %s", resultMap["PATH"])
	}

	if resultMap["HOME"] != "/home/user" {
		t.Errorf("Expected HOME=/home/user, got %s", resultMap["HOME"])
	}

	// Check that config env vars are added/override
	if resultMap["NEW_VAR"] != "new_value" {
		t.Errorf("Expected NEW_VAR=new_value, got %s", resultMap["NEW_VAR"])
	}

	if resultMap["EXISTING"] != "overridden" {
		t.Errorf("Expected EXISTING=overridden, got %s", resultMap["EXISTING"])
	}

	if resultMap["EMPTY"] != "" {
		t.Errorf("Expected EMPTY=, got %s", resultMap["EMPTY"])
	}
}

func TestEnvSummary(t *testing.T) {
	tests := []struct {
		name     string
		env      map[string]string
		expected string
	}{
		{
			name:     "empty env",
			env:      map[string]string{},
			expected: "",
		},
		{
			name: "single env var",
			env: map[string]string{
				"TEST": "value",
			},
			expected: "TEST=value",
		},
		{
			name: "multiple env vars",
			env: map[string]string{
				"VAR1": "value1",
				"VAR2": "value2",
			},
			expected: "VAR1=value1, VAR2=value2", // Note: order may vary
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := EnvSummary(tt.env)

			if tt.name == "empty env" {
				if result != tt.expected {
					t.Errorf("EnvSummary() = %q, want %q", result, tt.expected)
				}
				return
			}

			if tt.name == "single env var" {
				if result != tt.expected {
					t.Errorf("EnvSummary() = %q, want %q", result, tt.expected)
				}
				return
			}

			// For multiple vars, check that all expected parts are present
			if tt.name == "multiple env vars" {
				if !strings.Contains(result, "VAR1=value1") {
					t.Errorf("EnvSummary() missing VAR1=value1 in %q", result)
				}
				if !strings.Contains(result, "VAR2=value2") {
					t.Errorf("EnvSummary() missing VAR2=value2 in %q", result)
				}
				if !strings.Contains(result, ", ") {
					t.Errorf("EnvSummary() missing separator in %q", result)
				}
			}
		})
	}
}