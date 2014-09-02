package info

import (
	"gopkg.in/mgo.v2/bson"
)

// 建议
type Suggestion struct {
	Id         bson.ObjectId `bson:"_id"`
	UserId     bson.ObjectId `UserId`
	Addr       string        `Addr`
	Suggestion string        `Suggestion`
}
