package dir

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	items, err := parse([]byte(defaultData))

	if err != nil {
		t.Fatal(err)
	}

	if len(items) == 0 {
		t.FailNow()
	}
}

func TestNewStructure(t *testing.T) {
	s, err := NewStructure(strings.NewReader(defaultData))

	if err != nil {
		t.Fatal(err)
	}

	if len(s.Items) == 0 {
		t.FailNow()
	}
}
