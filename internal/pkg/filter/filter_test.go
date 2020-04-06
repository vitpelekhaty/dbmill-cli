package filter

import (
	"testing"
)

var testArray = []string{"(dbo|tmp)\\.Proc1", "dbo\\.Proc2", "^dbo\\..{1,}Asute.{1,}"}

var testCases = []struct {
	have string
	want bool
}{
	{
		have: "tmp.Proc1",
		want: true,
	},
	{
		have: "dbo.Get_Asute_Devices",
		want: true,
	},
	{
		have: "dbo.Proc2",
		want: true,
	},
	{
		have: "grf.Get_Devices",
		want: false,
	},
}

func TestNew(t *testing.T) {
	if _, err := New(testArray); err != nil {
		t.Fatal(err)
	}
}

func TestNewWithEmptyErray(t *testing.T) {
	if _, err := New(nil); err != nil {
		t.Fatal(err)
	}
}

func TestMatch(t *testing.T) {
	f, err := New(testArray)

	if err != nil {
		t.Fatal(err)
	}

	for _, test := range testCases {
		err = f.Match(test.have)

		if err != nil && err != ErrorNotMatched {
			t.Fatal(err)
		}

		if test.want && (err == ErrorNotMatched) {
			t.Errorf("match failed on %s, must %v", test.have, test.want)
		}

		if !test.want && (err == nil) {
			t.Errorf("match failed on %s, must %v", test.have, test.want)
		}
	}
}

func TestMatchWithEmptyArray(t *testing.T) {
	f, err := New(nil)

	if err != nil {
		t.Fatal()
	}

	var done bool

	for _, test := range testCases {
		err = f.Match(test.have)

		if err != nil && err != ErrorNotMatched {
			t.Fatal(err)
		}

		done = done || err == nil
	}

	if !done {
		t.Fail()
	}
}
