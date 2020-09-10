package database

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
)

var ErrNoPrimaryKey = errors.New("def must have a primary key")
var ErrNoColumns = errors.New("def must have columns")
var ErrArrayExpected = errors.New("def is not an array")

func CommonChain(tableDef TableDef, fn func(columnDef ColumnDef, v reflect.Value)) {
	v := reflect.ValueOf(tableDef)
	v = reflect.Indirect(v)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tagValue := f.Tag.Get(SqlTagName)

		if tagValue != "" {
			columnDef := ColumnDef{}
			columnDef.parse(tagValue)

			fn(columnDef, v.Field(i))
		}
	}
}

func makeUpdateString(tableName string, columnNames []string, primaryKeyNames []string) string {
	b := strings.Builder{}
	b.WriteString("UPDATE ")
	b.WriteString(tableName)
	b.WriteString(" SET ")

	for i, columnName := range columnNames {
		if i > 0 {
			b.WriteString(", ")
		}

		b.WriteString(columnName)
		b.WriteString("=?")
	}

	b.WriteString(" WHERE ")

	for i, primaryKeyName := range primaryKeyNames {
		if i > 0 {
			b.WriteString(" AND ")
		}

		b.WriteString(primaryKeyName)
		b.WriteString("=?")
	}

	return b.String()
}

func Update(db Q, tableDef TableDef, ops ColumnOperation) {
	whereValues := make([]interface{}, 0)
	whereNames := make([]string, 0)
	updateValues := make([]interface{}, 0)
	updateNames := make([]string, 0)

	CommonChain(tableDef, func(columnDef ColumnDef, v reflect.Value) {
		switch {
		case columnDef.Primary:
			whereValues = append(whereValues, v.Interface())
			whereNames = append(whereNames, columnDef.Name)
		case !columnDef.ReadOnly && (ops == nil || ops.Wants(columnDef)):
			updateValues = append(updateValues, v.Interface())
			updateNames = append(updateNames, columnDef.Name)
		}
	})

	str := makeUpdateString(tableDef.TableName(), updateNames, whereNames)

	_, err := db.Exec(str, append(updateValues, whereValues...)...)

	if err != nil {
		log.Output(2, err.Error())
		panic(err)
	}
}

func UpdateWhere(db Q, tableDef TableDef, ops ColumnOperation, where string, args ...interface{}) {
	updateValues := make([]interface{}, 0)
	updateNames := make([]string, 0)

	CommonChain(tableDef, func(columnDef ColumnDef, v reflect.Value) {
		switch {
		case ops == nil || ops.Wants(columnDef):
			updateValues = append(updateValues, v.Interface())
			updateNames = append(updateNames, columnDef.Name)
		}
	})

	w := strings.Builder{}

	w.WriteString("UPDATE ")
	w.WriteString(tableDef.TableName())
	w.WriteString(" SET ")

	for i, updateName := range updateNames {
		if i > 0 {
			w.WriteString(",")
		}

		w.WriteString(updateName)
		w.WriteString("=?")
	}

	if where != "" {
		w.WriteString(" WHERE ")
		w.WriteString(where)
	}

	_, err := db.Exec(w.String(), append(updateValues, args...)...)

	if err != nil {
		log.Output(2, err.Error())
		panic(err)
	}
}

func makeInsertString(table string, insertNames []string) string {
	b := strings.Builder{}

	b.WriteString("INSERT INTO ")
	b.WriteString(table)
	b.WriteString("(")

	for i, column := range insertNames {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString(column)
	}

	b.WriteString(")VALUES(")

	for i := range insertNames {
		if i > 0 {
			b.WriteString(",")
		}

		b.WriteString("?")
	}

	b.WriteString(")")

	return b.String()
}

func Insert(db Q, tableDef TableDef, ops ColumnOperation) {
	whereValues := make([]interface{}, 0)
	whereNames := make([]string, 0)
	insertValues := make([]interface{}, 0)
	insertNames := make([]string, 0)

	checkAddressable(tableDef)

	var id *int64

	CommonChain(tableDef, func(columnDef ColumnDef, v reflect.Value) {
		switch {
		case columnDef.Primary:
			whereValues = append(whereValues, v.Interface())
			whereNames = append(whereNames, columnDef.Name)

			if v.Type().Kind() == reflect.Int64 {
				id = v.Addr().Interface().(*int64)
			}
		case ops == nil || ops.Wants(columnDef):
			insertValues = append(insertValues, v.Interface())
			insertNames = append(insertNames, columnDef.Name)
		}
	})

	// multi-primary-key-table? set values to insert
	if len(whereNames) > 1 {
		insertNames = append(insertNames, whereNames...)
		insertValues = append(insertValues, whereValues...)
	}

	str := makeInsertString(tableDef.TableName(), insertNames)

	res, err := db.Exec(str, insertValues...)

	if err != nil {
		log.Output(2, err.Error())
		panic(err)
	}

	if len(whereValues) == 1 && id != nil {
		tempId, err := res.LastInsertId()

		if err != nil {
			log.Panicln(err)
		}

		*id = tempId
	}
}

