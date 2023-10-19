package rds

import (
	"testing"
)

func Test_genQuerySql(t *testing.T) {
	type args struct {
		table  []string
		column []string
		where  map[string]any
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
				where: map[string]any{
					"whereKey1": "whereVal1",
				},
			},
		},
		{
			name: "test2",
			args: args{
				table:  []string{"table1"},
				column: []string{"*"},
				where:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRst, gotErr := genQuerySql(tt.args.table, tt.args.column, tt.args.where, 0, 1)
			t.Logf("gotRst:%v, gotArgs:%v", gotRst, gotErr)
		})
	}
}

func Test_genInsertSql(t *testing.T) {
	type args struct {
		table string
		data  map[string]any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				table: "table1",
				data: map[string]any{
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
