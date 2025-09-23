//go:build windows

package main

import (
    "fmt"
    "os"
    "path/filepath"

    "gjg/internal/args"
    "gjg/internal/config"
    "gjg/internal/quote"
    "gjg/internal/resolve"
    "gjg/internal/runner"
)

const (
    exitConfigError = 201
    exitJavaMissing = 202
    exitJarMissing  = 203
)

func main() {
    // Parse special flags from CLI
    debug, dryRun, forwardArgs := args.ExtractSpecial(os.Args[1:])

    // Load configuration
    cfg, confPath, err := config.LoadFromExe()
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(exitConfigError)
    }

    // Resolve paths relative to exe dir
    exePath, _ := os.Executable()
    exeDir := filepath.Dir(exePath)

    javaPath, err := resolve.ResolveJava(cfg.JavaDir, cfg.JavaExecutable, exeDir)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(exitJavaMissing)
    }

    jarPath, err := resolve.ResolveJar(cfg.JarFile, exeDir)
    if err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
        os.Exit(exitJarMissing)
    }

    // Build final argument vector
    jvmTokens := args.Tokenize(cfg.JVMArgs)
    appTokens := args.Tokenize(cfg.AppArgs)

    argv := make([]string, 0, 4+len(jvmTokens)+len(appTokens)+len(forwardArgs))
    argv = append(argv, javaPath)
    argv = append(argv, jvmTokens...)
    argv = append(argv, "-jar", jarPath)
    argv = append(argv, appTokens...)
    argv = append(argv, forwardArgs...)

    // Prepare env
    env := config.MergeEnv(os.Environ(), cfg.Env)

    // Debug and Dry-run
    if debug {
        fmt.Fprintf(os.Stdout, "[GJG] Configuration loaded from: %s\n", confPath)
        fmt.Fprintf(os.Stdout, "[GJG] Java executable: %s\n", javaPath)
        fmt.Fprintf(os.Stdout, "[GJG] JAR file: %s\n", jarPath)
        if len(cfg.Env) > 0 {
            fmt.Fprintf(os.Stdout, "[GJG] Environment variables: %s\n", config.EnvSummary(cfg.Env))
        } else {
            fmt.Fprintf(os.Stdout, "[GJG] Environment variables: (none)\n")
        }
        fmt.Fprintf(os.Stdout, "[GJG] Executing: %s\n", quote.JoinWindows(argv))
    }

    if dryRun {
        // Do not execute
        return
    }

    // Execute
    code, err := runner.Run(argv, env, exeDir)
    if err != nil {
        // Best-effort: print error; if exit code was obtained, use it, else 1
        fmt.Fprintln(os.Stderr, err.Error())
        if code == 0 {
            os.Exit(1)
        }
        os.Exit(code)
    }
    os.Exit(code)
}