/**
Deletes by primary key
 */
func Delete(db Q, tableDef TableDef) {
	b := strings.Builder{}

	b.WriteString("DELETE FROM ")
	b.WriteString(tableDef.TableName())
	b.WriteString(" WHERE ")

	whereValues := make([]interface{}, 0)
	whereNames := make([]string, 0)

	CommonChain(tableDef, func(columnDef ColumnDef, v reflect.Value) {
		if columnDef.Primary {
			whereNames = append(whereNames, columnDef.Name)
			whereValues = append(whereValues, v.Interface())
		}
	})

	for i, whereName := range whereNames {
		if i > 0 {
			b.WriteString(" AND ")
		}

		b.WriteString(whereName)
		b.WriteString("=?")
	}

	_, err := db.Exec(b.String(), whereValues...)

	if err != nil {
		log.Output(2, err.Error())
		panic(err)
	}
}

func DeleteWhere(db Q, tableDef TableDef, where string, args ...interface{}) {
	b := strings.Builder{}

	b.WriteString("DELETE FROM ")
	b.WriteString(tableDef.TableName())
	b.WriteString(" WHERE ")
	b.WriteString(where)

	_, err := db.Exec(b.String(), args...)

	if err != nil {
		log.Output(2, err.Error())
		panic(err)
	}
}

func Scan(res *sql.Rows, obj TableDef, ops ColumnOperation) {
	if err := res.Scan(MakeValueRefs(obj, ops)...); err != nil {
		log.Output(2, err.Error())
		panic(err)
	}
}

func ScanCustom(res *sql.Rows, params... interface{}) {
	if err := res.Scan(params...); err != nil {
		log.Output(2, err.Error())
		panic(err)
	}
}

func ScanOptionalRow(res *sql.Row, tableDef TableDef, ops ColumnOperation) bool {
	switch err := res.Scan(MakeValueRefs(tableDef, ops)...); {
	case err == sql.ErrNoRows:
		return false
	case err != nil:
		log.Output(2, err.Error())
		panic(err)
	}

	return true
}

func MakeValueRefs(tableDef TableDef, ops ColumnOperation) []interface{} {
	v := reflect.ValueOf(tableDef)
	v = reflect.Indirect(v)
	t := v.Type()

	params := make([]interface{}, 0)

	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		tag := f.Tag.Get(SqlTagName)

		if tag != "" {
			columnDef := ColumnDef{}
			columnDef.parse(tag)

			if ops == nil || ops.Wants(columnDef) {
				params = append(params, v.Field(i).Addr().Interface())
			}
		}
	}

	return params
}

func ScanRow(res *sql.Row, tableDef TableDef, ops ColumnOperation) {
	params := MakeValueRefs(tableDef, ops)

	if err := res.Scan(params...); err != nil {
		log.Output(3, err.Error())
		panic(err)
	}
}

func SelectRow(db Q, tableDef TableDef, ops ColumnOperation) {
	whereValues := make([]interface{}, 0)
	whereNames := make([]string, 0)

	checkAddressable(tableDef)

	CommonChain(tableDef, func(columnDef ColumnDef, v reflect.Value) {
		switch {
		case columnDef.Primary:
			whereValues = append(whereValues, v.Addr().Interface())
			whereNames = append(whereNames, columnDef.Name)
		}
	})

	if len(whereNames) == 0 {
		log.Output(2, fmt.Sprintf("table %v must have primary key", tableDef.TableName()))
		panic(ErrNoPrimaryKey)
	}

	where := ""

	for i, name := range whereNames {
		if i > 0 {
			where += " AND "
		}

		where += name + "=?"
	}

	str := BuildSelect(tableDef, where,"", ops)
	res := db.QueryRow(str, whereValues...)

	ScanRow(res, tableDef, ops)
}

func SelectRowWhere(db Q, tableDef TableDef, ops ColumnOperation, where string, args... interface{}) bool {
	checkAddressable(tableDef)

	str := BuildSelect(tableDef, where,"", ops)
	res := db.QueryRow(str, args...)

	return ScanOptionalRow(res, tableDef, ops)
}

func checkAddressable(tableDef TableDef) {
	v := reflect.ValueOf(tableDef)

	if v.Type().Kind() != reflect.Ptr {
		err := fmt.Errorf("table '%v' must be a pointer", tableDef.TableName())
		log.Output(3, err.Error())
		panic(err)
	}
}

func SelectOptionalRow(db Q, tableDef TableDef, ops ColumnOperation) bool {
	whereValues := make([]interface{}, 0)
	whereNames := make([]string, 0)

	checkAddressable(tableDef)

	CommonChain(tableDef, func(columnDef ColumnDef, v reflect.Value) {
		switch {
		case columnDef.Primary:
			whereValues = append(whereValues, v.Addr().Interface())
			whereNames = append(whereNames, columnDef.Name)
		}
	})

	if len(whereNames) == 0 {
		log.Output(2, fmt.Sprintf("table %v must have primary key", tableDef.TableName()))
		panic(ErrNoPrimaryKey)
	}

	where := ""

	for i, name := range whereNames {
		if i > 0 {
			where += " AND "
		}

		where += name + "=?"
	}

	str := BuildSelect(tableDef, where,"", ops)
	res := db.QueryRow(str, whereValues...)

	return ScanOptionalRow(res, tableDef, ops)
}

