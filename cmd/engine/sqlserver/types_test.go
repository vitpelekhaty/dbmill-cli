package sqlserver

import (
	"database/sql"
	"testing"
)

func TestUserDefinedType_Type(t *testing.T) {
	var cases = []struct {
		udt  UserDefinedType
		want string
	}{
		{
			udt: UserDefinedType{
				Catalog:  "test",
				TypeName: "MyNVarchar",
				Schema:   "dbo",
				parentTypeName: sql.NullString{
					String: "nvarchar",
					Valid:  true,
				},
				maxLength: sql.NullString{
					String: "100",
					Valid:  true,
				},
				precision: sql.NullInt32{
					Int32: 0,
					Valid: false,
				},
				scale: sql.NullInt32{
					Int32: 0,
					Valid: false,
				},
				collation: sql.NullString{
					String: "",
					Valid:  false,
				},
				IsNullable:        false,
				IsTableType:       false,
				IsMemoryOptimized: false,
			},
			want: "[nvarchar](100)",
		},
		{
			udt: UserDefinedType{
				Catalog:  "test",
				TypeName: "MyMaxNVarchar",
				Schema:   "dbo",
				parentTypeName: sql.NullString{
					String: "nvarchar",
					Valid:  true,
				},
				maxLength: sql.NullString{
					String: "max",
					Valid:  true,
				},
				precision: sql.NullInt32{
					Int32: 0,
					Valid: false,
				},
				scale: sql.NullInt32{
					Int32: 0,
					Valid: false,
				},
				collation: sql.NullString{
					String: "",
					Valid:  false,
				},
				IsNullable:        false,
				IsTableType:       false,
				IsMemoryOptimized: false,
			},
			want: "[nvarchar](max)",
		},
		{
			udt: UserDefinedType{
				Catalog:  "test",
				TypeName: "MyFloat",
				Schema:   "dbo",
				parentTypeName: sql.NullString{
					String: "float",
					Valid:  true,
				},
				maxLength: sql.NullString{
					String: "10",
					Valid:  true,
				},
				precision: sql.NullInt32{
					Int32: 0,
					Valid: false,
				},
				scale: sql.NullInt32{
					Int32: 0,
					Valid: false,
				},
				collation: sql.NullString{
					String: "",
					Valid:  false,
				},
				IsNullable:        false,
				IsTableType:       false,
				IsMemoryOptimized: false,
			},
			want: "[float](10)",
		},
		{
			udt: UserDefinedType{
				Catalog:  "test",
				TypeName: "MyDouble",
				Schema:   "dbo",
				parentTypeName: sql.NullString{
					String: "double",
					Valid:  true,
				},
				maxLength: sql.NullString{
					String: "",
					Valid:  false,
				},
				precision: sql.NullInt32{
					Int32: 15,
					Valid: true,
				},
				scale: sql.NullInt32{
					Int32: 6,
					Valid: true,
				},
				collation: sql.NullString{
					String: "",
					Valid:  false,
				},
				IsNullable:        false,
				IsTableType:       false,
				IsMemoryOptimized: false,
			},
			want: "[double](15, 6)",
		},
		{
			udt: UserDefinedType{
				Catalog:  "test",
				TypeName: "MyAccountNumber",
				Schema:   "dbo",
				parentTypeName: sql.NullString{
					String: "double",
					Valid:  true,
				},
				maxLength: sql.NullString{
					String: "",
					Valid:  false,
				},
				precision: sql.NullInt32{
					Int32: 25,
					Valid: true,
				},
				scale: sql.NullInt32{
					Int32: 0,
					Valid: true,
				},
				collation: sql.NullString{
					String: "",
					Valid:  false,
				},
				IsNullable:        false,
				IsTableType:       false,
				IsMemoryOptimized: false,
			},
			want: "[double](25)",
		},
		{
			udt: UserDefinedType{
				Catalog:  "test",
				TypeName: "MyBit",
				Schema:   "dbo",
				parentTypeName: sql.NullString{
					String: "bit",
					Valid:  true,
				},
				maxLength: sql.NullString{
					String: "",
					Valid:  false,
				},
				precision: sql.NullInt32{
					Int32: 0,
					Valid: false,
				},
				scale: sql.NullInt32{
					Int32: 0,
					Valid: false,
				},
				collation: sql.NullString{
					String: "",
					Valid:  false,
				},
				IsNullable:        false,
				IsTableType:       false,
				IsMemoryOptimized: false,
			},
			want: "[bit]",
		},
	}

	var have string

	for _, test := range cases {
		have = test.udt.ParentType()

		if have != test.want {
			t.Errorf("TestUserDefinedType_Type failed for %s: have %s, want %s",
				test.udt.SchemaAndName(true), have, test.want)
		}
	}
}
