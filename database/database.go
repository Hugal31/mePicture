package database

import (
	"database/sql"
	"log"
	"sort"
	"strings"

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
	sort.Strings(tags)
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

func (db *DB) ListPicture() []string {
	var pictures []string

	rows, err := db.sql.Query("SELECT path FROM picture")
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
		pictures = append(pictures, name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	return pictures
}

/*
	This code can be used to perform an IN statement with a undefined number of argument

	statement := "SELECT id FROM tag WHERE name in (?" + strings.Repeat(",?", len(tags)-1) + ")"
	args := make([]interface{}, len(tags))
	for i, v := range tags {
		args[i] = v
	}
	rows, err := db.sql.Query(statement, args...)
 */

// TODO Optimise
func (db *DB) ListPictureWithTags(tags []string) []string {
	var pictures []string

	pictures = db.ListPicture()
	for i := 0; i < len(pictures); {
		pictureTags := db.ListPictureTags(pictures[i])
		var missTag bool
		for _, tag := range tags {
			missTag = true
			for _, pictureTag := range pictureTags {
				if strings.ToLower(tag) == strings.ToLower(pictureTag) {
					missTag = false
					break
				}
			}
			if missTag {
				break
			}
		}
		if missTag {
			pictures = append(pictures[:i], pictures[i+1:]...)
		} else {
			i++
		}
	}
	return pictures
}

func (db *DB) ListPictureTags(picture string) []string {
	pictureId := db.getPictureId(picture)

	var tags []string

	rows, err := db.sql.Query("SELECT tag_id FROM picture_tag WHERE picture_id = ?", pictureId)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var tagId int
		err = rows.Scan(&tagId)
		if err != nil {
			log.Fatal(err)
		}
		tags = append(tags, db.getTagName(tagId))
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(tags)
	return tags
}

func (db *DB) AddPicture(picture string) {
	db.AddPictures([]string{picture})
}

func (db *DB) AddPictures(pictures []string) {
	tx, err := db.sql.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := db.sql.Prepare("INSERT INTO picture(path) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, picture := range pictures {
		stmt.Exec(picture)
	}
	tx.Commit()
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
