package rdb

import "context"

func newDefaultSqlOption() *SqlOption {
	sqlOption := &SqlOption{}

	sqlOption.limit.start = 0
	sqlOption.limit.end = 1000

	sqlOption.column = []string{}
	sqlOption.where = make(map[string]interface{})
	sqlOption.group = []string{}

	return sqlOption
}

type SqlOption struct {
	ctx   context.Context
	table []string

	limit struct {
		start uint64
		end   uint64
	}
	column []string
	where  map[string]interface{}
	group  []string
}

type QueryOption interface {
	apply(*SqlOption)
}

type QueryOptionFunc func(*SqlOption)

func (q QueryOptionFunc) apply(sqlOption *SqlOption) {
	q(sqlOption)
}

func WithColumn(column []string) QueryOptionFunc {
	return QueryOptionFunc(func(s *SqlOption) {
		s.column = column
	})
}

func WithWhere(where map[string]interface{}) QueryOptionFunc {
	return QueryOptionFunc(func(s *SqlOption) {
		s.where = where
	})
}

func WithGroup(group []string) QueryOptionFunc {
	return QueryOptionFunc(func(s *SqlOption) {
		s.group = group
	})
}
