package sqlserver

import (
	"testing"
)

func TestIndex_String(t *testing.T) {
	var cases = []struct {
		index *Index
		want  string
	}{
		{
			index: &Index{
				Name:         "PK_Test1",
				Type:         "NONCLUSTERED",
				IsUnique:     true,
				IsPrimaryKey: true,
				Columns: IndexedColumns{
					"Session": &IndexedColumn{
						ID:         1,
						Name:       "Session",
						KeyOrdinal: 1,
					},
					"Division": &IndexedColumn{
						ID:         2,
						Name:       "Division",
						KeyOrdinal: 2,
					},
					"DeviceID": &IndexedColumn{
						ID:              3,
						Name:            "DeviceID",
						KeyOrdinal:      3,
						IsDescendingKey: true,
					},
				},

				owner: OwnerUserDefinedTableDataType,
			},
			want: "PRIMARY KEY [PK_Test1] ([Session], [Division], [DeviceID] DESC)",
		},
		{
			index: &Index{
				Name:         "PK_Test2",
				Type:         "NONCLUSTERED",
				IsUnique:     true,
				IsPrimaryKey: true,
				IgnoreDupKey: true,
				Columns: IndexedColumns{
					"Session": &IndexedColumn{
						ID:         1,
						Name:       "Session",
						KeyOrdinal: 1,
					},
					"Division": &IndexedColumn{
						ID:         2,
						Name:       "Division",
						KeyOrdinal: 2,
					},
				},

				owner: OwnerUserDefinedTableDataType,
			},
			want: "PRIMARY KEY [PK_Test2] ([Session], [Division]) WITH (IGNORE_DUP_KEY = ON)",
		},
		{
			index: &Index{
				Name:         "UK_Test3",
				Type:         "CLUSTERED",
				IsUnique:     true,
				IgnoreDupKey: true,
				Columns: IndexedColumns{
					"Session": &IndexedColumn{
						ID:         1,
						Name:       "Session",
						KeyOrdinal: 1,
					},
					"Division": &IndexedColumn{
						ID:         2,
						Name:       "Division",
						KeyOrdinal: 2,
					},
				},

				owner: OwnerUserDefinedTableDataType,
			},
			want: "INDEX [UK_Test3] UNIQUE CLUSTERED ([Session], [Division]) WITH (IGNORE_DUP_KEY = ON)",
		},
		{
			index: &Index{
				Name: "IX_Test4",
				Type: "NONCLUSTERED",
				Columns: IndexedColumns{
					"Session": &IndexedColumn{
						ID:         1,
						Name:       "Session",
						KeyOrdinal: 1,
					},
				},
				IncludedColumns: IndexedColumns{
					"Division": &IndexedColumn{
						ID:         2,
						Name:       "Division",
						KeyOrdinal: 2,
					},
				},

				owner: OwnerUserDefinedTableDataType,
			},
			want: "INDEX [IX_Test4] ([Session]) INCLUDE ([Division])",
		},
	}

	var have string

	for _, test := range cases {
		have = test.index.String()

		if test.want != have {
			t.Errorf("[%s].Strings() failed: have %s, want %s", test.index.Name, have, test.want)
		}
	}
}
