package sqlbuilder

import "fmt"

// SelectOption select statement builder option
type SelectOption func(builder *SelectBuilder)

// WithFieldEqual append where condition
//
// WHERE ${field} = ?
func WithFieldEqual(field string, arg interface{}) SelectOption {
	return func(builder *SelectBuilder) {
		builder.appendWhereCondition(fmt.Sprintf("%s=?", field), arg)
	}
}

// WithFieldLike append where condition
//
// WHERE ${field} LIKE '%?%'
func WithFieldLike(field string, arg interface{}) SelectOption {
	return func(builder *SelectBuilder) {
		builder.appendWhereCondition(fmt.Sprintf("%s LIKE ?", field), fmt.Sprintf("%%%s%%", arg))
	}
}

// WithFieldLessThan append where condition
//
// WHERE ${field} < ?
func WithFieldLessThan(field string, arg interface{}) SelectOption {
	return func(builder *SelectBuilder) {
		builder.appendWhereCondition(fmt.Sprintf("%s<?", field), arg)
	}
}

// WithFieldLessOrEqualThan append where condition
//
// WHERE ${field} <= ?
func WithFieldLessOrEqualThan(field string, arg interface{}) SelectOption {
	return func(builder *SelectBuilder) {
		builder.appendWhereCondition(fmt.Sprintf("%s<=?", field), arg)
	}
}

// WithFieldGreaterThan append where condition
//
// WHERE ${field} > ?
func WithFieldGreaterThan(field string, arg interface{}) SelectOption {
	return func(builder *SelectBuilder) {
		builder.appendWhereCondition(fmt.Sprintf("%s>?", field), arg)
	}
}

// WithFieldGreaterOrEqualThan append where condition
//
// WHERE ${field} >= ?
func WithFieldGreaterOrEqualThan(field string, arg interface{}) SelectOption {
	return func(builder *SelectBuilder) {
		builder.appendWhereCondition(fmt.Sprintf("%s>=?", field), arg)
	}
}

// WithOrderBy append order by condition
//
// ORDER BY ${field} ASC
func WithOrderBy(field string) SelectOption {
	return func(builder *SelectBuilder) {
		builder.orderByConditions = append(builder.orderByConditions, fmt.Sprintf("%s ASC", field))
	}
}

// WithOrderByDESC append order by condition
//
// ORDER BY ${field} DESC
func WithOrderByDESC(field string) SelectOption {
	return func(builder *SelectBuilder) {
		builder.orderByConditions = append(builder.orderByConditions, fmt.Sprintf("%s DESC", field))
	}
}

// WithLimit set limit
//
// LIMIT ${limit}
func WithLimit(limit uint32) SelectOption {
	return func(builder *SelectBuilder) {
		builder.limit = limit
	}
}

// WithOffset set offset
//
// OFFSET ${offset}
func WithOffset(offset uint32) SelectOption {
	return func(builder *SelectBuilder) {
		builder.offset = offset
	}
}

// WithForceIndex set force index
//
// FORCE INDEX(${index})
func WithForceIndex(index string) SelectOption {
	return func(builder *SelectBuilder) {
		builder.forceIndex = index
	}
}
