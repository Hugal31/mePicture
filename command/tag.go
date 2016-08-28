package command

import (
	"fmt"
	"os"

	"github.com/Hugal31/mePicture/database"
	"github.com/Hugal31/mePicture/tag"
)

func tagUsage() {
	fmt.Fprintln(os.Stderr, "Manage tags\n"+
		"\n"+
		"Usage:\n"+
		"\n"+
		"\tmePicture tag command [arguments]\n"+
		"\n"+
		"The commands are:\n"+
		"\n"+
		"\tadd tagName...      Add a tag\n"+
		"\tlist                List tags")
	os.Exit(1)
}

func checkTagNames(tagNames []string) {
	for _, tagName := range tagNames {
		if !tag.IsValid(tagName) {
			fmt.Fprintln(os.Stderr, "A tag name cannot contain the characters &, |, ( and )")
			os.Exit(1)
		}
	}
}

func ListTags() {
	db := database.Open()
	defer db.Close()

	tags := db.ListTags()
	for _, tag := range tags {
		println(tag.Name)
	}
}

func listTagsCommand([]string) {
	ListTags()
}

func AddTags(tagNames []string) {
	checkTagNames(tagNames)

	db := database.Open()
	defer db.Close()
	db.AddTags(tagNames)
}

func addTagCommand(args []string) {
	if len(args) == 0 {
		tagUsage()
	}
	AddTags(args)
}

func CommandTag(args []string) {
	if len(args) == 0 {
		tagUsage()
	}
	switch args[0] {
	case "list":
		listTagsCommand(args[1:])
		break
	case "add":
		addTagCommand(args[1:])
		break
	default:
		tagUsage()
	}
}
