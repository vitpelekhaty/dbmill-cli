package output

import (
	"testing"
)

func TestInit(t *testing.T) {
	if DefaultScriptsFolderOutput == nil {
		t.FailNow()
	}

	if len(DefaultScriptsFolderOutput.rules) == 0 {
		t.FailNow()
	}
}
