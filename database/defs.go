package database

import (
	"log"
	"reflect"
	"strings"
)

type ColumnDef struct {
	Name string
	Primary bool
}

func (columnDef *ColumnDef) parse(val string) {
	parts := strings.Split(val, ",")

	columnDef.Name = parts[0]

	if len(parts) > 1 {
		columnDef.Primary = true
	}
}

type TableDef interface {
	TableName() string
}

func Join(defs []ColumnDef) string {
	if len(defs) == 0 {
		log.Println(ErrNoColumns)
		panic(ErrNoColumns)
	}

	// TODO: performance?
	b := strings.Builder{}

	for i, def := range defs {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString(def.Name)
	}

	return b.String()
}

func JoinFull(tableDef TableDef, defs []ColumnDef) string {
	b := strings.Builder{}

	for i, def := range defs {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString(tableDef.TableName())
		b.WriteString(".")
		b.WriteString(def.Name)
	}

	return b.String()
}

func MakeColumnDefs(tableDef TableDef, operation ColumnOperation) []ColumnDef {
	v := reflect.ValueOf(tableDef)
	v = reflect.Indirect(v)

	t := v.Type()

	defs := make([]ColumnDef, 0)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		name := f.Tag.Get(SqlTagName)

		if name != "" {
			def := ColumnDef{}
			def.parse(name)

			if operation == nil || operation.Wants(def) {
				defs = append(defs, def)
			}
		}
	}

	return defs
}

type ColumnOperation interface {
	Wants(columnDef ColumnDef) bool
}

type ExcludeOperation map[string]bool

func (op ExcludeOperation) Wants(columnDef ColumnDef) bool {
	_, ok := op[columnDef.Name]
	return !ok
}

func Exclude(names ...string) ColumnOperation {
	operations := make(ExcludeOperation)

	for _, name := range names {
		operations[name] = true
	}

	return operations
}

type IncludeOperation map[string]bool

func (op IncludeOperation) Wants(columnDef ColumnDef) bool {
	_, ok := op[columnDef.Name]
	return ok
}

func Include(names ...string) ColumnOperation {
	operations := make(IncludeOperation)

	for _, name := range names {
		operations[name] = true
	}

	return operations
}

type ExcludePrimaryKeyOperation struct {}

func (op ExcludePrimaryKeyOperation) Wants(columnDef ColumnDef) bool {
	return !columnDef.Primary
}

func ExcludePrimaryKey() ColumnOperation {
	return ExcludePrimaryKeyOperation{}
}

type ChainOperation struct {
	first ColumnOperation
	second ColumnOperation
}

func (op ChainOperation) Wants(columnDef ColumnDef) bool {
	return op.first.Wants(columnDef) || op.second.Wants(columnDef)
}

func Chain(first ColumnOperation, second ColumnOperation) ColumnOperation {
	return ChainOperation{ first: first, second: second }
}