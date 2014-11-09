package info

import (
//	"time"
)

// 仅仅为了博客的主题

type BlogInfoCustom struct {
	UserId string
	Username string
	UserLogo string
	Title string
	SubTitle string
	Logo string
	OpenComment bool
	CommentType string
	ThemeId string
	SubDomain string
	Domain string
}

type Post struct {
	NoteId string
}
// 归档
type Archive struct {
	Year int	
	Posts []map[string]interface{}
}

type TagsCounts []TagCount
type TagCount struct {
	Tag string
	Count int
}
func (this TagsCounts) Len() int {
	return len(this)
}
func (this TagsCounts) Less(i, j int) bool {
	return this[i].Count > this[j].Count
}
func (this TagsCounts) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}