package testdata;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.time.LocalDateTime;
import java.time.format.DateTimeFormatter;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

/**
 * Modern test application for GJG launcher testing.
 * Uses modern Java features since we compile at test time.
 */
public class TestApp {
    public static void main(String[] args) {
        var timestamp = LocalDateTime.now().format(DateTimeFormatter.ISO_LOCAL_TIME);
        System.out.println("=== GJG Test Application (" + timestamp + ") ===");

        // Use modern Path API
        Path workingDir = Paths.get(System.getProperty("user.dir"));
        System.out.println("Working Directory: " + workingDir.toAbsolutePath());

        // Print JVM info using modern string methods
        printSystemProperty("java.version");
        printSystemProperty("java.library.path");
        printSystemProperty("user.name");

        // Use String.join for cleaner output
        var argsList = List.of(args);
        System.out.println("Command Line Args: [" + String.join(", ", argsList) + "]");

        // Print specific environment variables that tests set
        printEnvironmentVariable("MY_HOME");
        printEnvironmentVariable("DEBUG_MODE");
        printEnvironmentVariable("TEST_ENV");
        printEnvironmentVariable("WORKSPACE");

        // Use streams for GJG environment variables
        var gjgEnvVars = System.getenv().entrySet().stream()
            .filter(entry -> entry.getKey().startsWith("GJG_"))
            .collect(Collectors.toMap(Map.Entry::getKey, Map.Entry::getValue));

        if (!gjgEnvVars.isEmpty()) {
            System.out.println("GJG Environment Variables:");
            gjgEnvVars.forEach((key, value) ->
                System.out.println("  " + key + "=" + value));
        }

        // Handle special test modes with enhanced switch
        for (var arg : args) {
            switch (arg) {
                case "--test-exit-code" -> {
                    System.out.println("Test mode: custom exit code");
                    var exitCode = System.getProperty("gjg.test.exitcode", "42");
                    System.exit(Integer.parseInt(exitCode));
                }
                case "--test-error" -> {
                    System.out.println("Test mode: error output");
                    System.err.println("This is test error output");
                }
                case "--test-long-output" -> {
                    System.out.println("Test mode: long output");
                    for (int i = 1; i <= 5; i++) {
                        System.out.println("Line %d: This is test output line number %d".formatted(i, i));
                    }
                }
                case "--test-java-version" -> {
                    System.out.println("Java Runtime Version: " + Runtime.version());
                    System.out.println("Available Processors: " + Runtime.getRuntime().availableProcessors());
                }
            }
        }

        System.out.println("=== Test Application Complete ===");
    }

    private static void printSystemProperty(String key) {
        String value = System.getProperty(key);
        if (value != null) {
            System.out.println("System Property " + key + "=" + value);
        }
    }

    private static void printEnvironmentVariable(String key) {
        String value = System.getenv(key);
        if (value != null) {
            System.out.println("Environment Variable " + key + "=" + value);
        }
    }
}