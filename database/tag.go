package database

import (
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/Hugal31/mePicture/tag"
)

func (db *DB) getTagId(tagName string) (tagId int) {
	err := db.queryRow("SELECT id FROM tag WHERE name = ?", tagName).Scan(&tagId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Tag %s doesn't exists\n", tagName)
		os.Exit(1)
	}
	return
}

func (db *DB) getTagName(tagId int) (tag string) {
	err := db.queryRow("SELECT name FROM tag WHERE id = ?", tagId).Scan(&tag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}

func (db *DB) TagFromName(name string) tag.Tag {
	t := tag.Tag{Name: name}
	t.Id = db.getTagId(name)
	return t
}

func (db *DB) ListTags() tag.TagSlice {
	var tags tag.TagSlice

	rows, err := db.query("SELECT id,name FROM tag")
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
	stmt, err := db.prepare("INSERT INTO tag(name) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, tagName := range tagNames {
		stmt.Exec(tagName)
	}
}

// Get all the pictures which have this tag and perform PictureRemoveTag on it
// Then, delete the tag
func (db *DB) TagDelete(t tag.Tag) {
	rows, err := db.query("SELECT picture_id FROM picture_tag WHERE tag_id = ?", t.Id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for rows.Next() {
		var pictureId int
		rows.Scan(&pictureId)
		pic := db.PictureFromId(pictureId)
		db.PictureRemoveTag(pic, t)
	}

	db.exec("DELETE FROM tag WHERE id = ?", t.Id)
}
