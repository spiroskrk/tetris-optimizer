package tetris

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunPrintsSolvedBoard(t *testing.T) {
	path := writeTestInput(t, "##..\n##..\n....\n....\n")

	output, err := captureStdout(t, func() error {
		return Run(path)
	})
	if err != nil {
		t.Fatalf("Run returned an error: %v", err)
	}
	if output != "AA\nAA\n" {
		t.Fatalf("Run output = %q; want %q", output, "AA\nAA\n")
	}
}

func TestRunWrapsReadError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing-input.txt")

	err := Run(path)
	if err == nil {
		t.Fatal("Run returned no error for a missing file")
	}
	if !strings.Contains(err.Error(), "Could not read file") {
		t.Fatalf("Run error = %q; want a read-error context", err)
	}
}

func TestRunWrapsParseError(t *testing.T) {
	path := writeTestInput(t, "not a tetromino")

	err := Run(path)
	if err == nil {
		t.Fatal("Run returned no error for malformed input")
	}
	if !strings.Contains(err.Error(), "Could not parse file") {
		t.Fatalf("Run error = %q; want a parse-error context", err)
	}
}

func writeTestInput(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "input.txt")
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write test input: %v", err)
	}
	return path
}

func captureStdout(t *testing.T, run func() error) (string, error) {
	t.Helper()

	reader, writer, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}

	originalStdout := os.Stdout
	os.Stdout = writer
	t.Cleanup(func() {
		os.Stdout = originalStdout
	})

	runErr := run()
	if err := writer.Close(); err != nil {
		t.Fatalf("close stdout writer: %v", err)
	}
	os.Stdout = originalStdout

	output, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("read captured stdout: %v", err)
	}
	if err := reader.Close(); err != nil {
		t.Fatalf("close stdout reader: %v", err)
	}

	return string(output), runErr
}
