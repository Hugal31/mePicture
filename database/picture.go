package database

import (
	"log"
	"strings"
	"sort"

	"github.com/Hugal31/mePicture/picture"
	"github.com/Hugal31/mePicture/tag"
)

func (db *DB) ListPicture() picture.PictureSlice {
	var pictures picture.PictureSlice

	rows, err := db.sql.Query("SELECT id,path FROM picture")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var pic picture.Picture
		err = rows.Scan(&pic.Id, &pic.Name)
		if err != nil {
			log.Fatal(err)
		}

		// Retrieve tags
		tagRows, err := db.sql.Query("SELECT tag_id FROM picture_tag WHERE picture_id=?", pic.Id)
		if err != nil {
			log.Fatal(err)
		}
		for tagRows.Next() {
			var t tag.Tag
			tagRows.Scan(&t.Id)
			t.Name = db.getTagName(t.Id)
			pic.Tags = append(pic.Tags, t)
		}
		tagRows.Close()
		pictures = append(pictures, pic)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(pictures)
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
// TODO Use filter
func (db *DB) ListPictureWithTags(tags []string) picture.PictureSlice {
	pics := db.ListPicture()
	for i := 0; i < len(pics); {
		var missTag bool
		for _, t := range tags {
			missTag = true
			for _, pictureTag := range pics[i].Tags {
				if strings.ToLower(t) == strings.ToLower(pictureTag.Name) {
					missTag = false
					break
				}
			}
			if missTag {
				break
			}
		}
		if missTag {
			pics = append(pics[:i], pics[i+1:]...)
		} else {
			i++
		}
	}
	return pics
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

	for _, pic := range pictures {
		stmt.Exec(pic)
	}
	tx.Commit()
}
