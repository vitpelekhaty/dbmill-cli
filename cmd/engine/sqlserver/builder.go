package sqlserver

import (
	"strings"
)

// Builder расширенная реализация StringBuilder
type Builder struct {
	strings.Builder
}

// InsertDelimiter вставляет символ delimiter в качестве разделителя, если в builder уже имеется какой-то текст
func (builder *Builder) InsertDelimiter(delimiter rune) {
	if builder.Len() > 0 {
		builder.WriteRune(delimiter)
	}
}

// InsertDelimiter вставляет пробел в качестве разделителя, если в builder уже имеется какой-то текст
func (builder *Builder) InsertSpace() {
	builder.InsertDelimiter(' ')
}
