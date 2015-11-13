package db

import (
	"fmt"
	. "github.com/leanote/leanote/app/lea"
	"github.com/revel/revel"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

// Init mgo and the common DAO

// 数据连接
var Session *mgo.Session

// 各个表的Collection对象
var Notebooks *mgo.Collection
var Notes *mgo.Collection
var NoteContents *mgo.Collection
var NoteContentHistories *mgo.Collection

var ShareNotes *mgo.Collection
var ShareNotebooks *mgo.Collection
var HasShareNotes *mgo.Collection
var Blogs *mgo.Collection
var Users *mgo.Collection
var Groups *mgo.Collection
var GroupUsers *mgo.Collection

var Tags *mgo.Collection
var NoteTags *mgo.Collection
var TagCounts *mgo.Collection

var UserBlogs *mgo.Collection

var Tokens *mgo.Collection

var Suggestions *mgo.Collection

// Album & file(image)
var Albums *mgo.Collection
var Files *mgo.Collection
var Attachs *mgo.Collection

var NoteImages *mgo.Collection
var Configs *mgo.Collection
var EmailLogs *mgo.Collection

// blog
var BlogLikes *mgo.Collection
var BlogComments *mgo.Collection
var Reports *mgo.Collection
var BlogSingles *mgo.Collection
var Themes *mgo.Collection

// session
var Sessions *mgo.Collection

// 初始化时连接数据库
func Init(url, dbname string) {
	ok := true
	config := revel.Config
	if url == "" {
		url, ok = config.String("db.url")
		if !ok {
			url, ok = config.String("db.urlEnv")
			if ok {
				Log("get db conf from urlEnv: " + url)
			}
		} else {
			Log("get db conf from db.url: " + url)
		}

		if ok {
			// get dbname from urlEnv
			urls := strings.Split(url, "/")
			dbname = urls[len(urls)-1]
		}
	}
	if dbname == "" {
		dbname, _ = config.String("db.dbname")
	}

	// get db config from host, port, username, password
	if !ok {
		host, _ := revel.Config.String("db.host")
		port, _ := revel.Config.String("db.port")
		username, _ := revel.Config.String("db.username")
		password, _ := revel.Config.String("db.password")
		usernameAndPassword := username + ":" + password + "@"
		if username == "" || password == "" {
			usernameAndPassword = ""
		}
		url = "mongodb://" + usernameAndPassword + host + ":" + port + "/" + dbname
	}
	Log(url)

	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	var err error
	Session, err = mgo.Dial(url)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	Session.SetMode(mgo.Monotonic, true)

	// notebook
	Notebooks = Session.DB(dbname).C("notebooks")

	// notes
	Notes = Session.DB(dbname).C("notes")

	// noteContents
	NoteContents = Session.DB(dbname).C("note_contents")
	NoteContentHistories = Session.DB(dbname).C("note_content_histories")

	// share
	ShareNotes = Session.DB(dbname).C("share_notes")
	ShareNotebooks = Session.DB(dbname).C("share_notebooks")
	HasShareNotes = Session.DB(dbname).C("has_share_notes")

	// user
	Users = Session.DB(dbname).C("users")
	// group
	Groups = Session.DB(dbname).C("groups")
	GroupUsers = Session.DB(dbname).C("group_users")

	// blog
	Blogs = Session.DB(dbname).C("blogs")

	// tag
	Tags = Session.DB(dbname).C("tags")
	NoteTags = Session.DB(dbname).C("note_tags")
	TagCounts = Session.DB(dbname).C("tag_count")

	// blog
	UserBlogs = Session.DB(dbname).C("user_blogs")
	BlogSingles = Session.DB(dbname).C("blog_singles")
	Themes = Session.DB(dbname).C("themes")

	// find password
	Tokens = Session.DB(dbname).C("tokens")

	// Suggestion
	Suggestions = Session.DB(dbname).C("suggestions")

	// Album & file
	Albums = Session.DB(dbname).C("albums")
	Files = Session.DB(dbname).C("files")
	Attachs = Session.DB(dbname).C("attachs")

	NoteImages = Session.DB(dbname).C("note_images")

	Configs = Session.DB(dbname).C("configs")
	EmailLogs = Session.DB(dbname).C("email_logs")

	// 社交
	BlogLikes = Session.DB(dbname).C("blog_likes")
	BlogComments = Session.DB(dbname).C("blog_comments")

	// 举报
	Reports = Session.DB(dbname).C("reports")

	// session
	Sessions = Session.DB(dbname).C("sessions")
}

func close() {
	Session.Close()
}

// common DAO
// 公用方法

//----------------------

func Insert(collection *mgo.Collection, i interface{}) bool {
	err := collection.Insert(i)
	return Err(err)
}

//----------------------

// 适合一条记录全部更新
func Update(collection *mgo.Collection, query interface{}, i interface{}) bool {
	err := collection.Update(query, i)
	return Err(err)
}
func Upsert(collection *mgo.Collection, query interface{}, i interface{}) bool {
	_, err := collection.Upsert(query, i)
	return Err(err)
}
func UpdateAll(collection *mgo.Collection, query interface{}, i interface{}) bool {
	_, err := collection.UpdateAll(query, i)
	return Err(err)
}
func UpdateByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) bool {
	err := collection.Update(GetIdAndUserIdQ(id, userId), i)
	return Err(err)
}

func UpdateByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) bool {
	err := collection.Update(GetIdAndUserIdBsonQ(id, userId), i)
	return Err(err)
}
func UpdateByIdAndUserIdField(collection *mgo.Collection, id, userId, field string, value interface{}) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap(collection *mgo.Collection, id, userId string, v bson.M) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": v})
}

func UpdateByIdAndUserIdField2(collection *mgo.Collection, id, userId bson.ObjectId, field string, value interface{}) bool {
	return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap2(collection *mgo.Collection, id, userId bson.ObjectId, v bson.M) bool {
	return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": v})
}

//
func UpdateByQField(collection *mgo.Collection, q interface{}, field string, value interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": bson.M{field: value}})
	return Err(err)
}
func UpdateByQI(collection *mgo.Collection, q interface{}, v interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": v})
	return Err(err)
}

// 查询条件和值
func UpdateByQMap(collection *mgo.Collection, q interface{}, v interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": v})
	return Err(err)
}

//------------------------

// 删除一条
func Delete(collection *mgo.Collection, q interface{}) bool {
	err := collection.Remove(q)
	return Err(err)
}
func DeleteByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
	err := collection.Remove(GetIdAndUserIdQ(id, userId))
	return Err(err)
}
func DeleteByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
	err := collection.Remove(GetIdAndUserIdBsonQ(id, userId))
	return Err(err)
}

// 删除所有
func DeleteAllByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
	_, err := collection.RemoveAll(GetIdAndUserIdQ(id, userId))
	return Err(err)
}
func DeleteAllByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
	_, err := collection.RemoveAll(GetIdAndUserIdBsonQ(id, userId))
	return Err(err)
}

func DeleteAll(collection *mgo.Collection, q interface{}) bool {
	_, err := collection.RemoveAll(q)
	return Err(err)
}

//-------------------------

func Get(collection *mgo.Collection, id string, i interface{}) {
	collection.FindId(bson.ObjectIdHex(id)).One(i)
}
func Get2(collection *mgo.Collection, id bson.ObjectId, i interface{}) {
	collection.FindId(id).One(i)
}

func GetByQ(collection *mgo.Collection, q interface{}, i interface{}) {
	collection.Find(q).One(i)
}
func ListByQ(collection *mgo.Collection, q interface{}, i interface{}) {
	collection.Find(q).All(i)
}

func ListByQLimit(collection *mgo.Collection, q interface{}, i interface{}, limit int) {
	collection.Find(q).Limit(limit).All(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func GetByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	collection.Find(q).Select(selector).One(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func ListByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	collection.Find(q).Select(selector).All(i)
}
func GetByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) {
	collection.Find(GetIdAndUserIdQ(id, userId)).One(i)
}
func GetByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) {
	collection.Find(GetIdAndUserIdBsonQ(id, userId)).One(i)
}

// 按field去重
func Distinct(collection *mgo.Collection, q bson.M, field string, i interface{}) {
	collection.Find(q).Distinct(field, i)
}

//----------------------

func Count(collection *mgo.Collection, q interface{}) int {
	cnt, err := collection.Find(q).Count()
	if err != nil {
		Err(err)
	}
	return cnt
}

func Has(collection *mgo.Collection, q interface{}) bool {
	if Count(collection, q) > 0 {
		return true
	}
	return false
}

//-----------------

// 得到主键和userId的复合查询条件
func GetIdAndUserIdQ(id, userId string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id), "UserId": bson.ObjectIdHex(userId)}
}
func GetIdAndUserIdBsonQ(id, userId bson.ObjectId) bson.M {
	return bson.M{"_id": id, "UserId": userId}
}

// DB处理错误
func Err(err error) bool {
	if err != nil {
		fmt.Println(err)
		// 删除时, 查找
		if err.Error() == "not found" {
			return true
		}
		return false
	}
	return true
}

// 检查mognodb是否lost connection
// 每个请求之前都要检查!!
func CheckMongoSessionLost() {
	// fmt.Println("检查CheckMongoSessionLostErr")
	err := Session.Ping()
	if err != nil {
		Log("Lost connection to db!")
		Session.Refresh()
		err = Session.Ping()
		if err == nil {
			Log("Reconnect to db successful.")
		} else {
			Log("重连失败!!!! 警告")
		}
	}
}
