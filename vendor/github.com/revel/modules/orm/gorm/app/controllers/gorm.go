package gormcontroller

import (
	"database/sql"
	"fmt"

	"github.com/jinzhu/gorm"
	gormdb "github.com/revel/modules/orm/gorm/app"
	"github.com/revel/revel"
)

// Controller is a Revel controller with a pointer to the opened database
type Controller struct {
	*revel.Controller
    DB *gorm.DB
}

func (c *Controller) setDB() revel.Result {
	c.DB = gormdb.DB
	return nil
}

// TxnController is a Revel controller with database transaction support (begin, commit and rollback)
type TxnController struct {
	*revel.Controller
	Txn *gorm.DB
}

// Begin begins a DB transaction
func (c *TxnController) Begin() revel.Result {

	txn := gormdb.DB.Begin()
	if txn.Error != nil {
		c.Log.Panic("Transaction begine error","error",txn.Error)
	}

	c.Txn = txn
	return nil
}

// Commit commits the database transation
func (c *TxnController) Commit() revel.Result {
	if c.Txn == nil {
		return nil
	}

	c.Txn.Commit()
	if c.Txn.Error != nil && c.Txn.Error != sql.ErrTxDone {
		fmt.Println(c.Txn.Error)
		panic(c.Txn.Error)
	}

	c.Txn = nil
	return nil
}

// Rollback rolls back the transaction (eg. after a panic)
func (c *TxnController) Rollback() revel.Result {
	if c.Txn == nil {
		return nil
	}

	c.Txn.Rollback()
	if c.Txn.Error != nil && c.Txn.Error != sql.ErrTxDone {
		fmt.Println(c.Txn.Error)
		panic(c.Txn.Error)
	}

	c.Txn = nil
	return nil
}

func init() {
	revel.OnAppStart(func() {
		if revel.Config.BoolDefault("db.autoinit", true) {
			gormdb.InitDB()
			revel.InterceptMethod((*TxnController).Begin, revel.BEFORE)
			revel.InterceptMethod((*TxnController).Commit, revel.AFTER)
			revel.InterceptMethod((*TxnController).Rollback, revel.FINALLY)

			revel.InterceptMethod((*Controller).setDB, revel.BEFORE)
		}
	})
}
