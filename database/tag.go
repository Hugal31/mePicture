package database

import (
	"log"
	"sort"

	"github.com/Hugal31/mePicture/tag"
)

func (db *DB) ListTags() tag.TagSlice {
	var tags tag.TagSlice

	rows, err := db.sql.Query("SELECT id,name FROM tag")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var t tag.Tag
		err = rows.Scan(&t.Id, &t.Name)
		if err != nil {
			log.Fatal(err)
		}
		tags = append(tags, t)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(tags)
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
	defer stmt.Close()

	for _, tagName := range tagNames {
		stmt.Exec(tagName)
	}
	tx.Commit()
}