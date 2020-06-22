package sqlserver

import (
	"database/sql"
	"testing"
)

func TestColumn_String(t *testing.T) {
	var cases = []struct {
		column  Column
		options []ColumnOption
		want    string
	}{
		{
			column: Column{
				ID:       1,
				Name:     "key",
				TypeName: "int",
			},
			options: nil,
			want:    "[key] [int] NOT NULL",
		},
		{
			column: Column{
				ID:       2,
				Name:     "value",
				TypeName: "nvarchar",
				maxLength: sql.NullString{
					String: "2048",
					Valid:  true,
				},
				collation: sql.NullString{
					String: "SQL_Latin1_General_CP1_CI_AS",
					Valid:  true,
				},
				IsNullable:   true,
				IsANSIPadded: true,
			},
			options: nil,
			want:    "[value] [nvarchar](2048)",
		},
	}

	var have string

	for _, test := range cases {
		if len(test.options) > 0 {
			test.column.SetOptions(test.options...)
		}

		have = test.column.String()

		if have != test.want {
			t.Errorf("TestColumn_String for field %s failed: have %s, want %s", test.column.Name, have,
				test.want)
		}
	}
}
