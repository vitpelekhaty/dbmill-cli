package engine

import (
	"testing"
)

var RDBMSCases = []struct {
	connection string
	rdbms      RDBMSType
}{
	{
		connection: "sqlserver://localhost",
		rdbms:      RDBMSSQLServer,
	},
	{
		connection: "postgres://localhost",
		rdbms:      RDBMSUnknown,
	},
	{
		connection: "localhost",
		rdbms:      RDBMSUnknown,
	},
}

func TestRDBMS(t *testing.T) {
	for _, test := range RDBMSCases {
		dbt, _ := RDBMS(test.connection)

		if dbt != test.rdbms {
			t.Errorf(`RDBMS("%s") failed!`, test.connection)
		}
	}
}
