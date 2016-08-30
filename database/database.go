package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Hugal31/mePicture/config"
)

type DB struct {
	sql *sql.DB
	tx  *sql.Tx
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

func (db *DB) Begin() (err error) {
	db.tx, err = db.sql.Begin()
	return
}

func (db *DB) Commit() {
	db.tx.Commit()
	db.tx = nil
}

func (db *DB) Rollback() {
	db.tx.Rollback()
	db.tx = nil
}

func (db *DB) Close() {
	db.sql.Close()
}

func (db *DB) init() {
	sqlStmt := "CREATE TABLE IF NOT EXISTS picture (id INTEGER NOT NULL, path VARCHAR(4096) NOT NULL UNIQUE, PRIMARY KEY (id));" +
		"CREATE TABLE IF NOT EXISTS tag (id INTEGER NOT NULL, name VARCHAR(20) NOT NULL UNIQUE COLLATE NOCASE, PRIMARY KEY (id));" +
		"CREATE TABLE IF NOT EXISTS picture_tag (picture_id INTEGER NOT NULL, tag_id INTEGER NOT NULL," +
		"PRIMARY KEY (picture_id, tag_id), FOREIGN KEY(picture_id) REFERENCES picture(id), FOREIGN KEY(tag_id) REFERENCES tag(id));"
	_, err := db.sql.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (db *DB) query(query string, args ...interface{}) (rows *sql.Rows, err error) {
	if db.tx != nil {
		rows, err = db.tx.Query(query, args...)
	} else {
		rows, err = db.sql.Query(query, args...)
	}
	return
}

func (db *DB) queryRow(query string, args ...interface{}) (row *sql.Row) {
	if db.tx != nil {
		row = db.tx.QueryRow(query, args...)
	} else {
		row = db.sql.QueryRow(query, args...)
	}
	return
}

func (db *DB) prepare(query string) (stmt *sql.Stmt, err error) {
	if db.tx != nil {
		stmt, err = db.tx.Prepare(query)
	} else {
		stmt, err = db.sql.Prepare(query)
	}
	return
}

func (db *DB) exec(query string, args ...interface{}) (result sql.Result, err error) {
	if db.tx != nil {
		result, err = db.tx.Exec(query, args...)
	} else {
		result, err = db.sql.Exec(query, args...)
	}
	return
}

func (db *DB) addLink(pictureId int, tagId int) {
	db.exec("INSERT INTO picture_tag(picture_id, tag_id) VALUES (?, ?)", pictureId, tagId)
}

func (db *DB) removeLink(pictureId int, tagId int) {
	db.exec("DELETE FROM picture_tag WHERE picture_id=? AND tag_id=?", pictureId, tagId)
}
