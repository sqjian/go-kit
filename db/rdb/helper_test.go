package rdb

import (
	"testing"
)

func Test_genQuerySql(t *testing.T) {
	type args struct {
		table  []string
		column []string
		where  map[string]interface{}
		group  []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				table:  []string{"table1"},
				column: []string{"column1"},
				where: map[string]interface{}{
					"whereKey1": "whereVal1",
				},
				group: []string{"group1"},
			},
		},
		{
			name: "test2",
			args: args{
				table:  []string{"table1"},
				column: []string{"*"},
				where:  nil,
				group:  []string{"group1"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRst, gotErr := genQuerySql(tt.args.table, tt.args.column, tt.args.where, tt.args.group, 0, 1)
			t.Logf("gotRst:%v, gotArgs:%v", gotRst, gotErr)
		})
	}
}

func Test_genInsertSql(t *testing.T) {
	type args struct {
		table string
		data  map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				table: "table1",
				data: map[string]interface{}{
					"dataKey1": "dataVal1",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRst, gotErr := genInsertSql(tt.args.table, tt.args.data)
			t.Logf("gotRst:%v, gotArgs:%v", gotRst, gotErr)
		})
	}
}
