package testfixture

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("builds every template of the spec", func(t *testing.T) {
		dir := New(t)
		for _, id := range IDs {
			if _, err := os.Stat(filepath.Join(dir, id)); err != nil {
				t.Errorf("missing %q: %v", id, err)
			}
		}
		info, err := os.Stat(filepath.Join(dir, "Empty Bash"))
		if err != nil {
			t.Fatal(err)
		}
		if got := info.Mode().Perm(); got != 0o755 {
			t.Errorf("Empty Bash has %o, want 755", got)
		}
	})
}

func TestSnapshot(t *testing.T) {
	t.Run("changes when a file is added", func(t *testing.T) {
		dir := New(t)
		before := Snapshot(t, dir)
		Write(t, dir, "novo", "x", 0o644)
		if Snapshot(t, dir) == before {
			t.Error("the snapshot did not change")
		}
	})
}

func TestFakeProgram(t *testing.T) {
	t.Run("replies in order and records the calls", func(t *testing.T) {
		fake := FakeProgram(t, "prog",
			Reply{Stdout: "um\n"},
			Reply{Stdout: "dois", Stderr: "erro\n", Code: 3},
		)
		first := exec.Command("prog", "a b", "--x=y")
		first.Stdin = strings.NewReader("entrada\n")
		out, err := first.Output()
		if err != nil || string(out) != "um\n" {
			t.Fatalf("first call: %q, %v", out, err)
		}
		second := exec.Command("prog")
		var stderr strings.Builder
		second.Stderr = &stderr
		out, err = second.Output()
		var exit *exec.ExitError
		if !errors.As(err, &exit) || exit.ExitCode() != 3 || string(out) != "dois" || stderr.String() != "erro\n" {
			t.Fatalf("second call: %q, %q, %v", out, stderr.String(), err)
		}
		calls := fake.Calls(t)
		if len(calls) != 2 {
			t.Fatalf("got %d calls, want 2", len(calls))
		}
		if !slices.Equal(calls[0].Args, []string{"a b", "--x=y"}) || calls[0].Stdin != "entrada\n" {
			t.Errorf("first call recorded as %+v", calls[0])
		}
		if len(calls[1].Args) != 0 {
			t.Errorf("second call recorded args %q", calls[1].Args)
		}
	})
}
