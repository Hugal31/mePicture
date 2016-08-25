package database

import (
	"database/sql"
	"log"

	"github.com/Hugal31/mePicture/config"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	sql *sql.DB
}

func Open() *DB {
	var db DB
	var err error
	db.sql, err = sql.Open("sqlite3", config.GetConfig().DatabaseFile)
	if err != nil {
		log.Fatal(err)
	}

	db.init()

	return &db
}

func (db *DB) init() {
	sqlStmt := "CREATE TABLE IF NOT EXISTS picture (id INTEGER PRIMARY KEY, path VARCHAR(4096));" +
		"CREATE TABLE IF NOT EXISTS tag (id INTEGER PRIMARY KEY, name VARCHAR(20) UNIQUE);" +
		"CREATE TABLE IF NOT EXISTS picture_tag (id INTEGER PRIMARY KEY, picture_id INTEGER, tag_id INTEGER," +
		"FOREIGN KEY(picture_id) REFERENCES picture(id), FOREIGN KEY(tag_id) REFERENCES tag(id));"
	_, err := db.sql.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (db *DB) ListTags() []string {
	var tags []string

	rows, err := db.sql.Query("SELECT name FROM tag")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		tags = append(tags, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return tags
}

func (db *DB) AddTag(tagName string) {
	db.AddTags([]string{tagName})
}

func (db *DB) AddTags(tagNames []string) {
	tx, err := db.sql.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("INSERT INTO tag(name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	for _, tagName := range tagNames {
		stmt.Exec(tagName)
	}
	stmt.Close()
	tx.Commit()
}

func (db *DB) Close() {
	db.sql.Close()
}
