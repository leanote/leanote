package controllers

import (
	"github.com/revel/revel"
	//	"encoding/json"
	"github.com/leanote/leanote/app/info"
	"gopkg.in/mgo.v2/bson"
	//	. "github.com/leanote/leanote/app/lea"
	//	"io/ioutil"
)

// Album controller
type Album struct {
	BaseController
}

// 图片管理, iframe
func (c Album) Index() revel.Result {
	c.SetLocale()
	return c.RenderTemplate("album/index.html")
}

// all albums by userId
func (c Album) GetAlbums() revel.Result {
	re := albumService.GetAlbums(c.GetUserId())
	return c.RenderJSON(re)
}
func (c Album) DeleteAlbum(albumId string) revel.Result {
	re, msg := albumService.DeleteAlbum(c.GetUserId(), albumId)
	return c.RenderJSON(info.Re{Ok: re, Msg: msg})
}

// add album
func (c Album) AddAlbum(name string) revel.Result {
	album := info.Album{
		AlbumId: bson.NewObjectId(),
		Name:    name,
		Seq:     -1,
		UserId:  c.GetObjectUserId()}
	re := albumService.AddAlbum(album)

	if re {
		return c.RenderJSON(album)
	} else {
		return c.RenderJSON(false)
	}
}

// update alnum name
func (c Album) UpdateAlbum(albumId, name string) revel.Result {
	return c.RenderJSON(albumService.UpdateAlbum(albumId, c.GetUserId(), name))
}
