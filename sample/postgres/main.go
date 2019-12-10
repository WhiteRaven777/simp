package main

import (
	"fmt"
	"log"

	"github.com/WhiteRaven777/simp"
	_ "github.com/lib/pq"
)

func main() {
	dc := simp.DsnConf{
		UserName: "postgres",
		Password: "pass",
		DbName:   "sample",
		Params: map[string]string{
			"sslmode": "disable",
		},
	}
	if dsn, err := dc.DSN(simp.PostgreSQL); err != nil {
		log.Fatal(err.Error())
	} else {
		if db := simp.New(simp.PostgreSQL, dsn); db.Error() != nil {
			log.Fatal(db.Error())
		} else {
			if err = db.Ping(); err != nil {
				log.Fatal(err.Error())
			} else {
				fmt.Println(simp.PostgreSQL, "is connected")
			}
		}
	}
}
