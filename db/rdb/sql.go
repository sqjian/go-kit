package rdb

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Rdb struct {
	meta *DbMeta

	db          *sqlx.DB
	placeHolder string
}

func NewRdb(dbType Type, opts ...MetaOption) (*Rdb, error) {
	meta := newMeta()
	for _, opt := range opts {
		opt.apply(meta)
	}

	rdbInst := &Rdb{meta: meta}

	{
		db, dbErr := newDb(dbType, meta)
		if dbErr != nil {
			return nil, dbErr
		}
		rdbInst.db = db
	}

	{
		placeHolder, placeHolderErr := newPlaceHolder(dbType)
		if placeHolderErr != nil {
			return nil, placeHolderErr
		}
		rdbInst.placeHolder = placeHolder
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

func (r *Rdb) Query(ctx context.Context, table []string, opts ...QueryOptionFunc) ([]map[string]interface{}, error) {

	sqlOpt := newDefaultSqlOption()
	for _, opt := range opts {
		if opt != nil {
			opt(sqlOpt)
		}
	}
	sqlOpt.ctx = ctx
	sqlOpt.table = table

	instruct, instructErr := genQuerySql(table, sqlOpt.column, sqlOpt.where, sqlOpt.filter.offset, sqlOpt.filter.limit, WithPlaceholder(r.placeHolder))
	r.meta.Logger.Debugf("id:%v,fn:query=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:query=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	rows, rowsErr := r.db.QueryxContext(ctx, instruct.Sql, instruct.Args...)
	if rowsErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:query=>rowsErr:%v", ctx.Value("id"), rowsErr)
		return nil, rowsErr
	}

	var list []map[string]interface{}
	for rows.Next() {
		item := make(map[string]interface{})
		scanErr := rows.MapScan(item)
		if scanErr != nil {
			r.meta.Logger.Errorf("id:%v,fn:query=>scanErr:%v", ctx.Value("id"), scanErr)
			return nil, scanErr
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

func (r *Rdb) transaction(ctx context.Context, instructs ...*Instruct) (map[string]int64, error) {
	r.meta.Logger.Debugf("id:%v,fn:transaction=>instructs:%#v", ctx.Value("id"), instructsToString(instructs))

	tx, txErr := r.db.BeginTx(ctx, nil)
	if txErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:transaction=>txErr:%v", ctx.Value("id"), txErr)
		return nil, txErr
	}

	var affectedRows = make(map[string]int64)

	execErrs := make([]error, len(instructs))
	for _, instruct := range instructs {
		r.meta.Logger.Debugf("about to do ExecContext for instruct:%v", instruct)
		execRst, execErr := r.db.ExecContext(ctx, instruct.Sql, instruct.Args...)
		if execErr == nil {
			r.meta.Logger.Debugf("ExecContext successfully,about to get RowsAffected")
			affected, _ := execRst.RowsAffected()
			affectedRows[fmt.Sprintf("%v", instruct)] = affected
			r.meta.Logger.Debugf("instruct exec successfully,affected:%v", affected)
			continue
		}
		r.meta.Logger.Errorf("id:%v,fn:transaction=>execErr:%v", ctx.Value("id"), execErr)
		execErrs = append(execErrs, execErr)

		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			r.meta.Logger.Errorf("id:%v,fn:transaction=>rollbackErr:%v", ctx.Value("id"), rollbackErr)
			return affectedRows, fmt.Errorf("execErr:%v,rollbackErr:%v", execErr, rollbackErr)
		}
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:transaction=>commitErr:%v", ctx.Value("id"), commitErr)
		return affectedRows, fmt.Errorf("execErrs:%v,commitErr:%v", execErrs, commitErr)
	}

	return affectedRows, nil
}

func (r *Rdb) Delete(ctx context.Context, table string, where map[string]interface{}, opts ...QueryOptionFunc) (map[string]int64, error) {
	if where == nil {
		r.meta.Logger.Errorf("id:%v,fn:Delete=>nil where", ctx.Value("id"))
		return nil, errWrapper(IllegalParams)
	}
	sqlOpt := newDefaultSqlOption()
	for _, opt := range opts {
		if opt != nil {
			opt(sqlOpt)
		}
	}
	sqlOpt.ctx = ctx
	sqlOpt.table = []string{table}

	instruct, instructErr := genDeleteSql(table, where, WithPlaceholder(r.placeHolder))
	r.meta.Logger.Debugf("id:%v,fn:delete=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.meta.Logger.Debugf("id:%v,fn:delete=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	return r.transaction(ctx, instruct)
}

func (r *Rdb) Insert(ctx context.Context, table string, data map[string]interface{}, opts ...QueryOptionFunc) (map[string]int64, error) {
	if data == nil {
		r.meta.Logger.Errorf("id:%v,fn:Insert=>nil data", ctx.Value("id"))
		return nil, errWrapper(IllegalParams)
	}

	sqlOpt := newDefaultSqlOption()
	for _, opt := range opts {
		if opt != nil {
			opt(sqlOpt)
		}
	}
	sqlOpt.ctx = ctx
	sqlOpt.table = []string{table}

	instruct, instructErr := genInsertSql(table, data, WithPlaceholder(r.placeHolder))
	r.meta.Logger.Debugf("id:%v,fn:insert=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:insert=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	return r.transaction(ctx, instruct)
}

func (r *Rdb) Update(ctx context.Context, table string, data map[string]interface{}, where map[string]interface{}, opts ...QueryOptionFunc) (map[string]int64, error) {
	if data == nil {
		r.meta.Logger.Errorf("id:%v,fn:Update=>nil data", ctx.Value("id"))
		return nil, errWrapper(IllegalParams)
	}
	if where == nil {
		r.meta.Logger.Errorf("id:%v,fn:Update=>nil where", ctx.Value("id"))
		return nil, errWrapper(IllegalParams)
	}
	sqlOpt := newDefaultSqlOption()
	for _, opt := range opts {
		if opt != nil {
			opt(sqlOpt)
		}
	}
	sqlOpt.ctx = ctx
	sqlOpt.table = []string{table}

	instruct, instructErr := genUpdateSql(table, data, where, WithPlaceholder(r.placeHolder))
	r.meta.Logger.Debugf("id:%v,fn:update=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.meta.Logger.Debugf("id:%v,fn:update=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	return r.transaction(ctx, instruct)
}
