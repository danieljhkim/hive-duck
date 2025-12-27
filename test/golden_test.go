package test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestGolden runs all golden output tests in test/golden/*/
func TestGolden(t *testing.T) {
	goldenDir := "golden"

	entries, err := os.ReadDir(goldenDir)
	if err != nil {
		t.Fatalf("Failed to read golden directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		testName := entry.Name()
		testDir := filepath.Join(goldenDir, testName)

		t.Run(testName, func(t *testing.T) {
			runGoldenTest(t, testDir)
		})
	}
}

func runGoldenTest(t *testing.T, testDir string) {
	inputFile := filepath.Join(testDir, "input.sql")
	expectedFile := filepath.Join(testDir, "expected.out")
	argsFile := filepath.Join(testDir, "args.txt")
	configFile := filepath.Join(testDir, "config.yaml")

	// Check required files exist
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		t.Fatalf("Missing input.sql in %s", testDir)
	}
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		t.Fatalf("Missing expected.out in %s", testDir)
	}

	// Read expected output
	expectedBytes, err := os.ReadFile(expectedFile)
	if err != nil {
		t.Fatalf("Failed to read expected.out: %v", err)
	}
	expected := strings.TrimSpace(string(expectedBytes))

	// Build command arguments
	args := []string{"-f", inputFile}

	// Add extra args from args.txt if it exists
	if argsBytes, err := os.ReadFile(argsFile); err == nil {
		extraArgs := parseArgs(string(argsBytes))
		args = append(args, extraArgs...)
	}

	// Add config if it exists
	if _, err := os.Stat(configFile); err == nil {
		args = append(args, "--config", configFile)
	}

	// Run hive-duck
	cmd := exec.Command("go", append([]string{"run", "../cmd/hive-duck"}, args...)...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		t.Fatalf("Command failed: %v\nStderr: %s", err, stderr.String())
	}

	// Compare output
	actual := strings.TrimSpace(stdout.String())

	if actual != expected {
		t.Errorf("Output mismatch:\n--- Expected ---\n%s\n--- Actual ---\n%s\n--- Diff ---\n%s",
			expected, actual, diffStrings(expected, actual))
	}
}

// parseArgs parses space/newline-separated arguments from a string
func parseArgs(s string) []string {
	var args []string
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		// Split by spaces but respect quotes (simple parsing)
		parts := strings.Fields(line)
		args = append(args, parts...)
	}
	return args
}

// diffStrings provides a simple line-by-line diff
func diffStrings(expected, actual string) string {
	expectedLines := strings.Split(expected, "\n")
	actualLines := strings.Split(actual, "\n")

	var diff strings.Builder
	maxLines := len(expectedLines)
	if len(actualLines) > maxLines {
		maxLines = len(actualLines)
	}

	for i := 0; i < maxLines; i++ {
		var expLine, actLine string
		if i < len(expectedLines) {
			expLine = expectedLines[i]
		}
		if i < len(actualLines) {
			actLine = actualLines[i]
		}

		if expLine != actLine {
			diff.WriteString("Line ")
			diff.WriteString(strings.Repeat(" ", 0))
			diff.WriteString(string(rune('0' + i + 1)))
			diff.WriteString(":\n")
			diff.WriteString("  - ")
			diff.WriteString(expLine)
			diff.WriteString("\n")
			diff.WriteString("  + ")
			diff.WriteString(actLine)
			diff.WriteString("\n")
		}
	}

	return diff.String()
}
