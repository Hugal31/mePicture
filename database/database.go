package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/Hugal31/mePicture/config"
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

func (db *DB) getPictureId(picture string) (pictureId int) {
	err := db.sql.QueryRow("SELECT id FROM picture WHERE path = ?", picture).Scan(&pictureId)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (db *DB) getPictureName(pictureId int) (picture string) {
	err := db.sql.QueryRow("SELECT path FROM picture WHERE id = ?", pictureId).Scan(&picture)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (db *DB) getTagId(tag string) (tagId int) {
	err := db.sql.QueryRow("SELECT id FROM tag WHERE name = ?", tag).Scan(&tagId)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (db *DB) getTagName(tagId int) (tag string) {
	err := db.sql.QueryRow("SELECT name FROM tag WHERE id = ?", tagId).Scan(&tag)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (db *DB) addLink(pictureId int, tagId int) {
	db.sql.Exec("INSERT INTO picture_tag(picture_id, tag_id) VALUES (?, ?)", pictureId, tagId)
}

func (db *DB) AddTagPicture(picture string, tags []string) {
	pictureId := db.getPictureId(picture)

	stmt, err := db.sql.Prepare("SELECT id FROM tag WHERE name = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for _, tag := range tags {
		var tagId int
		err := stmt.QueryRow(tag).Scan(&tagId)
		if err != nil {
			log.Fatal(err)
		}
		db.addLink(pictureId, tagId)
	}
}

func (db *DB) Close() {
	db.sql.Close()
}
