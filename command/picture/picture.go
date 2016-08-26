package picture

import (
	"path/filepath"
	"fmt"
	"log"
	"os"
	"github.com/Hugal31/mePicture/database"
	"github.com/Hugal31/mePicture/config"
	"github.com/Hugal31/mePicture/picture"
)

func usage() {
	fmt.Fprintln(os.Stderr, "Manage pictures\n"+
		"\n"+
		"Usage:\n"+
		"\n"+
		"\tmePicture picture command [arguments]\n"+
		"\n"+
		"The commands are:\n"+
		"\n"+
		"\tadd target tagName...      Tag pictures\n"+
		"\tlist [tagName...]          List pictures, filter with the tags given in parameter")
	os.Exit(1)
}

func PictureAddTags(path string, tags []string) {
	db := database.Open()
	defer db.Close()

	db.AddTags(tags)

	// TODO Check if is not a directory
	if _, err := os.Stat(path); err != nil {
		log.Fatal(err)
	}
	rel, err := getPicturePath(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s is not in the picture root\n", path)
	}
	db.AddPicture(rel)

	db.AddTagPicture(rel, tags)
}

// Handler for list command
func CommandPicture(args []string) {
	if len(args) == 0 {
		usage()
	}

	switch args[0] {
	case "add":
		pictureAddCommand(args[1:])
		break
	case "list":
		pictureListCommand(args[1:])
		break
	default:
		usage()
	}
}

func pictureAddCommand(args []string) {
	if len(args) < 2 {
		usage()
	}
	PictureAddTags(args[0], args[1:])
}

func pictureListCommand(args []string) {
	db := database.Open()
	defer db.Close()

	var pictures picture.PictureSlice
	if len(args) != 0 {
		pictures = db.ListPictureWithTags(args)
	} else {
		pictures = db.ListPicture()
	}
	for _, pic := range pictures {
		fmt.Print(pic.Name, "\t\t")
		for i, tag :=  range pic.Tags {
			if i == 0 {
				fmt.Print(tag.Name)
			} else {
				fmt.Printf(", %s", tag.Name)
			}
		}
		fmt.Println()
	}
}

func getPicturePath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Rel(config.GetConfig().PicturesRoot, path)
}
