package rdb

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Rdb struct {
	config *Config

	db          *sqlx.DB
	placeHolder string
}

func NewRdb(dbType Type, opts ...ConfigOptionFunc) (*Rdb, error) {
	config := newConfig()
	for _, opt := range opts {
		opt(config)
	}

	rdbInst := &Rdb{config: config}

	{
		db, dbErr := newDb(dbType, config)
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
		rdbInst.db.SetConnMaxLifetime(config.MaxLifeTime)
		rdbInst.db.SetMaxIdleConns(config.MaxIdleConns)
	}

	{
		pingErr := rdbInst.db.Ping()
		if pingErr != nil {
			return nil, pingErr
		}
	}

	return rdbInst, nil
}

func (r *Rdb) Engine() *sqlx.DB {
	return r.db
}

func (r *Rdb) Query(ctx context.Context, table []string, opts ...QueryOptionFunc) ([]map[string]any, error) {

	sqlOpt := newDefaultSqlOption()
	for _, opt := range opts {
		if opt != nil {
			opt(sqlOpt)
		}
	}
	sqlOpt.ctx = ctx
	sqlOpt.table = table

	instruct, instructErr := genQuerySql(table, sqlOpt.column, sqlOpt.where, sqlOpt.filter.offset, sqlOpt.filter.limit, WithPlaceholder(r.placeHolder))
	r.config.Logger.Debugf("id:%v,fn:query=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.config.Logger.Errorf("id:%v,fn:query=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	rows, rowsErr := r.db.QueryxContext(ctx, instruct.Sql, instruct.Args...)
	if rowsErr != nil {
		r.config.Logger.Errorf("id:%v,fn:query=>rowsErr:%v", ctx.Value("id"), rowsErr)
		return nil, rowsErr
	}

	var list []map[string]any
	for rows.Next() {
		item := make(map[string]any)
		scanErr := rows.MapScan(item)
		if scanErr != nil {
			r.config.Logger.Errorf("id:%v,fn:query=>scanErr:%v", ctx.Value("id"), scanErr)
			return nil, scanErr
		}
		list = append(list, item)
	}
	closeErr := rows.Close()
	if closeErr != nil {
		r.config.Logger.Errorf("id:%v,fn:query=>closeErr:%v", ctx.Value("id"), closeErr)
		return nil, closeErr
	}
	return list, nil
}

func (r *Rdb) transaction(ctx context.Context, instructs ...*Instruct) (map[string]int64, error) {
	r.config.Logger.Debugf("id:%v,fn:transaction=>instructs:%#v", ctx.Value("id"), instructsToString(instructs))

	tx, txErr := r.db.BeginTx(ctx, nil)
	if txErr != nil {
		r.config.Logger.Errorf("id:%v,fn:transaction=>txErr:%v", ctx.Value("id"), txErr)
		return nil, txErr
	}

	var affectedRows = make(map[string]int64)

	execErrs := make([]error, len(instructs))
	for _, instruct := range instructs {
		r.config.Logger.Debugf("about to do ExecContext for instruct:%v", instruct)
		execRst, execErr := r.db.ExecContext(ctx, instruct.Sql, instruct.Args...)
		if execErr == nil {
			r.config.Logger.Debugf("ExecContext successfully,about to get RowsAffected")
			affected, _ := execRst.RowsAffected()
			affectedRows[fmt.Sprintf("%v", instruct)] = affected
			r.config.Logger.Debugf("instruct exec successfully,affected:%v", affected)
			continue
		}
		r.config.Logger.Errorf("id:%v,fn:transaction=>execErr:%v", ctx.Value("id"), execErr)
		execErrs = append(execErrs, execErr)

		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			r.config.Logger.Errorf("id:%v,fn:transaction=>rollbackErr:%v", ctx.Value("id"), rollbackErr)
			return affectedRows, fmt.Errorf("execErr:%v,rollbackErr:%v", execErr, rollbackErr)
		}
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		r.config.Logger.Errorf("id:%v,fn:transaction=>commitErr:%v", ctx.Value("id"), commitErr)
		return affectedRows, fmt.Errorf("execErrs:%v,commitErr:%v", execErrs, commitErr)
	}

	return affectedRows, nil
}

func (r *Rdb) Delete(ctx context.Context, table string, where map[string]any, opts ...QueryOptionFunc) (map[string]int64, error) {
	if where == nil {
		r.config.Logger.Errorf("id:%v,fn:Delete=>nil where", ctx.Value("id"))
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
	r.config.Logger.Debugf("id:%v,fn:delete=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.config.Logger.Debugf("id:%v,fn:delete=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	return r.transaction(ctx, instruct)
}

func (r *Rdb) Insert(ctx context.Context, table string, data map[string]any, opts ...QueryOptionFunc) (map[string]int64, error) {
	if data == nil {
		r.config.Logger.Errorf("id:%v,fn:Insert=>nil data", ctx.Value("id"))
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
	r.config.Logger.Debugf("id:%v,fn:insert=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.config.Logger.Errorf("id:%v,fn:insert=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	return r.transaction(ctx, instruct)
}

func (r *Rdb) Update(ctx context.Context, table string, data map[string]any, where map[string]any, opts ...QueryOptionFunc) (map[string]int64, error) {
	if data == nil {
		r.config.Logger.Errorf("id:%v,fn:Update=>nil data", ctx.Value("id"))
		return nil, errWrapper(IllegalParams)
	}
	if where == nil {
		r.config.Logger.Errorf("id:%v,fn:Update=>nil where", ctx.Value("id"))
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
	r.config.Logger.Debugf("id:%v,fn:update=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.config.Logger.Debugf("id:%v,fn:update=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	return r.transaction(ctx, instruct)
}
