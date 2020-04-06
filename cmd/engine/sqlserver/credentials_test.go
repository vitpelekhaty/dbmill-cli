package sqlserver

import (
	"net/url"
	"testing"
)

type credentials struct {
	username string
	password string
}

var cases4Credentials = []struct {
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

	for _, test := range cases4Credentials {
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

var cases4SetCredentials1 = []struct {
	connection string
	user       credentials
	want       string
}{
	{
		connection: "sqlserver://username:password@localhost",
		user: credentials{
			username: "username1",
			password: "password1",
		},
		want: "sqlserver://username1:password1@localhost",
	},
	{
		connection: "sqlserver://username@localhost",
		user: credentials{
			username: "username1",
			password: "password1",
		},
		want: "sqlserver://username1:password1@localhost",
	},
	{
		connection: "sqlserver://localhost",
		user: credentials{
			username: "username1",
			password: "password1",
		},
		want: "sqlserver://username1:password1@localhost",
	},
	{
		connection: "sqlserver://username:password@localhost",
		user: credentials{
			username: "username1",
			password: "",
		},
		want: "sqlserver://username1@localhost",
	},
	{
		connection: "sqlserver://username@localhost",
		user: credentials{
			username: "",
			password: "",
		},
		want: "sqlserver://localhost",
	},
	{
		connection: "sqlserver://localhost?password=password",
		user: credentials{
			username: "username1",
			password: "password1",
		},
		want: "sqlserver://username1:password1@localhost?password=password",
	},
}

func TestSetCredentials1(t *testing.T) {
	var done bool

	for _, test := range cases4SetCredentials1 {
		connection, err := SetCredentials(test.connection, test.user.username, test.user.password)

		if err != nil {
			t.Fatal(err)
		}

		done = connection == test.want

		if !done {
			t.Errorf(`SetCredentials("%s", "%s", "%s") failed: have connection - %s, want - %s`,
				test.connection, test.user.username, test.user.password, connection, test.want)
		}
	}
}

var cases4SetCredentials2 = []struct {
	connection string
	user       credentials
}{

	{
		connection: "sqlserver://localhost?user id=username&password=password",
		user: credentials{
			username: "username1",
			password: "password2",
		},
	},
	{
		connection: "sqlserver://localhost?user id=username&password=password",
		user: credentials{
			username: "username1",
			password: "",
		},
	},
	{
		connection: "sqlserver://localhost?user id=username&password=password",
		user: credentials{
			username: "",
			password: "",
		},
	},
	{
		connection: "sqlserver://localhost?user id=username",
		user: credentials{
			username: "username1",
			password: "password1",
		},
	},
}

func TestSetCredentials2(t *testing.T) {
	var done bool

	for _, test := range cases4SetCredentials2 {
		connection, err := SetCredentials(test.connection, test.user.username, test.user.password)

		if err != nil {
			t.Fatal(err)
		}

		u, err := url.Parse(connection)

		if err != nil {
			t.Fatal(err)
		}

		q, err := url.ParseQuery(u.RawQuery)

		if err != nil {
			t.Fatal(err)
		}

		username := q.Get("user id")
		password := q.Get("password")

		done = username == test.user.username && password == test.user.password

		if !done {
			t.Errorf(`SetCredentials("%s", "%s", "%s") failed: have username - %s, want - %s; have password - %s, want - %s`,
				test.connection, test.user.username, test.user.password, username, test.user.username, password, test.user.password)
		}
	}
}
