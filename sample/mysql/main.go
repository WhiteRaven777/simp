package main

import (
	"fmt"
	"log"

	"github.com/WhiteRaven777/simp"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dc := simp.DsnConf{
		UserName: "root",
		Password: "pass",
		DbName:   "sample",
		Params: map[string]string{
			"parseTime":            "true",
			"loc":                  "UTC",
			"charset":              "utf8mb4",
			"autocommit":           "false",
			"clientFoundRows":      "true",
			"allowNativePasswords": "True",
		},
	}
	if dsn, err := dc.DSN(simp.MySQL); err != nil {
		log.Fatal(err.Error())
	} else {
		if db := simp.New(simp.MySQL, dsn); db.Error() != nil {
			log.Fatal(db.Error())
		} else {
			if err = db.Ping(); err != nil {
				log.Fatal(err.Error())
			} else {
				fmt.Println(simp.MySQL, "is connected")
			}
		}
	}
}
