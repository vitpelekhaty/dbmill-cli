package sqlserver

import (
	"database/sql"
	"testing"
)

func TestColumn_String(t *testing.T) {
	var cases = []struct {
		column Column
		want   string
	}{
		{
			column: Column{
				ID:                 1,
				Name:               "key",
				TypeName:           "int",
				IsNullable:         false,
				IsANSIPadded:       false,
				IsRowGUIDCol:       false,
				IsIdentity:         false,
				isComputed:         false,
				IsFileStream:       false,
				IsReplicated:       false,
				IsNonSQLSubscribed: false,
				IsMergePublished:   false,
				IsDTSReplicated:    false,
				IsXMLDocument:      false,
				IsSparse:           false,
				IsColumnSet:        false,
				IsHidden:           false,
				IsMasked:           false,
			},
			want: "[key] [int] NOT NULL",
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
				IsNullable:         true,
				IsANSIPadded:       true,
				IsRowGUIDCol:       false,
				IsIdentity:         false,
				isComputed:         false,
				IsFileStream:       false,
				IsReplicated:       false,
				IsNonSQLSubscribed: false,
				IsMergePublished:   false,
				IsDTSReplicated:    false,
				IsXMLDocument:      false,
				IsSparse:           false,
				IsColumnSet:        false,
				IsHidden:           false,
				IsMasked:           false,
			},
			want: "[value] [nvarchar](2048)",
		},
	}

	var have string

	for _, test := range cases {
		have = test.column.String()

		if have != test.want {
			t.Errorf("TestColumn_String for field %s failed: have %s, want %s", test.column.Name, have,
				test.want)
		}
	}
}
