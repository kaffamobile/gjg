//go:build windows

package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    cwd, _ := os.Getwd()
    fmt.Printf("[STUB] CWD=%s\n", cwd)
    fmt.Printf("[STUB] ARGV=%s\n", strings.Join(os.Args, "|"))
    // Print selected envs for verification
    for _, k := range []string{"MY_HOME", "DEBUG_MODE"} {
        if v, ok := os.LookupEnv(k); ok {
            fmt.Printf("[STUB] ENV %s=%s\n", k, v)
        }
    }

    // Simulate running a JAR: detect -jar and print app output
    // Capture args after the JAR file as application args
    appArgs := []string{}
    for i := 0; i < len(os.Args); i++ {
        if os.Args[i] == "-jar" && i+1 < len(os.Args) {
            // Remaining after jar path are application args
            if i+2 <= len(os.Args) {
                appArgs = append(appArgs, os.Args[i+2:]...)
            }
            break
        }
    }
    // Emit a line identical to the Hello app for smoke verification
    fmt.Printf("Hello from JAR! Args:%v\n", appArgs)
    // Optional exit code from env
    if v, ok := os.LookupEnv("STUB_EXIT"); ok {
        // attempt to parse int
        var code int
        fmt.Sscanf(v, "%d", &code)
        os.Exit(code)
    }
}
