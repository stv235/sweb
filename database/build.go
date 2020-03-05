package database

import "strings"

func BuildSelect(tableDef TableDef, where string, orderBy string, operations ColumnOperation) string {
	b := strings.Builder{}

	b.WriteString("SELECT ")
	b.WriteString(Join(MakeColumnDefs(tableDef, operations)))
	b.WriteString(" FROM ")
	b.WriteString(tableDef.TableName())

	if where != "" {
		b.WriteString(" WHERE ")
		b.WriteString(where)
	}

	if orderBy != "" {
		b.WriteString(" ORDER BY ")
		b.WriteString(orderBy)
	}

	return b.String()
}

func BuildJoinSelect(tableDef TableDef, tableNames string, where string, orderBy string, operations ColumnOperation) string {
	b := strings.Builder{}

	b.WriteString("SELECT ")
	b.WriteString(JoinFull(tableDef, MakeColumnDefs(tableDef, operations)))
	b.WriteString(" FROM ")
	b.WriteString(tableDef.TableName())

	if tableNames != "" {
		b.WriteString(",")
		b.WriteString(tableNames)
	}

	if where != "" {
		b.WriteString(" WHERE ")
		b.WriteString(where)
	}

	if orderBy != "" {
		b.WriteString(" ORDER BY ")
		b.WriteString(orderBy)
	}

	return b.String()
}
