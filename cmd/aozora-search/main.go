package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func showAuthors(db *sql.DB) error {
	rows, err := db.Query(`
	SELECT
		a.author_id,
		a.author
	FROM
		authors a
	ORDER BY
		CAST(a.author_id AS INTEGER)
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var authorID, author string
		err = rows.Scan(&authorID, &author)
		if err != nil {
			return err
		}
		fmt.Printf("%s: %s\n", authorID, author)
	}
	return nil
}

const usage = `
Usage of ./aozora-search [sub-command] [...]:
  -d string
        database (default "database.sqlite")

Sub-commands:
    authors
    titles  [AuthorID]
    content [AuthorID] [TitleID]
    query   [Query]
`

func main() {
	var dsn string
	flag.StringVar(&dsn, "d", "database.sqlite", "database")
	flag.Usage = func() {
		fmt.Print(usage)
	}
	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(2)
	}

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	switch flag.Arg(0) {
	case "authors":
		err = showAuthors(db)
	}

	if err != nil {
		log.Fatal(err)
	}
}
