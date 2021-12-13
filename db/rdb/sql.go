package rdb

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

type rdb struct {
	meta *Meta

	db *sql.DB
}

func NewRdb(dbType Type, opts ...Option) (*rdb, error) {
	meta := newDefaultMeta()
	for _, opt := range opts {
		opt.apply(meta)
	}

	rdbInst := &rdb{meta: meta}

	{
		db, dbErr := newDb(dbType, meta)
		if dbErr != nil {
			return nil, dbErr
		}
		rdbInst.db = db
	}

	{
		rdbInst.db.SetConnMaxLifetime(meta.MaxLifeTime)
		rdbInst.db.SetMaxIdleConns(meta.MaxIdleConns)
	}

	{
		pingErr := rdbInst.db.Ping()
		if pingErr != nil {
			return nil, pingErr
		}
	}

	return rdbInst, nil
}

func (r *rdb) columns(ctx context.Context, table string) ([]string, error) {
	rawSql := fmt.Sprintf("SELECT * FROM %v WHERE 1 = 0", table)
	r.meta.Logger.Debugf("id:%v,fn:columns=>rawSql:%v", ctx.Value("id"), rawSql)

	rows, err := r.db.QueryContext(ctx, rawSql)
	if err != nil {
		r.meta.Logger.Errorf("id:%v,fn:columns=>err:%v", ctx.Value("id"), err)
		return nil, err
	}
	return rows.Columns()
}

func (r *rdb) Query(ctx context.Context, table string, where map[string]interface{}) ([]map[string]interface{}, error) {

	if where == nil {
		return nil, fmt.Errorf("nil where")
	}

	columns, columnsErr := r.columns(ctx, table)
	if columnsErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:query=>columnsErr:%v", ctx.Value("id"), columnsErr)
		return nil, columnsErr
	}
	for expectColumn, _ := range where {
		matched := false
		for _, realColumn := range columns {
			if expectColumn == realColumn {
				matched = true
				break
			}
		}
		if !matched {
			return nil, fmt.Errorf("can't find column:%v in table columns:%v", expectColumn, columns)
		}
	}

	rawSql := func() string {
		s := fmt.Sprintf("SELECT * FROM %v ", table)
		var w []string
		for k, v := range where {
			w = append(w, fmt.Sprintf("%v = %#v", k, v))
		}
		s = s + "WHERE " + strings.Join(w, " AND ")
		return s
	}()
	r.meta.Logger.Debugf("id:%v,fn:query=>rawSql:%v", ctx.Value("id"), rawSql)

	rows, rowsErr := r.db.Query(rawSql)
	if rowsErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:query=>rowsErr:%v", ctx.Value("id"), rowsErr)
		return nil, rowsErr
	}

	columnLength := len(columns)
	cache := make([]interface{}, columnLength)
	for index, _ := range cache {
		var a interface{}
		cache[index] = &a
	}

	var list []map[string]interface{}
	for rows.Next() {
		scanErr := rows.Scan(cache...)
		if scanErr != nil {
			r.meta.Logger.Errorf("id:%v,fn:query=>scanErr:%v", ctx.Value("id"), scanErr)
			return nil, scanErr
		}
		item := make(map[string]interface{})
		for i, data := range cache {
			item[columns[i]] = *data.(*interface{})
		}
		list = append(list, item)
	}
	closeErr := rows.Close()
	if closeErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:query=>closeErr:%v", ctx.Value("id"), closeErr)
		return nil, closeErr
	}
	return list, nil
}

func (r *rdb) transaction(ctx context.Context, query string, args ...interface{}) (int64, error) {
	r.meta.Logger.Debugf("id:%v,fn:transaction=>query:%v,args:%v", ctx.Value("id"), query, args)

	tx, txErr := r.db.Begin()
	if txErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:transaction=>txErr:%v", ctx.Value("id"), txErr)
		return 0, txErr
	}
	defer func(tx *sql.Tx) {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			r.meta.Logger.Errorf("id:%v,fn:transaction=>rollbackErr:%v", ctx.Value("id"), rollbackErr)
		}
	}(tx)

	stmt, stmtErr := tx.PrepareContext(ctx, query)
	if stmtErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:transaction=>stmtErr:%v", ctx.Value("id"), stmtErr)
		return 0, stmtErr
	}
	defer func(stmt *sql.Stmt) {
		stmtCloseErr := stmt.Close()
		if stmtCloseErr != nil {
			r.meta.Logger.Errorf("id:%v,fn:transaction=>stmtCloseErr:%v", ctx.Value("id"), stmtCloseErr)
		}
	}(stmt)

	execRst, execErr := stmt.ExecContext(ctx)
	if execErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:transaction=>execErr:%v", ctx.Value("id"), execErr)
		return 0, execErr
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:transaction=>commitErr:%v", ctx.Value("id"), commitErr)
		return 0, commitErr
	}

	return execRst.RowsAffected()
}
func (r *rdb) Delete(ctx context.Context, table string, where map[string]interface{}) (int64, error) {
	if where == nil {
		r.meta.Logger.Errorf("id:%v,fn:Delete=>nil where", ctx.Value("id"))
		return 0, errWrapper(IllegalParams)
	}

	rawSql := func() string {
		s := fmt.Sprintf("DELETE FROM %v WHERE ", table)
		var w []string
		for k, v := range where {
			w = append(w, fmt.Sprintf("%v = %#v", k, v))
		}
		s += strings.Join(w, " AND ")
		return s
	}()

	return r.transaction(ctx, rawSql)
}

func (r *rdb) Insert(ctx context.Context, table string, data map[string]interface{}) (int64, error) {
	if data == nil {
		r.meta.Logger.Errorf("id:%v,fn:Insert=>nil data", ctx.Value("id"))
		return 0, errWrapper(IllegalParams)
	}

	rawSql := func() string {
		s := fmt.Sprintf("INSERT INTO %v ", table)
		var ks, vs []string
		for k, v := range data {
			ks = append(ks, fmt.Sprintf("`%v`", k))
			vs = append(vs, fmt.Sprintf("%v", v))
		}
		s += fmt.Sprintf("(%v) ", strings.Join(ks, ","))
		s += fmt.Sprintf("VALUES (%v) ", strings.Join(vs, ","))
		return s
	}()

	return r.transaction(ctx, rawSql)
}

func (r *rdb) Update(ctx context.Context, table string, data map[string]interface{}, where map[string]interface{}) (int64, error) {
	if data == nil {
		r.meta.Logger.Errorf("id:%v,fn:Update=>nil data", ctx.Value("id"))
		return 0, errWrapper(IllegalParams)
	}
	if where == nil {
		r.meta.Logger.Errorf("id:%v,fn:Update=>nil where", ctx.Value("id"))
		return 0, errWrapper(IllegalParams)
	}

	rawSql := func() string {
		s := fmt.Sprintf("UPDATE %v ", table)

		var dataKvs []string
		for dataKey, dataVal := range data {
			dataKvs = append(dataKvs, fmt.Sprintf("%v=%#v", dataKey, dataVal))
		}
		s = fmt.Sprintf("%v SET %v ", s, strings.Join(dataKvs, ","))

		var whereKvs []string
		for whereKey, whereVal := range where {
			whereKvs = append(whereKvs, fmt.Sprintf("%v=%#v", whereKey, whereVal))
		}
		s = fmt.Sprintf("%v WHERE %v ", s, strings.Join(whereKvs, ","))
		return s
	}()

	return r.transaction(ctx, rawSql)
}
