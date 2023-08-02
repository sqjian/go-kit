package rdb

import "context"

func newDefaultSqlOption() *SqlOption {
	sqlOption := &SqlOption{}

	sqlOption.filter.offset = 0
	sqlOption.filter.limit = 1000

	sqlOption.column = []string{}
	sqlOption.where = make(map[string]any)

	return sqlOption
}

type SqlOption struct {
	ctx   context.Context
	table []string

	filter struct {
		offset uint64
		limit  uint64
	}
	column []string
	where  map[string]any
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

func WithWhere(where map[string]any) QueryOptionFunc {
	return QueryOptionFunc(func(s *SqlOption) {
		s.where = where
	})
}
