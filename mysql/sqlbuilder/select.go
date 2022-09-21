package sqlbuilder

import (
	"bytes"
	"fmt"
	"strings"
)

// SelectBuilder select statement builder
type SelectBuilder struct {
	table             string
	whereConditions   []string
	whereArgs         []interface{}
	orderByConditions []string
	limit             uint32
	offset            uint32
	forceIndex        string
}

// NewSelect new select statement builder
func NewSelect(table string, opts ...SelectOption) *SelectBuilder {
	builder := &SelectBuilder{
		table: table,
	}

	for _, opt := range opts {
		opt(builder)
	}

	return builder
}

// Build select statement
func (s *SelectBuilder) Build(fields ...string) (string, []interface{}) {
	if len(fields) == 0 {
		fields = []string{"*"}
	}

	sql := bytes.NewBufferString("SELECT ")
	sql.WriteString(strings.Join(fields, ", "))
	sql.WriteString(fmt.Sprintf(" FROM %s", s.table))
	sql.WriteString(s.buildForceIndex())
	sql.WriteString(s.buildWhere())
	sql.WriteString(s.buildOrderBy())
	sql.WriteString(s.buildLimit())
	sql.WriteString(s.buildOffset())

	return sql.String(), s.whereArgs
}

func (s *SelectBuilder) buildForceIndex() string {
	if s == nil || s.forceIndex == "" {
		return ""
	}

	return fmt.Sprintf(" FORCE INDEX(%s)", s.forceIndex)
}

func (s *SelectBuilder) buildWhere() string {
	if s == nil || len(s.whereConditions) == 0 {
		return ""
	}

	return fmt.Sprintf(" WHERE %s", strings.Join(s.whereConditions, " AND "))
}

func (s *SelectBuilder) buildOrderBy() string {
	if s == nil || len(s.orderByConditions) == 0 {
		return ""
	}

	return fmt.Sprintf(" ORDER BY %s", strings.Join(s.orderByConditions, ", "))
}

func (s *SelectBuilder) buildLimit() string {
	if s == nil || s.limit == 0 {
		return ""
	}

	return fmt.Sprintf(" LIMIT %d", s.limit)
}

func (s *SelectBuilder) buildOffset() string {
	if s == nil {
		return ""
	}

	return fmt.Sprintf(" OFFSET %d", s.offset)
}

func (s *SelectBuilder) appendWhereCondition(condition string, arg interface{}) {
	s.whereConditions = append(s.whereConditions, condition)
	s.whereArgs = append(s.whereArgs, arg)
}
