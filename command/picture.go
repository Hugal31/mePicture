package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Hugal31/mePicture/config"
	"github.com/Hugal31/mePicture/database"
	"github.com/Hugal31/mePicture/picture"
	"github.com/Hugal31/mePicture/tag"
)

func pictureUsage() {
	fmt.Fprintln(os.Stderr, "Manage pictures\n"+
		"\n"+
		"Usage:\n"+
		"\n"+
		"\tmePicture picture command [arguments]\n"+
		"\n"+
		"The commands are:\n"+
		"\n"+
		"\tadd target tagName...        Tag pictures\n"+
		"\tlist [-p] [-f] [tagName...]  List pictures, filter with the tags given in parameter\n"+
		"\t                             Add -p to prevent the display of tag namse\n"+
		"\t                             Add -f to show full path\n"+
		"\tremove target tagName...     Remove tags from picture\n"+
		"\tdelete target...             Remove all tags from target\n"+
		"\n"+
		"target:  Image file or directory")
	os.Exit(1)
}

// Handler for list command
func CommandPicture(args []string) {
	if len(args) == 0 {
		pictureUsage()
	}

	switch args[0] {
	case "add":
		pictureAddCommand(args[1:])
		break
	case "list":
		pictureListCommand(args[1:])
		break
	case "remove":
		pictureRemoveCommand(args[1:])
		break
	case "delete":
		pictureDeleteCommand(args[1:])
		break
	default:
		pictureUsage()
	}
}

func addPictureTag(path string, tags tag.TagSlice, db *database.DB) {
	rel := getPicturePath(path)
	pic := db.PictureAdd(rel)
	db.PictureAddTags(&pic, tags)
}

func walk(path string, file os.FileInfo, tags tag.TagSlice, db *database.DB, f func(string, tag.TagSlice, *database.DB)) {
	if file.IsDir() {
		subFiles, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, subfile := range subFiles {
			walk(path+string(os.PathSeparator)+subfile.Name(), subfile, tags, db, f)
		}
	} else if picture.IsValidExt(filepath.Ext(path)) {
		f(path, tags, db)
	}
}

func PictureAddTags(path string, tagNames []string) {
	checkTagNames(tagNames)

	file, err := os.Stat(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	db := database.Open()
	defer db.Close()

	db.Begin()
	db.AddTags(tagNames)
	tags := db.TagsFromNames(tagNames)
	walk(path, file, tags, db, addPictureTag)
	db.Commit()
}

func pictureAddCommand(args []string) {
	if len(args) < 2 {
		pictureUsage()
	}
	PictureAddTags(args[0], args[1:])
}

func PictureList(tags []string) {
	db := database.Open()
	defer db.Close()

	var pictures picture.PictureSlice
	if len(tags) != 0 {
		pictures = db.ListPictureWithTags(tags)
	} else {
		pictures = db.ListPicture()
	}
	for _, pic := range pictures {
		if displayFullPath {
			fmt.Printf("%s%c%s", config.GetConfig().PicturesRoot, os.PathSeparator, pic.Name)
		} else {
			fmt.Print(pic.Name)
		}
		if !hideTagNames {
			fmt.Print(", \t\t")
			for i, t := range pic.Tags {
				if i == 0 {
					fmt.Print(t.Name)
				} else {
					fmt.Printf(", %s", t.Name)
				}
			}
		}
		fmt.Println()
	}
}

func pictureListCommand(args []string) {
	PictureList(args)
}

func pictureRemoveTags(path string, tags tag.TagSlice, db *database.DB) {
	rel := getPicturePath(path)
	pic := db.PictureFromPath(rel)

	for _, t := range tags {
		db.PictureRemoveTag(&pic, &t)
	}
}

func PictureRemoveTags(path string, tagNames []string) {
	file, err := os.Stat(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	db := database.Open()
	defer db.Close()

	db.Begin()
	tags := db.TagsFromNames(tagNames)
	walk(path, file, tags, db, pictureRemoveTags)
	db.Commit()
}

func pictureRemoveCommand(args []string) {
	if len(args) < 2 {
		pictureUsage()
	}

	PictureRemoveTags(args[0], args[1:])
}

func pictureDeleteCommand(args []string) { // TODO Implement recursive
	if len(args) < 1 {
		pictureUsage()
	}

	db := database.Open()
	defer db.Close()

	for _, path := range args {
		path := getPicturePath(path)
		pic := db.PictureFromPath(path)
		db.PictureDelete(&pic)
	}
}

func getPicturePath(path string) string {
	path, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rel, err := filepath.Rel(config.GetConfig().PicturesRoot, path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is not in the picture root\n", path)
		os.Exit(1)
	}
	return rel
}
