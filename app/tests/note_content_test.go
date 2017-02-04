package tests

import (
    "github.com/leanote/leanote/app/db"
    "github.com/revel/revel"
    "testing"
    //  . "github.com/leanote/leanote/app/lea"
    "github.com/leanote/leanote/app/service"
    // "regexp"
    //  "gopkg.in/mgo.v2"
    // "fmt"
    // "strings"
)


// 可在server端调试

func init() {
    revel.Init("dev", "github.com/leanote/leanote", "/Users/life/Documents/Go/package_base/src")
    db.Init("mongodb://localhost:27017/leanote", "leanote")
    service.InitService()
    service.ConfigS.InitGlobalConfigs()
}

func TestApiFixNoteContent2(t *testing.T) {
    note2 := service.NoteS.GetNote("585df83771c1b17e8a000000", "585df81199c37b6176000004")
    note := service.NoteS.GetNoteContent("585df83771c1b17e8a000000", "585df81199c37b6176000004")
    contentFixed := service.NoteS.FixContent(note.Content, false)
    t.Log(note2.Title)
    t.Log(note.Content)
    t.Log(contentFixed)
}