package gorp

import (
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"    // mysql package
	_ "github.com/jinzhu/gorm/dialects/postgres" // postgres package
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // mysql package
	"github.com/revel/revel"
	sq "gopkg.in/Masterminds/squirrel.v1"
	"gopkg.in/gorp.v2"
	"github.com/revel/revel/logger"
)

var (
	// The database map to use to populate data
	Db = &DbGorp{}
	moduleLogger logger.MultiLogger
)
func init() {
	revel.RegisterModuleInit(func(module *revel.Module){
		moduleLogger = module.Log
		moduleLogger.Debug("Assigned Logger")
	})
}
func (dbResult *DbGorp)InitDb(open bool) (err error) {
	dbInfo := dbResult.Info

	switch dbInfo.DbDriver {
	default:
		dbResult.SqlStatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Question)
		dbInfo.Dialect = gorp.SqliteDialect{}
		if len(dbInfo.DbConnection) == 0 {
			dbInfo.DbConnection = fmt.Sprintf(dbInfo.DbHost)
		}
	case "postgres":
		dbResult.SqlStatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		dbInfo.Dialect = gorp.PostgresDialect{}
		if len(dbInfo.DbConnection) == 0 {
			dbInfo.DbConnection = fmt.Sprintf("host=%s port=8500 user=%s dbname=%s sslmode=disable password=%s", dbInfo.DbHost, dbInfo.DbUser, dbInfo.DbName, dbInfo.DbPassword)
		}
	case "mysql":
		dbResult.SqlStatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Question)
		dbInfo.Dialect = gorp.MySQLDialect{}
		if len(dbInfo.DbConnection) == 0 {
			dbInfo.DbConnection = fmt.Sprintf("%s:%s@%s/%s?charset=utf8&parseTime=True&loc=Local", dbInfo.DbUser, dbInfo.DbPassword, dbInfo.DbHost, dbInfo.DbName)
		}
	}

	if open {
		err = dbResult.OpenDb()
	}
	return
}

// Initialize the database from revel.Config
func InitDb(dbResult *DbGorp) (error) {
	params := DbInfo{}
	params.DbDriver = revel.Config.StringDefault("db.driver", "sqlite3")
	params.DbHost = revel.Config.StringDefault("db.host", "localhost")
	if params.DbDriver == "sqlite3" && params.DbHost == "localhost" {
		params.DbHost = "/tmp/app.db"
	}
	params.DbUser = revel.Config.StringDefault("db.user", "default")
	params.DbPassword = revel.Config.StringDefault("db.password", "")
	params.DbName = revel.Config.StringDefault("db.name", "default")
	params.DbConnection = revel.Config.StringDefault("db.connection", "")
	dbResult.Info = &params

	return dbResult.InitDb(true)
}
