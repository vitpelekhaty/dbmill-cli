package output

import (
	"strings"
	"testing"
)

func TestParseScriptsFolderOutput(t *testing.T) {
	rules, err := parse([]byte(defaultScriptsFolderOutput))

	if err != nil {
		t.Fatal(err)
	}

	if len(rules) == 0 {
		t.FailNow()
	}
}

func TestNewScriptsFolderOutput(t *testing.T) {
	s, err := NewScriptsFolderOutput(strings.NewReader(defaultScriptsFolderOutput))

	if err != nil {
		t.Fatal(err)
	}

	if len(s.rules) == 0 {
		t.FailNow()
	}
}
