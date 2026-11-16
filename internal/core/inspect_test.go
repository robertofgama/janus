package core

import (
	"path/filepath"
	"strings"
	"testing"

	"gitlab.com/rg-dev/janus/internal/testfixture"
)

func TestInspect(t *testing.T) {
	dir := testfixture.New(t)

	t.Run("classifies the templates for the preview (FR-CORE-23)", func(t *testing.T) {
		html := filepath.Join(dir, "info.dev/react; script; simple.html")
		cases := []struct {
			path string
			kind PreviewKind
			size int64
		}{
			{filepath.Join(dir, "Empty File"), PreviewEmpty, 0},
			{filepath.Join(dir, "3d/fatiar.3mf"), PreviewBinary, int64(len(testfixture.Binary))},
			{html, PreviewText, int64(len(testfixture.Read(t, html)))},
		}
		for _, c := range cases {
			kind, size, err := Inspect(c.path)
			if err != nil || kind != c.kind || size != c.size {
				t.Errorf("%s: got %d, %d, %v; want %d, %d", c.path, kind, size, err, c.kind, c.size)
			}
		}
	})

	t.Run("looks for the null byte only in the first 8000 bytes (FR-CORE-23)", func(t *testing.T) {
		late := testfixture.Write(t, dir, "tarde", strings.Repeat("a", 8000)+"\x00", 0o644)
		early := testfixture.Write(t, dir, "cedo", strings.Repeat("a", 7999)+"\x00", 0o644)
		if kind, _, _ := Inspect(late); kind != PreviewText {
			t.Errorf("a null byte at 8001 gave %d, want text", kind)
		}
		if kind, _, _ := Inspect(early); kind != PreviewBinary {
			t.Errorf("a null byte at 8000 gave %d, want binary", kind)
		}
	})
}
