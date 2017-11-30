modules/gorp
===============

This composite orm module combines Gorp with Squirrel to give you a complete solution  
* [Gorp](https://github.com/go-gorp/gorp) *saves you time, minimizes the drudgery of 
getting data in and out of your database, and helps your code focus on algorithms, 
not infrastructure*. 
* [Squirrel](https://github.com/Masterminds/squirrel)  *helps you build SQL queries from 
composable parts*

## Activation
```ini
module.gorp = github.com/revel/modules/orm/gorp
```

## Drivers

* sqlite3
* postgres
* mysql

## Configuration file

```ini
# Database config
db.autoinit=true # default=true
db.driver=postgres # mysql, postgres, sqlite3

# The database connection properties individually or use the db.connection
db.host=localhost  # Use db.host /tmp/app.db is your driver is sqlite
db.user=dbuser
db.name=dbname
db.password=dbpassword

# Database connection string (host, user, dbname and other params)
db.connection=localhost port=8500 user=user dbname=mydb sslmode=disable password=ack
# If true then the database will be initialized on startup.
db.autoinit=true 
```
## Decelerations
A global `Db *DbGorp` object is created in `github.com/revel/modules/gorp/app`.
The `Db` is initialized from the app.conf if `db.autoinit=true`.
 ```go
// DB Gorp
type DbGorp struct {
	Gorp *gorp.DbMap
	// The Sql statement builder to use to build select statements
	SqlStatementBuilder sq.StatementBuilderType
}

var (
	// The database map to use to populate data
	Db = &DbGorp{}
)

```
## Usage
If `db.autoinit=true` in app.conf then you can add your tables to Gorp on app start.
Note that the tables are added as a function using `gorp.Db.SetDbInit` - this is for database thread pooling
```go
import (
	"github.com/revel/revel"
	"github.com/revel/modules/gorp/app"
)
func init() {
	revel.OnAppStart(func(){
		// Register tables
		gorp.Db.SetDbInit(func(dbGorp *gorp.DbGorp) error {
			// Register tables
			gorp.Db.Map.AddTableWithName(model.MyTable{}, "my_table")
			return nil
		})		
	},5)
}
```
### Controller
Controllers, with the `gorpController.Controller` embedded, 
have a gorp.DbGorp populated in the `Controller.Db`. This database is
the global one that is created on startup

```go
package controllers

import (
	"github.com/revel/revel"
	"github.com/revel/modules/orm/gorp/app/controllers"
)

type App struct {
	gorpController.Controller
}
type TableRow struct {
  Id int `db:"id,int64"`
}
func (c App) Index() revel.Result {
  sql,args,_ := c.Db.SqlStatementBuilder.Select("*").From("table").Limit(1).ToSql()
  row := &TableRow{}
  if err:= c.Db.Map.SelectOne(row,sql,args...); err!=nil {
    c.RenderError(err)
  }
  return c.Render(row)
}
```

### Multiple databases
The gorp module can populate a `DbGorp` object for you from a `gorp.DbInfo` object. So if you don't want
to use the global database (in the gorp module) you can initialize another anywhere in your project.
 ```go
import (
	"github.com/revel/revel"
	"github.com/revel/modules/orm/gorp/app"
)
var (
SecondDb = &gorp.DbGorp{} 
)
func init() {
	revel.OnAppStart(func(){
		// Create a DbInfo object with a minimum of a driver and other details
		params := gorp.DbInfo{Driver:"postgres",DbUser:revel.Config.StringDefault("seconddb.user", "default")}
		secondDb.Info = params
		if err:=secondDb.InitDb(true); err!=nil {
		  revel.Panicf("Second database failed to open %s", err.Error())
		}
	},0)
}
```

### Multi Channel Connections
This is not connection pooling - this is for distributing work across multiple channels to get a 
lot of stuff done fast. It creates a bunch workers and each worker has
its own connection (On the start of the worker a status is sent in case you want to do some prework). 
Tasks are sent through the `DbWorkContainer.InputChannel` which distributes the
task to whatever worker is available. 

If you are using any tables that requires GORP to have initialized tables you 
must register the tables using `gorp.Db.SetDbInit`. This is the only way that this service
can properly initialize the newly thread created GORP instances. Here is an example.  
```go
import (
	"github.com/revel/revel"
	"github.com/revel/modules/gorp/app"
)
func init() {
	revel.OnAppStart(func(){
		// Register tables
		gorp.Db.SetDbInit(func(dbGorp *gorp.DbGorp) error {
			// Register tables
			gorp.Db.Map.AddTableWithName(model.MyTable{}, "my_table")
			return nil
		})		
	},5)
}
```
 

In order to achieve this there is a `gorp.DbWorkerContainer` which is initialized by
`NewDbWorker(db *DbGorp, callBack DbCallback, numWorkers int) (container *DbWorkerContainer, err error)`
The `DbCallback` can be initialized by
`gorp.MakeCallback(status func(phase WorkerPhase, worker *DbWorker), work func(value interface{}, worker *DbWorker)) DbCallback`
()the status function is optional) or implemented by your structure

Once the `gorp.DbWorkerContainer` is created tasks can be submitted to it by using the 
`gorp.DbWorkerContainer.InputChannel<-task` this will call the `gorp.DbCallback.Work` function
passing in an instance of *DbWorker

```go
  // Assume sourceDb is a *gorp.DbGorp instance
	workerPool := gorp.NewDbWorker(d.Db, gorp.MakeCallback(func(phase gorp.WorkerPhase, worker *gorp.DbWorker) {
	    // On start initialize some data to be used later
		if phase == gorp.Start {
			dataList, err := model.FetchData(worker.Db)
			worker.SharedData["dataList"] = dataList

		} else if phase == gorp.JobLongrunning {
			revel.AppLog.Error("Long running process detected", "worker", worker.Id)
		}
	}, func(work interface{}, worker *gorp.DbWorker) {
		dataList := worker.SharedData["dataList"].(model.DataList)
        // Whatever is sent into the workerPool.InputChannel<- will be the value
        value := work.(*SomeObject)
        
	}), 100) 
	// Set the timeout for watchdog notifications
	workerPool.LongWorkTimeout=300
	err := workerPool.Start() // Start a 100 worker threads
	if err!=nil {
	  return err
	}
	defer workerPool.Close(0) // Close, wait for channels to exit (non zero would exit after timeout)
	// var tasks[] a large list of work to be done
	tasksBlock := make([]*SomeObject,100)
	for i,task := range tasks {
	  if i>0 && i%100==0 {
        workerPool.InputChannel<-taskBlock
        tasksBlock = make([]*SomeObject)  
	  }
	  taskBlock = append(taskBlock, task)
	}
	workerPool.InputChannel<-taskBlock
	// Pool is closed on defer function, it will not return till pool closes
	return nil
```
#### Implementation notes
If your "work" is short but there is a lot of it then it is highly recommended you pass in lists of 
items to work on. Channels are great at providing an easy way
to move data to and from threads but it is a process of synchronizing between two threads. If you pass
in a single object at a time you pay that cost on every row you pass. If you do it once every 10,000
rows then the cost is minimal.

#### WorkParallel function
There is a handy function called 
`WorkParallel(db *DbGorp, tasks []func(worker *DbWorker), returnResults bool, maxNumWorkers int, timeouts int) (err error)`
which makes it simple to do create a group of calls 
```go
  // Assume sourceDb is a *gorp.DbGorp instance
	task := func(query string) func(db *gorp.DbWorker) {
		return func(worker *gorp.DbWorker) {
			_, e := worker.Db.Map.Exec("ANALYZE " + query)
		}
	}
	gorp.WorkParallel(sourceDb, []func(worker *gorp.DbWorker){
		task("history"),
		task("summary"),
		task("daily"),
		task("monthly"),
	}, false, 0, 0)

```
#### Returning data from workers
`gorp.DbWorker` contains an OutputChannel which you can send the results back to, you must read
from the output the same number of times that you wrote. The output size is the same size as the
 
```go
  // Assume sourceDb is a *gorp.DbGorp instance
	task := func(query string) func(db *gorp.DbWorker) {
		return func(worker *gorp.DbWorker) {
			result, err := worker.Db.Map.Exec("ANALYZE " + query)
			worker.OutputChannel <- []interface{}{result,err}
		}
	}
	r,e:=gorp.WorkParallel(sourceDb, []func(worker *gorp.DbWorker){
		task("history"),
		task("daily"),
		task("monthly"),
	},true, 0, 0)

	println(r, e)
	for _,result := range r {

		key,err:= result.([]interface{})[0],result.([]interface{})[1]
		fmt.Println("Returned result",key,"error",err)
	}

```
#### Watchdog Timeouts
`gorp.DbWorkerContainer` contains a couple of timeout settings (in seconds) used to monitor
the startup and running of the workers for the duration of the container 
(this is typically called a watchdog timeout).
- `StartWorkTimeout` If greater then 0 this is the timeout in seconds that it takes to start a worker. 
- `LongWorkTimeout` If greater then 0 this is the timeout in seconds that it takes before a 
notification is sent to the `DbCallbackImplied.StatusFn func(phase WorkerPhase, worker *DbWorker)`
if a worker runs past X seconds on a single task. 
Each worker will have their own watchdog channel and it will send a `gorp.JobLongrunning` and
the `gorp.DBWorker` to the status function so you can log or investigate long running processes 
