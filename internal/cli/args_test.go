package cli

import (
	"errors"
	"maps"
	"slices"
	"testing"

	"gitlab.com/rg-dev/janus/internal/failure"
)

var testOptions = map[string]bool{"dir": true, "name": true, "json": false}

func TestParseArgs(t *testing.T) {
	t.Run("reads options in both forms and in any position (FR-CLI-02)", func(t *testing.T) {
		cases := []struct {
			args    []string
			options map[string]string
			rest    []string
		}{
			{[]string{"id", "--dir", "x"}, map[string]string{"dir": "x"}, []string{"id"}},
			{[]string{"--dir=x", "id"}, map[string]string{"dir": "x"}, []string{"id"}},
			{[]string{"id", "--name=a=b"}, map[string]string{"name": "a=b"}, []string{"id"}},
			{[]string{"--json"}, map[string]string{"json": ""}, nil},
			{[]string{"--name", "y", "--", "-x.md"}, map[string]string{"name": "y"}, []string{"-x.md"}},
			{[]string{"--name="}, map[string]string{"name": ""}, nil},
			{[]string{"--", "--dir"}, map[string]string{}, []string{"--dir"}},
		}
		for _, c := range cases {
			p, err := parseArgs("new", c.args, testOptions)
			if err != nil || !maps.Equal(p.options, c.options) || !slices.Equal(p.args, c.rest) {
				t.Errorf("parseArgs(%q) = %v %q, %v", c.args, p.options, p.args, err)
			}
		}
	})

	t.Run("rejects what the contract does not accept (FR-CLI-01)", func(t *testing.T) {
		cases := []struct {
			args    []string
			problem failure.UsageProblem
			arg     string
		}{
			{[]string{"--x"}, failure.UnknownOption, "--x"},
			{[]string{"-name", "x"}, failure.UnknownOption, "-name"},
			{[]string{"-"}, failure.UnknownOption, "-"},
			{[]string{"--json=1"}, failure.UnknownOption, "--json=1"},
			{[]string{"--dir", "a", "--dir", "b"}, failure.RepeatedOption, "--dir"},
			{[]string{"--dir=a", "--dir", "b"}, failure.RepeatedOption, "--dir"},
			{[]string{"id", "--name"}, failure.MissingValue, "--name"},
		}
		for _, c := range cases {
			_, err := parseArgs("new", c.args, testOptions)
			var e *failure.UsageError
			if !errors.As(err, &e) || e.Problem != c.problem || e.Arg != c.arg || e.Command != "new" {
				t.Errorf("parseArgs(%q) = %v (%+v), want problem %d on %q", c.args, err, e, c.problem, c.arg)
			}
		}
	})
}
