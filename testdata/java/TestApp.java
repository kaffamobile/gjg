package testdata;

import java.io.File;
import java.util.Arrays;
import java.util.Map;

/**
 * Test application for GJG launcher testing.
 * Compatible with Java 8+ and provides comprehensive testing capabilities.
 */
public class TestApp {
    public static void main(String[] args) {
        System.out.println("=== GJG Test Application ===");

        // Print working directory
        System.out.println("Working Directory: " + System.getProperty("user.dir"));

        // Print JVM arguments that are commonly tested
        printSystemProperty("java.version");
        printSystemProperty("java.library.path");
        printSystemProperty("user.name");

        // Print command line arguments
        System.out.println("Command Line Args: " + Arrays.toString(args));

        // Print specific environment variables that tests set
        printEnvironmentVariable("MY_HOME");
        printEnvironmentVariable("DEBUG_MODE");
        printEnvironmentVariable("TEST_ENV");
        printEnvironmentVariable("WORKSPACE");

        // Print all environment variables starting with GJG_ for testing
        System.out.println("GJG Environment Variables:");
        for (Map.Entry<String, String> entry : System.getenv().entrySet()) {
            if (entry.getKey().startsWith("GJG_")) {
                System.out.println("  " + entry.getKey() + "=" + entry.getValue());
            }
        }

        // Handle special test modes
        for (String arg : args) {
            if ("--test-exit-code".equals(arg)) {
                System.out.println("Test mode: custom exit code");
                String exitCode = System.getProperty("gjg.test.exitcode", "42");
                System.exit(Integer.parseInt(exitCode));
            } else if ("--test-error".equals(arg)) {
                System.out.println("Test mode: error output");
                System.err.println("This is test error output");
            } else if ("--test-long-output".equals(arg)) {
                System.out.println("Test mode: long output");
                for (int i = 1; i <= 5; i++) {
                    System.out.println("Line " + i + ": This is test output line number " + i);
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