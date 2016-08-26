package tag

type Tag struct {
	Id   int
	Name string
}

type TagSlice []Tag

func (tags TagSlice) Len() int {
	return len(tags)
}

func (tags TagSlice) Swap(i, j int) {
	tmp := tags[i]
	tags[i] = tags[j]
	tags[j] = tmp
}

func (tags TagSlice) Less(i, j int) bool {
	return tags[i].Name < tags[j].Name
}
