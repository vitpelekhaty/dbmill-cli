package strings

import (
	"testing"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder("test")

	if builder.String() != "test" {
		t.FailNow()
	}
}

func TestBuilder_WriteDelimiter(t *testing.T) {
	builder := NewBuilder("test")

	builder.WriteDelimiter('/')
	builder.WriteString("test")

	if builder.String() != "test/test" {
		t.FailNow()
	}
}

func TestBuilder_WriteDelimiter2(t *testing.T) {
	var builder Builder

	builder.WriteDelimiter('/')
	builder.WriteString("test")

	if builder.String() != "test" {
		t.FailNow()
	}
}

func TestBuilder_WriteSpace(t *testing.T) {
	builder := NewBuilder("test")

	builder.WriteSpace()
	builder.WriteString("test")

	if builder.String() != "test test" {
		t.FailNow()
	}
}

func TestBuilder_WriteSpace2(t *testing.T) {
	var builder Builder

	builder.WriteSpace()
	builder.WriteString("test")

	if builder.String() != "test" {
		t.FailNow()
	}
}
