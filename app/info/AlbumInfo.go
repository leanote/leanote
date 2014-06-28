package info

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Album struct {
	AlbumId     bson.ObjectId `bson:"_id,omitempty"` //
	UserId      bson.ObjectId `bson:"UserId"`
	Name        string        `Name` // album name
	Type        int           `Type` // type, the default is image: 0
	Seq         int           `Seq`
	CreatedTime time.Time     `CreatedTime`
}
