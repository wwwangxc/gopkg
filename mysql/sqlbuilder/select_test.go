package sqlbuilder

import (
	"reflect"
	"testing"
)

func TestNewSelect(t *testing.T) {
	type args struct {
		table string
		opts  []SelectOption
	}
	tests := []struct {
		name string
		args args
		want *SelectBuilder
	}{
		{
			name: "option empty",
			args: args{
				table: "table_name",
			},
			want: &SelectBuilder{
				table: "table_name",
			},
		},
		{
			name: "normal",
			args: args{
				table: "table_name",
				opts: []SelectOption{
					WithFieldEqual("f_equal", 1),
					WithFieldLike("f_like", "field"),
					WithFieldLessThan("f_less", 111),
					WithFieldLessOrEqualThan("f_less_equal", 222),
					WithFieldGreaterThan("f_greater", 333),
					WithFieldGreaterOrEqualThan("f_greater_equal", 444),
					WithOrderBy("f_asc"),
					WithOrderByDESC("f_desc"),
					WithLimit(555),
					WithOffset(666),
					WithForceIndex("index"),
				},
			},
			want: &SelectBuilder{
				table: "table_name",
				whereConditions: []string{
					"f_equal=?",
					"f_like LIKE ?",
					"f_less<?",
					"f_less_equal<=?",
					"f_greater>?",
					"f_greater_equal>=?",
				},
				whereArgs: []interface{}{
					1,
					"%field%",
					111,
					222,
					333,
					444,
				},
				orderByConditions: []string{
					"f_asc ASC",
					"f_desc DESC",
				},
				limit:      555,
				offset:     666,
				forceIndex: "index",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSelect(tt.args.table, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSelect() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSelectBuilder_Build(t *testing.T) {
	type args struct {
		table string
		opts  []SelectOption
	}
	type args1 struct {
		fields []string
	}
	tests := []struct {
		name  string
		args  args
		args1 args1
		want  string
		want1 []interface{}
	}{
		{
			name: "all condition empty",
			args: args{
				table: "table_name",
			},
			args1: args1{},
			want:  "SELECT * FROM table_name OFFSET 0",
			want1: []interface{}{},
		},
		{
			name: "select fields",
			args: args{
				table: "table_name",
			},
			args1: args1{
				fields: []string{"f1", "f2"},
			},
			want:  "SELECT f1, f2 FROM table_name OFFSET 0",
			want1: []interface{}{},
		},
		{
			name: "where",
			args: args{
				table: "table_name",
				opts: []SelectOption{
					WithFieldEqual("f_equal", 1),
					WithFieldLike("f_like", "field"),
					WithFieldLessThan("f_less", 111),
					WithFieldLessOrEqualThan("f_less_equal", 222),
					WithFieldGreaterThan("f_greater", 333),
					WithFieldGreaterOrEqualThan("f_greater_equal", 444),
				},
			},
			args1: args1{
				fields: []string{"f1", "f2"},
			},
			want: "SELECT f1, f2 FROM table_name " +
				"WHERE f_equal=? AND f_like LIKE ? AND f_less<? AND f_less_equal<=? AND f_greater>? AND f_greater_equal>=? OFFSET 0",
			want1: []interface{}{1, "%field%", 111, 222, 333, 444},
		},
		{
			name: "order by",
			args: args{
				table: "table_name",
				opts: []SelectOption{
					WithFieldEqual("f_equal", 1),
					WithFieldLike("f_like", "field"),
					WithFieldLessThan("f_less", 111),
					WithFieldLessOrEqualThan("f_less_equal", 222),
					WithFieldGreaterThan("f_greater", 333),
					WithFieldGreaterOrEqualThan("f_greater_equal", 444),
					WithOrderBy("f_asc"),
					WithOrderByDESC("f_desc"),
				},
			},
			args1: args1{
				fields: []string{"f1", "f2"},
			},
			want: "SELECT f1, f2 FROM table_name " +
				"WHERE f_equal=? AND f_like LIKE ? AND f_less<? AND f_less_equal<=? AND f_greater>? AND f_greater_equal>=? " +
				"ORDER BY f_asc ASC, f_desc DESC " +
				"OFFSET 0",
			want1: []interface{}{1, "%field%", 111, 222, 333, 444},
		},
		{
			name: "limit",
			args: args{
				table: "table_name",
				opts: []SelectOption{
					WithFieldEqual("f_equal", 1),
					WithFieldLike("f_like", "field"),
					WithFieldLessThan("f_less", 111),
					WithFieldLessOrEqualThan("f_less_equal", 222),
					WithFieldGreaterThan("f_greater", 333),
					WithFieldGreaterOrEqualThan("f_greater_equal", 444),
					WithOrderBy("f_asc"),
					WithOrderByDESC("f_desc"),
					WithLimit(555),
				},
			},
			args1: args1{
				fields: []string{"f1", "f2"},
			},
			want: "SELECT f1, f2 FROM table_name " +
				"WHERE f_equal=? AND f_like LIKE ? AND f_less<? AND f_less_equal<=? AND f_greater>? AND f_greater_equal>=? " +
				"ORDER BY f_asc ASC, f_desc DESC " +
				"LIMIT 555 " +
				"OFFSET 0",
			want1: []interface{}{1, "%field%", 111, 222, 333, 444},
		},
		{
			name: "offset",
			args: args{
				table: "table_name",
				opts: []SelectOption{
					WithFieldEqual("f_equal", 1),
					WithFieldLike("f_like", "field"),
					WithFieldLessThan("f_less", 111),
					WithFieldLessOrEqualThan("f_less_equal", 222),
					WithFieldGreaterThan("f_greater", 333),
					WithFieldGreaterOrEqualThan("f_greater_equal", 444),
					WithOrderBy("f_asc"),
					WithOrderByDESC("f_desc"),
					WithLimit(555),
					WithOffset(666),
				},
			},
			args1: args1{
				fields: []string{"f1", "f2"},
			},
			want: "SELECT f1, f2 FROM table_name " +
				"WHERE f_equal=? AND f_like LIKE ? AND f_less<? AND f_less_equal<=? AND f_greater>? AND f_greater_equal>=? " +
				"ORDER BY f_asc ASC, f_desc DESC " +
				"LIMIT 555 " +
				"OFFSET 666",
			want1: []interface{}{1, "%field%", 111, 222, 333, 444},
		},
		{
			name: "full condition",
			args: args{
				table: "table_name",
				opts: []SelectOption{
					WithFieldEqual("f_equal", 1),
					WithFieldLike("f_like", "field"),
					WithFieldLessThan("f_less", 111),
					WithFieldLessOrEqualThan("f_less_equal", 222),
					WithFieldGreaterThan("f_greater", 333),
					WithFieldGreaterOrEqualThan("f_greater_equal", 444),
					WithOrderBy("f_asc"),
					WithOrderByDESC("f_desc"),
					WithLimit(555),
					WithOffset(666),
					WithForceIndex("index"),
				},
			},
			args1: args1{
				fields: []string{"f1", "f2"},
			},
			want: "SELECT f1, f2 FROM table_name " +
				"FORCE INDEX(index) " +
				"WHERE f_equal=? AND f_like LIKE ? AND f_less<? AND f_less_equal<=? AND f_greater>? AND f_greater_equal>=? " +
				"ORDER BY f_asc ASC, f_desc DESC " +
				"LIMIT 555 " +
				"OFFSET 666",
			want1: []interface{}{1, "%field%", 111, 222, 333, 444},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewSelect(tt.args.table, tt.args.opts...)
			got, got1 := builder.Build(tt.args1.fields...)
			if got != tt.want {
				t.Errorf("SelectBuilder.Build() got = %v, want %v", got, tt.want)
			}
			if len(got1) == 0 && len(tt.want1) == 0 {
				return
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("SelectBuilder.Build() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
