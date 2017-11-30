package gorp
import (
	gorpa "gopkg.in/gorp.v2"
	"database/sql"
	sq "gopkg.in/Masterminds/squirrel.v1"
)
type (
	// This is a small wrapped around gorp.Transaction so you can make use of the builder statements as well
	Transaction struct {
		Map *gorpa.Transaction
	}
)

func (txn *Transaction) Rollback() (err error) {
	return txn.Map.Rollback()
}
func (txn *Transaction) Commit() (err error) {
	return txn.Map.Commit()
}
func (txn *Transaction) Select(i interface{}, builder sq.SelectBuilder) (l []interface{}, err error) {
	query, args, err := builder.ToSql()
	if err == nil {
		list, err := txn.Map.Select(i, query, args...)
		if err != nil && gorpa.NonFatalError(err) {
			return list, nil
		}
		if err==sql.ErrNoRows {
			err = nil
		}
		return list, err
	}
	return
}

func (txn *Transaction) SelectOne(i interface{}, builder sq.SelectBuilder) (err error) {
	query, args, err := builder.ToSql()
	if err == nil {
		err = txn.Map.SelectOne(i, query, args...)
		if err != nil && gorpa.NonFatalError(err) {
			return nil
		}
	}
	return
}

func (txn *Transaction) SelectInt(builder sq.SelectBuilder) (i int64, err error) {
	query, args, err := builder.ToSql()
	if err == nil {
		i, err = txn.Map.SelectInt(query, args...)
	}
	return
}

func (txn *Transaction) ExecUpdate(builder sq.UpdateBuilder) (r sql.Result, err error) {
	query, args, err := builder.ToSql()
	if err == nil {
		r, err = txn.Map.Exec(query, args...)
	}
	return
}
func (txn *Transaction) ExecInsert(builder sq.InsertBuilder) (r sql.Result, err error) {
	query, args, err := builder.ToSql()
	if err == nil {
		r, err = txn.Map.Exec(query, args...)
	}
	return
}

// Shifted some common functions up a level

func (txn *Transaction) Insert(list ...interface{}) error {
	return txn.Map.Insert(list...)
}
func (txn *Transaction) Update(list ...interface{}) (int64, error) {
	return txn.Map.Update(list...)
}
func (txn *Transaction) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return txn.Map.Get(i,keys...)
}
func (txn *Transaction) Delete(i ...interface{}) (int64, error) {
	return txn.Map.Delete(i...)
}
