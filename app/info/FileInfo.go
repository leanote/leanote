package info

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type File struct {
	FileId         bson.ObjectId `bson:"_id,omitempty"` //
	UserId         bson.ObjectId `bson:"UserId"`
	AlbumId        bson.ObjectId `bson:"AlbumId"`
	Name           string        `Name`  // file name
	Title          string        `Title` // file name or user defined for search
	Size           int64           `Size`  // file size (byte)
	Type           string        `Type`  // file type, such as image/jpg
	Path           string        `Path`  // the file path, based on /upload
	IsDefaultAlbum bool          `IsDefaultAlbum`
	CreatedTime    time.Time     `CreatedTime`
}
