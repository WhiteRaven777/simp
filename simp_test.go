package simp

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func TestDriverName(t *testing.T) {
	if MySQL.String() != "mysql" {
		t.Error("MySQL is wrong.")
	}
	if PostgreSQL.String() != "postgres" {
		t.Error("PostgreSQL is wrong.")
	}
}

func TestDSN(t *testing.T) {
	dc := DsnConf{}
	dc = DsnConf{
		UserName: "root",
		Password: "pass",
		DbName:   "database",
		Params: map[string]string{
			"parseTime":            "true",
			"charset":              "utf8mb4",
			"autocommit":           "false",
			"clientFoundRows":      "true",
			"allowNativePasswords": "True",
		},
	}
	if dsn, err := dc.DSN(MySQL); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(dsn)
	}
	// ---
	dc = DsnConf{}
	dc = DsnConf{
		UserName: "postgres",
		Password: "pass",
		DbName:   "database",
		Params: map[string]string{
			"sslmode": "disable",
		},
	}
	if dsn, err := dc.DSN(PostgreSQL); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(dsn)
	}
	// ---
	dc = DsnConf{}
	dc = DsnConf{
		UserName: "root",
		DbName:   "database",
	}
	if dsn, err := dc.DSN(MySQL); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(dsn)
	}
	// ---
	dc = DsnConf{}
	dc = DsnConf{
		UserName: "postgres",
		DbName:   "database",
	}
	if dsn, err := dc.DSN(PostgreSQL); err != nil {
		t.Error(err.Error())
	} else {
		t.Log(dsn)
	}
}
