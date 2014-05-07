package info

import (
	"labix.org/v2/mgo/bson"
)

// 建议
type Suggestion struct {
	Id         bson.ObjectId `bson:"_id"`
	UserId     bson.ObjectId `UserId`
	Addr       string        `Addr`
	Suggestion string        `Suggestion`
}
