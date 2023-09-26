package rds

import (
	"fmt"
	"github.com/elgris/sqrl"
	"strings"
)

type InstructOptionFunc func(*instructOpts)

func newDefaultInstructOpts() *instructOpts {
	return &instructOpts{
		placeholder: sqrl.Question,
	}
}

type instructOpts struct {
	placeholder sqrl.PlaceholderFormat
}

func WithPlaceholder(placeholder string) InstructOptionFunc {
	return func(i *instructOpts) {
		switch placeholder {
		case "?":
			i.placeholder = sqrl.Question
		case "$":
			i.placeholder = sqrl.Dollar
		}
	}
}

type Instruct struct {
	Sql  string
	Args []any
}

func NewInstruct(sql sqrl.Sqlizer, opts ...InstructOptionFunc) (*Instruct, error) {

	instructOptsInst := newDefaultInstructOpts()
	for _, opt := range opts {
		opt(instructOptsInst)
	}

	sqlStr, sqlArgs, sqlErr := sql.ToSql()
	if sqlErr != nil {
		return nil, sqlErr
	}

	var sqlPlaceHolderErr error
	sqlStr, sqlPlaceHolderErr = instructOptsInst.placeholder.ReplacePlaceholders(sqlStr)
	if sqlPlaceHolderErr != nil {
		return nil, sqlPlaceHolderErr
	}

	return &Instruct{sqlStr, sqlArgs}, nil
}

func InstructsToString(instructs []*Instruct) string {
	var rst []string
	for _, instruct := range instructs {
		rst = append(rst, fmt.Sprintf("%v", *instruct))
	}
	return strings.Join(rst, ",")
}

func GenQuerySql(table []string, column []string, where map[string]any, offset, limit uint64, opts ...InstructOptionFunc) (instruct *Instruct, err error) {

	var sql *sqrl.SelectBuilder

	if len(column) == 0 {
		column = append(column, "*")
	}
	sql = sqrl.Select(column...)

	if len(table) == 0 {
		return nil, fmt.Errorf("empty table")
	}
	sql = sql.From(table...)

	if len(where) != 0 {
		sql = sql.Where(sqrl.Eq(where))
	}

	sql = sql.Offset(offset)

	if limit != 0 {
		sql = sql.Limit(limit)
	}

	return NewInstruct(sql, opts...)
}

func GenInsertSql(table string, data map[string]any, opts ...InstructOptionFunc) (instruct *Instruct, err error) {

	var sql *sqrl.InsertBuilder

	sql = sqrl.Insert(table)

	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	var columns []string
	var values []any
	for column, value := range data {
		columns = append(columns, column)
		values = append(values, value)
	}
	sql = sql.Columns(columns...)
	sql = sql.Values(values...)

	return NewInstruct(sql, opts...)
}

func GenDeleteSql(table string, where map[string]any, opts ...InstructOptionFunc) (instruct *Instruct, err error) {

	var sql *sqrl.DeleteBuilder

	sql = sqrl.Delete(table)

	if len(where) == 0 {
		return nil, fmt.Errorf("empty where")
	}
	sql = sql.Where(where)

	return NewInstruct(sql, opts...)
}

func GenUpdateSql(table string, data map[string]any, where map[string]any, opts ...InstructOptionFunc) (instruct *Instruct, err error) {

	var sql *sqrl.UpdateBuilder

	sql = sqrl.Update(table)

	if len(where) == 0 {
		return nil, fmt.Errorf("empty where")
	}
	sql = sql.Where(where)

	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}
	sql = sql.SetMap(data)

	return NewInstruct(sql, opts...)
}
