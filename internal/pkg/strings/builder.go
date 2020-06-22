package strings

import (
	"strings"
)

// Builder расширенная реализация StringBuilder
type Builder struct {
	strings.Builder
}

// NewBuilder конструктор
func NewBuilder(value string) *Builder {
	var builder Builder
	builder.WriteString(value)

	return &builder
}

// WriteDelimiter вставляет символ delimiter в качестве разделителя, если в builder уже имеется какой-то текст
func (builder *Builder) WriteDelimiter(delimiter rune) {
	if builder.Len() > 0 {
		builder.WriteRune(delimiter)
	}
}

// WriteSpace вставляет пробел в качестве разделителя, если в builder уже имеется какой-то текст
func (builder *Builder) WriteSpace() {
	builder.WriteDelimiter(' ')
}
