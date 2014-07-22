package main

import (
	"github.com/robfig/revel"

	"leanote/app/info"
	. "leanote/app/lea"
	"leanote/app/service"
	//	"github.com/leanote/leanote/app/lea/memcache"
	//	"github.com/leanote/leanote/app/lea/netutil"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"leanote/app/lea/html2image"
	"time"
	//	"labix.org/v2/mgo"
	//	"encoding/json"
	//	"strings"
)

var userId = "5295d4c95b1dd58edb4a7f4f"
var userIdO = bson.ObjectIdHex(userId)

var userId3 = "52b43ae4cfeeae33ef073b2b"
var userId3O = bson.ObjectIdHex(userId3)

var toUserId = "52b4376bcfeeae33ef073b21"
var toUserIdO = bson.ObjectIdHex(toUserId)

var notebookId = "52b2d051ea3ba3d3fb35910c"
var notebookIdO = bson.ObjectIdHex(notebookId)
var notebookId2O = bson.ObjectIdHex("52b2d051ea3ba3d3fb35910b")

var noteId = "52b2dd34ea3ba3d3fb35910d"
var noteIdO = bson.ObjectIdHex(noteId)

func testNoteService() {
	println("testNoteService")
	noteService := &service.NoteService{}
	note := noteService.GetBlogNote("535f9e6b19807a4c8d000000")
	LogJ(note)
	/*
		_, notes := noteService.SearchNote("go", "52d3e8ac99c37b7f0d000001", 1, 30, "", false, true);
		LogJ(notes)
	*/
	return

	/*
		noteService.AddNote(info.Note{UserId: userIdO,
			NotebookId: notebookIdO,
			Title: "life you", Tags: []string{"red", "yellow"}})
	*/

	noteService.AddNoteContent(info.NoteContent{UserId: userIdO,
		NoteId:  bson.ObjectIdHex("52b4531dcfeeae33ef073b33"),
		Content: "xxxxxxxxxxxxxxxxxxxxlifeyou can m<div><p></p></div>"})

	//	noteService.UpdateNoteContent(userId, userId, "52b2dd34ea3ba3d3fb35910d", "life2---------")

	//    noteService.AddNote(info.Note{Title: "life", Tags: []string{"life", "life2"}})

	//	println(bson.IsObjectIdHex(id))

	//	note := noteService.Get(id)
	//	fmt.Println(note)

	//	noteService.UpdateTags(id, []string{"lifedd", "life2"});
}

func testNotebookService() {
	service := &service.NotebookService{}
	Log(service.IsBlog("52ccb959bcbf21610d000001"))

	return

	notebooks := service.GetNotebooks(userId3)
	LogJ(notebooks)

	//	b, _ := json.MarshalIndent(notebooks, "", "	")
	//	fmt.Println(string(b))

	// service.UpdateNotebookTitle("52b2cf9eea3ba3d3fb359108", userId, "JS")

	//	service.SortNotebooks(userId, map[string]int{"52b2d051ea3ba3d3fb35910c": 4, "52b2d051ea3ba3d3fb35910b": 3})

	/*
		service.AddNotebook(info.Notebook{UserId: userId3O,
			ParentNotebookId: bson.ObjectIdHex("52b43c1fcfeeae33ef073b2e"),
			Title: "Mac-life",
			Seq: 0})
	*/

	println("xxx")
}

func testShareService() {
	service := &service.ShareService{}

	//	service.AddShareNote("52bd127dbcbf216d0b000000", 1, "5295d4c95b1dd58edb4a7f4f", "c@a.com")
	LogJ(service.ListNoteShareUserInfo("52cd11a1bcbf215680000000", "5295d4c95b1dd58edb4a7f4f"))
	return

	/*
		service.AddShareNotebook(info.ShareNotebook{UserId: userId3O, ToUserId: userIdO,
			NotebookId: bson.ObjectIdHex("52b43c38cfeeae33ef073b30")})
	*/

	notebooks, _ := service.GetShareNotebooks(userId)
	LogJ(notebooks)

	//	noteService.AddNoteContent(info.NoteContent{UserId: userIdO,
	//		NoteId: bson.ObjectIdHex("52b2dd34ea3ba3d3fb35910d"), Content:"xxxxxxxxxxxxxxxxxxxxlifeyou can m<div><p></p></div>"})

	//    noteService.AddNote(info.Note{Title: "life", Tags: []string{"life", "life2"}})

	//	println(bson.IsObjectIdHex(id))

	//	note := noteService.Get(id)
	//	fmt.Println(note)

	//	noteService.UpdateTags(id, []string{"lifedd", "life2"});
}

