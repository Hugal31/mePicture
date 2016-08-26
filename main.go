package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Hugal31/mePicture/command/picture"
	"github.com/Hugal31/mePicture/command/tag"
)

var help bool

func usage() {
	fmt.Fprintln(os.Stderr, "mePicture is a tool for managing pictures with tag and settup wallpapers slideshows\n"+
		"\n"+
		"Usage:\n"+
		"\n"+
		"\tmePicture command [arguments]\n"+
		"\n"+
		"The commands are:\n"+
		"\n"+
		"\tpicture     List, add and delete pictures\n"+
		"\ttag         List, add and delete tags")
	os.Exit(1)
}

func init() {
	flag.BoolVar(&help, "h", false, "Help")
	flag.Parse()
}

func main() {
	if help || flag.NArg() == 0 {
		usage()
	}

	switch flag.Arg(0) {
	case "picture":
		picture.CommandPicture(flag.Args()[1:])
		break
	case "tag":
		tag.CommandTag(flag.Args()[1:])
		break
	default:
		usage()
	}
}
