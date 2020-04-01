package dir

import "testing"

func TestDefault(t *testing.T) {
	if Default == nil {
		t.FailNow()
	}

	if len(Default.Items) == 0 {
		t.FailNow()
	}
}
