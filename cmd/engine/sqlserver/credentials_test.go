package sqlserver

import (
	"testing"
)

type credentials struct {
	username string
	password string
}

var cases = []struct {
	connection string
	want       credentials
}{
	{
		connection: "sqlserver://username:password@localhost",
		want: credentials{
			username: "username",
			password: "password",
		},
	},
	{
		connection: "sqlserver://username@localhost",
		want: credentials{
			username: "username",
			password: "",
		},
	},
	{
		connection: "sqlserver://localhost",
		want: credentials{
			username: "",
			password: "",
		},
	},
	{
		connection: "sqlserver://localhost/instance",
		want: credentials{
			username: "",
			password: "",
		},
	},
	{
		connection: "sqlserver://localhost:3112",
		want: credentials{
			username: "",
			password: "",
		},
	},
	{
		connection: "sqlserver://localhost:3112/instance",
		want: credentials{
			username: "",
			password: "",
		},
	},
	{
		connection: "sqlserver://localhost:3112/instance?user id=username",
		want: credentials{
			username: "username",
			password: "",
		},
	},
	{
		connection: "sqlserver://localhost:3112/instance?user id=username&password=password",
		want: credentials{
			username: "username",
			password: "password",
		},
	},
	{
		connection: "sqlserver://username:password@localhost:3112/instance?user id=username2&password=password2",
		want: credentials{
			username: "username",
			password: "password",
		},
	},
	{
		connection: "sqlserver://username:my%7Bpass@localhost",
		want: credentials{
			username: "username",
			password: "my{pass",
		},
	},
}

func TestCredentials(t *testing.T) {
	var done bool

	for _, test := range cases {
		username, password, err := Credentials(test.connection)

		if err != nil {
			t.Fatal(err)
		}

		done = username == test.want.username && password == test.want.password

		if !done {
			t.Errorf(`credentials("%s") failed: have username - %s, want - %s; have password - %s, want - %s`,
				test.connection, username, test.want.username, password, test.want.password)
		}
	}
}
