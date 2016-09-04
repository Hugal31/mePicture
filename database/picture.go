package database

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/Hugal31/mePicture/picture"
	"github.com/Hugal31/mePicture/tag"
)

func (db *DB) getPictureId(path string) (pictureId int) {
	err := db.queryRow("SELECT id FROM picture WHERE path = ?", path).Scan(&pictureId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Picture %s doesn't exists\n", path)
		os.Exit(1)
	}
	return
}

func (db *DB) getPictureName(pictureId int) (picture string) {
	err := db.queryRow("SELECT path FROM picture WHERE id = ?", pictureId).Scan(&picture)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return
}

func (db *DB) fillPictureTags(pic *picture.Picture) {
	// Retrieve tags
	tagRows, err := db.query("SELECT tag_id FROM picture_tag WHERE picture_id=?", pic.Id)
	if err != nil {
		log.Fatal(err)
	}
	for tagRows.Next() {
		var t tag.Tag
		tagRows.Scan(&t.Id)
		t.Name = db.getTagName(t.Id)
		pic.Tags = append(pic.Tags, t)
	}
}

func (db *DB) PictureFromPath(path string) picture.Picture {
	pic := picture.Picture{Id: db.getPictureId(path), Name: path}
	db.fillPictureTags(&pic)
	return pic
}

func (db *DB) PictureFromId(id int) picture.Picture {
	pic := picture.Picture{Id: id, Name: db.getPictureName(id)}
	db.fillPictureTags(&pic)
	return pic
}

func (db *DB) ListPicture() picture.PictureSlice {
	var pictures picture.PictureSlice

	rows, err := db.query("SELECT id,path FROM picture")
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
		tagRows, err := db.query("SELECT tag_id FROM picture_tag WHERE picture_id=?", pic.Id)
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
	rows, err := db.Query(statement, args...)
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

func (db *DB) PictureAdd(path string) picture.Picture {
	db.PicturesAdd([]string{path})
	return db.PictureFromPath(path)
}

func (db *DB) PicturesAdd(pictures []string) {
	stmt, err := db.prepare("INSERT INTO picture(path) VALUES (?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	for _, pic := range pictures {
		stmt.Exec(pic)
	}
}

func (db *DB) PictureDelete(pic *picture.Picture) {
	db.exec("DELETE FROM picture_tag WHERE picture_id = ?;"+
		"DELETE FROM picture WHERE id = ?", pic.Id, pic.Id)
}

func (db *DB) PictureAddTags(pic *picture.Picture, tags tag.TagSlice) {
	for _, t := range tags {
		isFound := false
		for _, picTag := range pic.Tags {
			if picTag.Id == t.Id {
				isFound = true
				break
			}
		}
		if !isFound {
			db.addLink(pic, &t)
			pic.Tags = append(pic.Tags, t)
		}
	}
}

func (db *DB) PictureRemoveTag(pic *picture.Picture, t *tag.Tag) {
	for i := 0; i < pic.Tags.Len(); i++ {
		if pic.Tags[i].Id == t.Id {
			db.removeLink(pic, t)
			pic.Tags = append(pic.Tags[:i], pic.Tags[i+1:]...)
			break
		}
	}
	if len(pic.Tags) == 0 {
		db.PictureDelete(pic)
	}
}
