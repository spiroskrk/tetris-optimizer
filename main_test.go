package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestCommandLine(t *testing.T) {
	goCache := filepath.Join(t.TempDir(), "go-cache")
	validInput := writeCLIInput(t, "valid.txt", "##..\n##..\n....\n....\n")
	badInput := writeCLIInput(t, "bad.txt", "not a tetromino")
	missingInput := filepath.Join(t.TempDir(), "missing.txt")

	tests := []struct {
		name           string
		args           []string
		wantStdout     string
		stdoutContains string
	}{
		{
			name:       "valid input",
			args:       []string{validInput},
			wantStdout: "AA\nAA\n",
		},
		{
			name:       "missing argument",
			wantStdout: "Please insert only one Argument.\n",
		},
		{
			name:       "too many arguments",
			args:       []string{validInput, validInput},
			wantStdout: "Please insert only one Argument.\n",
		},
		{
			name:           "malformed input",
			args:           []string{badInput},
			stdoutContains: "Could not parse file",
		},
		{
			name:           "missing input file",
			args:           []string{missingInput},
			stdoutContains: "Could not read file",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdout, stderr, err := runGoCommand(t, goCache, test.args...)
			if err != nil {
				t.Fatalf("go run . exited with an error: %v\nstderr: %s", err, stderr)
			}
			if stderr != "" {
				t.Errorf("stderr = %q; want empty stderr", stderr)
			}
			if test.wantStdout != "" && stdout != test.wantStdout {
				t.Errorf("stdout = %q; want %q", stdout, test.wantStdout)
			}
			if test.stdoutContains != "" && !strings.Contains(stdout, test.stdoutContains) {
				t.Errorf("stdout = %q; want it to contain %q", stdout, test.stdoutContains)
			}
		})
	}
}

func writeCLIInput(t *testing.T, name, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
		t.Fatalf("write CLI input: %v", err)
	}
	return path
}

func runGoCommand(t *testing.T, goCache string, args ...string) (string, string, error) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	commandArgs := append([]string{"run", "."}, args...)
	command := exec.CommandContext(ctx, "go", commandArgs...)
	command.Env = append(os.Environ(), "GOCACHE="+goCache)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err := command.Run()
	if ctx.Err() != nil {
		t.Fatalf("go run . timed out: %v", ctx.Err())
	}

	return stdout.String(), stderr.String(), err
}