func testAuthService() {
	userService := &service.UserService{}
	LogJ(userService.GetUserInfo("52d26b4e99c37b609a000001"))
	// userService.AddUser(info.User{UserId: bson.ObjectIdHex("52d26b4e99c37b609a000001"), Email: "leanote@leanote.com", Pwd:"abc"})
	return

	authService := &service.AuthService{}
	authService.Register("f@a.com", "abc")

	//	fmt.Println(authService.LogonGetUserInfo("a@a.com", "abc"))

	//	noteService.AddNoteContent(info.NoteContent{UserId: userIdO,
	//		NoteId: bson.ObjectIdHex("52b2dd34ea3ba3d3fb35910d"), Content:"xxxxxxxxxxxxxxxxxxxxlifeyou can m<div><p></p></div>"})

	//    noteService.AddNote(info.Note{Title: "life", Tags: []string{"life", "life2"}})

	//	println(bson.IsObjectIdHex(id))

	//	note := noteService.Get(id)
	//	fmt.Println(note)

	//	noteService.UpdateTags(id, []string{"lifedd", "life2"});
}

func testTagService() {
	service := &service.TagService{}
	//	service.AddTags("5295d4c95b1dd58edb4a7f4f", []string{"life", "blue", "yellow"})
	//	service.AddTags("5295d4c95b1dd58edb4a7f4f", []string{"what", "can", "make"})
	LogJ(service.GetTags("5295d4c95b1dd58edb4a7f4f"))
}

func testHtml2Image() {
	start := time.Now()
	Log("start...")
	//	TestFillString()

	html2image.ToImage("uid", "username", "noteId", "开发一款属于自己的编程语言，开发一款属于自己的编程语言听起来是不是很酷？", `
		<div class="each-post">
			<p>
			一个合格的 Techspace 需要有足够专业的器材、场地和资源，你可以和你的团队在里面进行激光切割、快速贴片甚至加工木材等操作，在相对独立的空间内又能同周围的同道友人互相激发切磋。国内现有的 Techspace 没几家，不久前我去深圳特地拜访了当地的 Techspace，很喜欢那里的氛围，希望国内其他地方也能有更多这类空间供创客发挥。
假如你有一个比较成型的想法，想在硬件领域做点事情，核心团队也基本组好，硬件软件交互基本都有专人了。</p>
			<p><a>这时候你的首要目标</a>，就是找个地方，按照你的计划，尽早做出一个可用的原型。Techspace可能是一个合适的地方。
	一个合格的 Techspace 需要有足够专业的器材、场地和资源，你可以和你的团队在里面进行激光切割、快速贴片甚至加工木材等操作，在相对独立的空间内又能同周围的同道友人互相激发切磋。国内现有的 Techspace 没几家，不久前我去深圳特地拜访了当地的 Techspace，很喜欢那里的氛围，希望国内其他地方也能有更多这类空间供创客发挥。
	深圳 Techspace 位于工业区，园区内有奥迪、BMW 等企业的厂房，在大门口我停下来问保安，M10 栋在哪里？保安答，去 Techspace？最里面靠右手那栋。惊叹于保安的机智，我也在想，莫非有许多朋友都慕名来到这巨大园区寻访 Techspace？
穿过一片工业区里的高楼和各种建筑材料
			</p>
			<!--
			<pre class="">cd jpeg-9a/<br>./configure --enable-shared --enable-static
make<br>make --- install</pre>
			<p>
			life you can, !!@kk
			</p>
			
			-->
		</div>
	`, "/Users/life/Desktop/a.png")
	fmt.Printf("time cost %v\n", time.Now().Sub(start))
}

func testLea() {
	names := ListDir("/Users/life/Documents/Go/package/src/leanote")
	fmt.Println(names)
}

func main() {
	revel.BasePath = "/Users/life/Documents/Go/package/src/leanote"
	testLea()
	//	a, b := SplitFilename("http://ab/c/a.gif#??")
	//	println(a)
	//	println(b)

	//	path, ok := netutil.WriteUrl("http://a.36krcnd.com/photo/2014/9bd1a07c0973d79ca05ad13c3c2e16b8.png!slider", "/tmp")
	//	println(path)
	//	testHtml2Image();

	//	println(IsObjectId("52d26b4e99c37b609a000001"))

	//	b := `请点击链接验证邮箱: <a href="">http://leanote.com/user/activeEmail?token=d8ca086cce5550a6227f9dc84dbac09d</a>. 48小时后过期.`
	//	SendEmail("lifephp@gmail.com", "leanote-验证邮箱", "验证邮", b)

	//_, err := mgo.Dial("mongodb://leanote:nKFAkxKnWkEQy8Vv2LlM@115.28.133.226:27017/leanote")

	//	testNotebookService();
	// testNoteService();
	//	testShareService()
	//	testAuthService()

	//	testTagService();
	/*
		filename := "你好59.26.png"
		ext := SubstringByte(filename, strings.LastIndex(filename, "."))
		ext = strings.ToLower(ext)
		print(ext)
		52d26ab199c37b5f80000001

		Log(bson.NewObjectId())
		Log(bson.NewObjectId())
		Log(bson.NewObjectId())

	*/
	//	Log(TransferExt("/你好a/b/a.gif", ".jpg"))
	//	TransToGif("/Users/life/Desktop/a2.png", 0, false)

	//	Log(IsUsername("xx**x"))

	/*
		memcache.Set("xx", map[string]string{"A":"you"}, 0)
		a := memcache.Get("xx")
		Log(a)
	*/
}
