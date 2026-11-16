package cli

import (
	"os"
	"testing"
)

func TestIsTerminal(t *testing.T) {
	t.Run("says no for /dev/null and for a pipe (FR-TUI-01)", func(t *testing.T) {
		null, err := os.Open(os.DevNull)
		if err != nil {
			t.Fatal(err)
		}
		defer null.Close()
		r, w, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		defer r.Close()
		defer w.Close()
		for _, f := range []*os.File{null, r, w} {
			if isTerminal(f) {
				t.Errorf("%s is not a terminal", f.Name())
			}
		}
	})
}
