package command

import (
	"fmt"
	"os"

	"github.com/ogier/pflag"
)

var help bool
var hideTagNames bool
var displayFullPath bool

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
	pflag.BoolVarP(&help, "help", "h", false, "Help")
	pflag.BoolVarP(&hideTagNames, "plain", "p", false, "Hide tag names")
	pflag.BoolVarP(&displayFullPath, "fullpath", "f", false, "Display full path")
	pflag.Parse()
}

func Run() {
	if help || pflag.NArg() == 0 {
		usage()
	}

	switch pflag.Arg(0) {
	case "picture":
		CommandPicture(pflag.Args()[1:])
		break
	case "tag":
		CommandTag(pflag.Args()[1:])
		break
	default:
		usage()
	}
}
