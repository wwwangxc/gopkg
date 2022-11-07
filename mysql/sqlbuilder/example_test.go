package sqlbuilder_test

import (
	"fmt"

	"github.com/wwwangxc/gopkg/mysql/sqlbuilder"
)

func Example() {
	builder := sqlbuilder.NewSelect("table_name",
		sqlbuilder.WithFieldEqual("f_equal", 1),
		sqlbuilder.WithFieldLike("f_like", "val"),
		sqlbuilder.WithFieldLessThan("f_less", 111),
		sqlbuilder.WithFieldLessOrEqualThan("f_less_eq", 222),
		sqlbuilder.WithFieldGreaterThan("f_greater", 333),
		sqlbuilder.WithFieldGreaterOrEqualThan("f_greater_equal", 444),
		sqlbuilder.WithOrderBy("f_asc"),
		sqlbuilder.WithOrderByDESC("f_desc"),
		sqlbuilder.WithLimit(100),
		sqlbuilder.WithOffset(1),
		sqlbuilder.WithForceIndex("index_name"),
	)

	// sql:
	//     SELECT select_field FROM table_name
	//     FORCE INDEX(index_name)
	//     WHERE f_equal=? AND f_like LIKE ? AND f_less<? AND f_less_equal<=? AND f_greater>? AND f_greater_equal>=?
	//     ORDER BY f_asc ASC, f_desc DESC
	//     LIMIT 100 OFFSET 1
	//
	// args:
	//     [1, "%val%", 111, 222, 333, 444]
	sql, args := builder.Build("select_field")
	fmt.Printf("generate sql: %s\n", sql)
	fmt.Println(args)
}
