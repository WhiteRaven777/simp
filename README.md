# simp
Simp is a simple SQL operation package.
There is nothing flashy.
It is a package for those who wish to run mere handwritten queries frankly.

# Overview
* Support transactions
* Support commit and rollback
* Support for database open and close
* Execute query manually
* Supported Databases.
    * MySQL
    * PostgreSQL

# Description
WhiteRaven777/simp - I made this package with the desire to run SQL queries more freely.
I know that there are various ORMs.
However, I could not do it well because everything was too functional.
When I tried treating SQL simply, I felt naturally in this form.
People who use only a part of advanced ORM functionality, those who frequently perform complex queries to reduce the number of requests, this package may fit.

# Installation
go get -u github.com/WhiteRaven777/simp

# Documentation
API document and more examples are available here:
http://godoc.org/github.com/WhiteRaven777/simp

# Requirement
This package uses the following standard package.
* database/sql
* errors
* sync
* time

# Usage
## MySQL
```go
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
	
	var db *simp.DB
	if dsn, err := dc.DSN(simp.MySQL); err != nil {
		log.Fatal(err.Error())
	} else {
		if db = simp.New(simp.MySQL, dsn); db.Error() != nil {
			panic(db.Error())
		} else {
			if err = db.Ping(); err != nil {
				fmt.Println("MySQL Connect Error", db.Error().Error())
				panic(db.Error())
			} else {
				fmt.Println("*** Open MySQL Connect ***")
				db.SetConnMaxLifetime(60 * time.Second)
				db.SetMaxIdleConns(5)
			}
		}
	}
	
	// ---
	
	type User struct {
		Id   int
		Name string
	}

	defer func() {
		if err := recover(); err != nil {
			db.Rollback()
		} else {
			db.Commit()
		}
	}()
	if err := db.Begin(); err != nil {
		fmt.Println("error - Begin()", err.Error())
		panic(err)
	}

	var query string
	query = `
		INSERT INTO user (
			id,
			name
		) VALUES (
			?, ?
		), (
			?, ?
		)`

	var row int64
	if r, err := db.Exec(
		query,
		1, "Tom",
		2, "Jerry",
	); err != nil {
		fmt.Println("Insert error", err.Error())
		panic(err.Error())
	} else {
		row, _ = r.RowsAffected()
	}
	fmt.Printf("%d row(s) were inserted.\n", row)

	// ---

	query = `
		SELECT
			id,
			name
		FROM
			user`

	var users []User
	if rows, err := db.Query(query); err != nil {
		fmt.Println("Query error", err.Error())
	} else {
		var user User
		for rows.Next() {
			if err = rows.Scan(
				&user.Id,
				&user.Name,
			); err != nil {
				fmt.Println("Scan error", err.Error())
				continue
			}
			users = append(users, user)
		}
	}

	// ---

	query = `
		SELECT
			count(id)
		FROM
			user`

	var cnt int
	if err := db.QueryRow(query).Scan(&cnt); err != nil {
		fmt.Println("Query error", err.Error())
	}
}
```

# License
Simp is licensed under the MIT