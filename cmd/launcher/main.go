package main

import (
	"fmt"
	"gjg/internal/args"
	"gjg/internal/config"
	"gjg/internal/runner"
	"os"
	"path/filepath"
	"time"
)

var version = "dev"

func main() {
	debug, dryRun, forwardArgs := args.ExtractSpecial(os.Args[1:])

	var logFile *os.File
	if debug {
		exe, err := os.Executable()
		if err == nil {
			logFile, _ = os.Create(filepath.Join(filepath.Dir(exe), "gjg-debug.log"))
			if logFile != nil {
				defer logFile.Close()
			}
			logf(logFile, "Starting Launcher on Version: %s", version)
		}
	}

	cfg, confPath, err := config.Load()
	if err != nil {
		logf(logFile, "Error loading config: %s", err)
		os.Exit(1)
	}

	jvmTokens := args.Tokenize(cfg.JVMArgs)
	appTokens := args.Tokenize(cfg.AppArgs)

	argv := make([]string, 0, 4+len(jvmTokens)+len(appTokens)+len(forwardArgs))
	argv = append(argv, cfg.JavaExecutableAbsolutePath)
	argv = append(argv, jvmTokens...)
	argv = append(argv, "-jar", cfg.JarFileAbsolutePath)
	argv = append(argv, appTokens...)
	argv = append(argv, forwardArgs...)

	if debug {
		logf(logFile, "Configuration loaded from: %s", confPath)
		logf(logFile, "Java executable: %s", cfg.JavaExecutableAbsolutePath)
		logf(logFile, "JAR file: %s", cfg.JarFileAbsolutePath)
		logf(logFile, "Working directory: %s", filepath.Dir(confPath))

		if cfg.JVMArgs != "" {
			logf(logFile, "JVM arguments: %s", cfg.JVMArgs)
		}

		if len(forwardArgs) > 0 {
			logf(logFile, "Forward arguments: %v", forwardArgs)
		}

		logf(logFile, "Executing: %v", argv)
	}

	if dryRun {
		logf(logFile, "Dry-run mode - not executing")
		os.Exit(0)
	}

	workDir := filepath.Dir(confPath)
	code, err := runner.Run(argv, cfg.Env, workDir)
	if err != nil {
		logf(logFile, "ERROR: Execution failed: %v", err)
		if code == 0 {
			os.Exit(1)
		}
		os.Exit(code)
	}

	if debug && code != 0 {
		logf(logFile, "Process exited with code: %d", code)
	}

	os.Exit(code)
}

func logf(logFile *os.File, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("[%s] [GJG] "+format, append([]interface{}{timestamp}, args...)...)
	if logFile != nil {
		_, _ = logFile.WriteString(msg + "\n")
	}
	fmt.Println(msg)
}
