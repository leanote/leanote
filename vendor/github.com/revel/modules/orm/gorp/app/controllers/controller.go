package gorpController

import (
	"database/sql"
	"github.com/revel/modules/orm/gorp/app"
	"github.com/revel/revel"
)

// Controller definition for database transaction
// This controller is only useful if you intend to use the database instance
// defined in github.com/revel/modules/orm/gorp/app.Db
type Controller struct {
	*revel.Controller
	Txn *gorp.Transaction
	Db  *gorp.DbGorp
}

// Begin a transaction
func (c *Controller) Begin() revel.Result {
	c.Db = gorp.Db
	txn, err := gorp.Db.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

// Rollback if it's still going (must have panicked).
func (c *Controller) Rollback() revel.Result {
	if c.Txn != nil {
		if err := c.Txn.Rollback(); err != nil {
			if err != sql.ErrTxDone {
				panic(err)
			}
		}
		c.Txn = nil
	}
	return nil
}

// Commit the transaction.
func (c *Controller) Commit() revel.Result {
	if c.Txn != nil {
		if err := c.Txn.Commit(); err != nil {
			if err != sql.ErrTxDone {
				panic(err)
			}
		}
		c.Txn = nil
	}
	return nil
}

func init() {
	// Run this as soon as possible
	revel.OnAppStart(func() {
		if revel.Config.BoolDefault("db.autoinit", false) {
			if err := gorp.InitDb(gorp.Db); err != nil {
				// Force a failure
				revel.RevelLog.Panicf("gorp:Unable to initialize database")
			}
			revel.InterceptMethod((*Controller).Begin, revel.BEFORE)
			revel.InterceptMethod((*Controller).Commit, revel.AFTER)
			revel.InterceptMethod((*Controller).Rollback, revel.FINALLY)
		}
	}, 0,
	)
}
