package command

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/kritama/memovee-cli/internal/version"
)

func TestVersionHumanOutput(t *testing.T) {
	setVersion(t, "1.2.3", "abc123", "2026-09-01T00:00:00Z")

	exitCode, stdout, stderr := runCommand("version")

	if exitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want %d", exitCode, ExitSuccess)
	}
	if stdout != "memovee 1.2.3 (revision abc123, built 2026-09-01T00:00:00Z)\n" {
		t.Fatalf("stdout = %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}

func TestVersionJSONOutput(t *testing.T) {
	setVersion(t, "1.2.3", "abc123", "2026-09-01T00:00:00Z")

	exitCode, stdout, stderr := runCommand("--json", "version")

	if exitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want %d", exitCode, ExitSuccess)
	}
	want := "{\"schema_version\":\"1\",\"ok\":true," +
		"\"result\":{\"command\":\"version\",\"version\":\"1.2.3\"," +
		"\"revision\":\"abc123\",\"build_time\":\"2026-09-01T00:00:00Z\"}}\n"
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}

func TestGlobalOptionsAreAcceptedBeforeCommand(t *testing.T) {
	exitCode, _, stderr := runCommand(
		"--no-color",
		"--non-interactive",
		"--yes",
		"--config",
		"config.json",
		"version",
	)

	if exitCode != ExitSuccess {
		t.Fatalf("exit code = %d, want %d; stderr = %q", exitCode, ExitSuccess, stderr)
	}
}

func TestHelp(t *testing.T) {
	for _, args := range [][]string{{"help"}, {"--help"}, {"-h"}} {
		exitCode, stdout, stderr := runCommand(args...)

		if exitCode != ExitSuccess {
			t.Errorf("Run(%q) exit code = %d, want %d", args, exitCode, ExitSuccess)
		}
		if stdout != usage {
			t.Errorf("Run(%q) stdout = %q, want usage", args, stdout)
		}
		if stderr != "" {
			t.Errorf("Run(%q) stderr = %q, want empty", args, stderr)
		}
	}
}

func TestMissingCommandHumanError(t *testing.T) {
	exitCode, stdout, stderr := runCommand()

	if exitCode != ExitUsage {
		t.Fatalf("exit code = %d, want %d", exitCode, ExitUsage)
	}
	if stdout != "" {
		t.Fatalf("stdout = %q, want empty", stdout)
	}
	want := "error (usage): a command is required\n" +
		"next: run `memovee help` for usage\n"
	if stderr != want {
		t.Fatalf("stderr = %q, want %q", stderr, want)
	}
}

func TestUnknownCommandJSONError(t *testing.T) {
	exitCode, stdout, stderr := runCommand("--json", "unknown")

	if exitCode != ExitUsage {
		t.Fatalf("exit code = %d, want %d", exitCode, ExitUsage)
	}
	want := "{\"schema_version\":\"1\",\"ok\":false," +
		"\"error\":{\"category\":\"usage\",\"exit_code\":2," +
		"\"message\":\"unknown command \\\"unknown\\\"\"," +
		"\"next\":\"run `memovee help` for usage\"}}\n"
	if stdout != want {
		t.Fatalf("stdout = %q, want %q", stdout, want)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}

func TestVersionRejectsArguments(t *testing.T) {
	exitCode, _, stderr := runCommand("version", "extra")

	if exitCode != ExitUsage {
		t.Fatalf("exit code = %d, want %d", exitCode, ExitUsage)
	}
	if !strings.Contains(stderr, "command \"version\" does not accept arguments") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestInvalidGlobalOption(t *testing.T) {
	exitCode, _, stderr := runCommand("--unknown", "version")

	if exitCode != ExitUsage {
		t.Fatalf("exit code = %d, want %d", exitCode, ExitUsage)
	}
	if !strings.Contains(stderr, "flag provided but not defined") {
		t.Fatalf("stderr = %q", stderr)
	}
}

func TestJSONFlagAppliesToLaterOptionErrors(t *testing.T) {
	exitCode, stdout, stderr := runCommand("--json", "--unknown", "version")

	if exitCode != ExitUsage {
		t.Fatalf("exit code = %d, want %d", exitCode, ExitUsage)
	}
	if !strings.Contains(stdout, `"category":"usage"`) {
		t.Fatalf("stdout = %q", stdout)
	}
	if stderr != "" {
		t.Fatalf("stderr = %q, want empty", stderr)
	}
}

func TestOutputFailureUsesInternalExitCode(t *testing.T) {
	exitCode := Run([]string{"version"}, IO{
		Stdin:  strings.NewReader(""),
		Stdout: errorWriter{},
		Stderr: errorWriter{},
	})

	if exitCode != ExitInternal {
		t.Fatalf("exit code = %d, want %d", exitCode, ExitInternal)
	}
}

func runCommand(args ...string) (int, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	exitCode := Run(args, IO{
		Stdin:  strings.NewReader(""),
		Stdout: &stdout,
		Stderr: &stderr,
	})
	return exitCode, stdout.String(), stderr.String()
}

func setVersion(t *testing.T, value, revision, buildTime string) {
	t.Helper()
	previousVersion := version.Version
	previousRevision := version.Revision
	previousBuildTime := version.BuildTime
	t.Cleanup(func() {
		version.Version = previousVersion
		version.Revision = previousRevision
		version.BuildTime = previousBuildTime
	})

	version.Version = value
	version.Revision = revision
	version.BuildTime = buildTime
}

type errorWriter struct{}

func (errorWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}