func Select(db Q, tableDefs interface{}, op ColumnOperation) {
	st := reflect.TypeOf(tableDefs)

	if st.Kind() != reflect.Ptr {
		log.Output(2, "slice must be a pointer")
		panic(ErrArrayExpected)
	}

	st = st.Elem()

	if st.Kind() != reflect.Slice {
		log.Output(2, ErrArrayExpected.Error())
		panic(ErrArrayExpected)
	}

	t := st.Elem()

	if t.Kind() != reflect.Ptr {
		err := fmt.Errorf("array elements must be pointers")
		log.Output(2, err.Error())
		panic(err)
	}

	t = t.Elem()

	tableDef := reflect.New(t).Interface().(TableDef)

	str := BuildSelect(tableDef, "", "", op)
	res := Query(db, str)

	sv := reflect.ValueOf(tableDefs)
	sv = reflect.Indirect(sv)

	for res.Next() {
		Scan(res, tableDef, op)
		sv = reflect.Append(sv, reflect.ValueOf(tableDef))
		tableDef = reflect.New(t).Interface().(TableDef)
	}

	reflect.ValueOf(tableDefs).Elem().Set(sv)
}

func SelectOrder(db Q, tableDefs interface{}, op ColumnOperation, order string) {
	st := reflect.TypeOf(tableDefs)

	if st.Kind() != reflect.Ptr {
		log.Output(2, "slice must be a pointer")
		panic(ErrArrayExpected)
	}

	st = st.Elem()

	if st.Kind() != reflect.Slice {
		log.Output(2, ErrArrayExpected.Error())
		panic(ErrArrayExpected)
	}

	t := st.Elem()

	if t.Kind() != reflect.Ptr {
		err := fmt.Errorf("array elements must be pointers")
		log.Output(2, err.Error())
		panic(err)
	}

	t = t.Elem()

	tableDef := reflect.New(t).Interface().(TableDef)

	str := BuildSelect(tableDef, "", order, op)
	res := Query(db, str)

	sv := reflect.ValueOf(tableDefs)
	sv = reflect.Indirect(sv)

	for res.Next() {
		Scan(res, tableDef, op)
		sv = reflect.Append(sv, reflect.ValueOf(tableDef))
		tableDef = reflect.New(t).Interface().(TableDef)
	}

	reflect.ValueOf(tableDefs).Elem().Set(sv)
}

func SelectWhere(db Q, tableDefs interface{}, op ColumnOperation, where string, args... interface{}) {
	st := reflect.TypeOf(tableDefs)

	if st.Kind() != reflect.Ptr {
		log.Output(2, "slice must be a pointer")
		panic(ErrArrayExpected)
	}

	st = st.Elem()

	if st.Kind() != reflect.Slice {
		log.Output(2, ErrArrayExpected.Error())
		panic(ErrArrayExpected)
	}

	t := st.Elem()

	if t.Kind() != reflect.Ptr {
		err := fmt.Errorf("array elements must be pointers")
		log.Output(2, err.Error())
		panic(err)
	}

	t = t.Elem()

	tableDef := reflect.New(t).Interface().(TableDef)

	str := BuildSelect(tableDef, where, "", op)
	res := Query(db, str, args...)

	sv := reflect.ValueOf(tableDefs)
	sv = reflect.Indirect(sv)

	for res.Next() {
		Scan(res, tableDef, op)
		sv = reflect.Append(sv, reflect.ValueOf(tableDef))
		tableDef = reflect.New(t).Interface().(TableDef)
	}

	reflect.ValueOf(tableDefs).Elem().Set(sv)
}

func SelectWhereOrder(db Q, tableDefs interface{}, op ColumnOperation, order string, where string, args... interface{}) {
	st := reflect.TypeOf(tableDefs)

	if st.Kind() != reflect.Ptr {
		log.Output(2, "slice must be a pointer")
		panic(ErrArrayExpected)
	}

	st = st.Elem()

	if st.Kind() != reflect.Slice {
		log.Output(2, ErrArrayExpected.Error())
		panic(ErrArrayExpected)
	}

	t := st.Elem()

	if t.Kind() != reflect.Ptr {
		err := fmt.Errorf("array elements must be pointers")
		log.Output(2, err.Error())
		panic(err)
	}

	t = t.Elem()

	tableDef := reflect.New(t).Interface().(TableDef)

	str := BuildSelect(tableDef, where, order, op)
	res := Query(db, str, args...)

	sv := reflect.ValueOf(tableDefs)
	sv = reflect.Indirect(sv)

	for res.Next() {
		Scan(res, tableDef, op)
		sv = reflect.Append(sv, reflect.ValueOf(tableDef))
		tableDef = reflect.New(t).Interface().(TableDef)
	}

	reflect.ValueOf(tableDefs).Elem().Set(sv)
}