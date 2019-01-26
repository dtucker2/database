package query

import (
	"reflect"
	"strings"

	"github.com/jinzhu/inflection"
)

// QueryBuilder provides methods to generate MySQL queries from pointers to structs.
type QueryBuilder struct{}

// NewQueryBuilder returns a pointer to a new instance of the QueryBuilder struct.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// BuildInsertQuery constructs and returns a MySQL INSERT query and arguments from the passed object (struct pointer).
func (builder *QueryBuilder) BuildInsertQuery(object interface{}) (string, []interface{}) {
	columnNames, args := builder.getColumnNamesAndValues(object)
	values := strings.Repeat("?,", len(args))
	if len(values) > 0 {
		values = values[:len(values)-1]
	}
	return strings.Join([]string{
		"INSERT INTO",
		builder.getTableName(object),
		"(" + strings.Join(columnNames, ",") + ")",
		"VALUES",
		"(" + values + ")",
	}, " "), args
}

// BuildUpdateQuery constructs and returns a MySQL UPDATE query and arguments from the passed object (struct pointer).
func (builder *QueryBuilder) BuildUpdateQuery(object interface{}) (string, []interface{}) {
	columnNames, args := builder.getColumnNamesAndValues(object)
	args = append(args, builder.getPrimaryKeyColumnValue(object))
	return strings.Join([]string{
		"UPDATE",
		builder.getTableName(object),
		"SET",
		builder.buildUpdateValues(columnNames),
		"WHERE",
		builder.getPrimaryKeyColumnName(object) + "=?",
	}, " "), args
}

func (builder *QueryBuilder) getTableName(object interface{}) string {
	typ := reflect.TypeOf(object).Elem()
	return inflection.Plural(typ.Name())
}

func (builder *QueryBuilder) getColumnNamesAndValues(object interface{}) ([]string, []interface{}) {
	typ := reflect.TypeOf(object).Elem()
	val := reflect.ValueOf(object).Elem()
	columnNames := make([]string, 0)
	values := make([]interface{}, 0)
	for i := 0; i < typ.NumField(); i++ {
		columnNames = append(columnNames, typ.Field(i).Name)
		values = append(values, val.Field(i).Interface())
	}
	return columnNames, values
}

func (builder *QueryBuilder) buildUpdateValues(columnNames []string) string {
	str := ""
	for _, columnName := range columnNames {
		str += columnName + "=?,"
	}
	return str[:len(str)-1]
}

func (builder *QueryBuilder) getPrimaryKeyColumnName(object interface{}) string {
	return "Id"
}

func (builder *QueryBuilder) getPrimaryKeyColumnValue(object interface{}) interface{} {
	typ := reflect.TypeOf(object).Elem()
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Name == "Id" {
			return reflect.ValueOf(object).Elem().Field(i).Interface()
		}
	}
	return nil
}
