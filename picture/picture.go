package picture

import "github.com/Hugal31/mePicture/tag"

type Picture struct {
	Id   int
	Name string
	Tags []tag.Tag
}

type PictureSlice []Picture

func (tags PictureSlice) Len() int {
	return len(tags)
}

func (tags PictureSlice) Swap(i, j int) {
	tmp := tags[i]
	tags[i] = tags[j]
	tags[j] = tmp
}

func (tags PictureSlice) Less(i, j int) bool {
	return tags[i].Name < tags[j].Name
}
