package rdb

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Rdb struct {
	meta *DbMeta

	db *sqlx.DB
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

func (r *Rdb) columns(opts *SqlOption) ([]string, error) {
	rawSql := fmt.Sprintf("SELECT * FROM %v WHERE 1 = 0", opts.table)
	r.meta.Logger.Debugf("id:%v,fn:columns=>rawSql:%v", opts.ctx.Value("id"), rawSql)

	rows, err := r.db.QueryContext(opts.ctx, rawSql)
	if err != nil {
		r.meta.Logger.Errorf("id:%v,fn:columns=>err:%v", opts.ctx.Value("id"), err)
		return nil, err
	}
	return rows.Columns()
}

func (r *Rdb) checkWhere(opts *SqlOption) error {
	columns, columnsErr := r.columns(opts)
	if columnsErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:query=>columnsErr:%v", opts.ctx.Value("id"), columnsErr)
		return columnsErr
	}
	for expectColumn, _ := range opts.where {
		matched := false
		for _, realColumn := range columns {
			if expectColumn == realColumn {
				matched = true
				break
			}
		}
		if !matched {
			return fmt.Errorf("can't find column:%v in table columns:%v", expectColumn, columns)
		}
	}
	return nil
}

func (r *Rdb) Query(ctx context.Context, table []string, opts ...QueryOptionFunc) ([]map[string]interface{}, error) {

	sqlOpt := newDefaultSqlOption()
	{
		for _, opt := range opts {
			if opt != nil {
				opt(sqlOpt)
			}
		}
		sqlOpt.ctx = ctx
		sqlOpt.table = table

		if len(sqlOpt.where) != 0 {
			err := r.checkWhere(sqlOpt)
			if err != nil {
				return nil, err
			}
		}
	}

	instruct, instructErr := genQuerySql(table, sqlOpt.column, sqlOpt.where, sqlOpt.group, sqlOpt.limit.start, sqlOpt.limit.end)
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

	for _, instruct := range instructs {
		execRst, execErr := r.db.ExecContext(ctx, instruct.Sql, instruct.Args...)
		if execErr == nil {
			affected, _ := execRst.RowsAffected()
			affectedRows[fmt.Sprintf("%v", instruct)] = affected
			continue
		}
		r.meta.Logger.Errorf("id:%v,fn:transaction=>execErr:%v", ctx.Value("id"), execErr)

		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			r.meta.Logger.Errorf("id:%v,fn:transaction=>rollbackErr:%v", ctx.Value("id"), rollbackErr)
			return affectedRows, rollbackErr
		}
	}

	commitErr := tx.Commit()
	if commitErr != nil {
		r.meta.Logger.Errorf("id:%v,fn:transaction=>commitErr:%v", ctx.Value("id"), commitErr)
		return affectedRows, commitErr
	}

	return affectedRows, nil
}

func (r *Rdb) Delete(ctx context.Context, table string, where map[string]interface{}, opts ...QueryOptionFunc) (map[string]int64, error) {
	if where == nil {
		r.meta.Logger.Errorf("id:%v,fn:Delete=>nil where", ctx.Value("id"))
		return nil, errWrapper(IllegalParams)
	}
	sqlOpt := newDefaultSqlOption()
	{
		for _, opt := range opts {
			if opt != nil {
				opt(sqlOpt)
			}
		}
		sqlOpt.ctx = ctx
		sqlOpt.table = []string{table}

		if len(sqlOpt.where) != 0 {
			err := r.checkWhere(sqlOpt)
			if err != nil {
				return nil, err
			}
		}
	}

	instruct, instructErr := genDeleteSql(table, where)
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
	{
		for _, opt := range opts {
			if opt != nil {
				opt(sqlOpt)
			}
		}
		sqlOpt.ctx = ctx
		sqlOpt.table = []string{table}

		if len(sqlOpt.where) != 0 {
			err := r.checkWhere(sqlOpt)
			if err != nil {
				return nil, err
			}
		}
	}

	instruct, instructErr := genInsertSql(table, data)
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
	{
		for _, opt := range opts {
			if opt != nil {
				opt(sqlOpt)
			}
		}
		sqlOpt.ctx = ctx
		sqlOpt.table = []string{table}

		if len(sqlOpt.where) != 0 {
			err := r.checkWhere(sqlOpt)
			if err != nil {
				return nil, err
			}
		}
	}

	instruct, instructErr := genUpdateSql(table, data, where)
	r.meta.Logger.Debugf("id:%v,fn:update=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
	if instructErr != nil {
		r.meta.Logger.Debugf("id:%v,fn:update=>instruct:%v,instructErr:%v", ctx.Value("id"), instruct, instructErr)
		return nil, instructErr
	}

	return r.transaction(ctx, instruct)
}
