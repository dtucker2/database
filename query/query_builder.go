package query

import (
	"reflect"
	"strings"

	"github.com/jinzhu/inflection"
)

// QueryBuilder provides methods to generate MySQL queries from structs.
type QueryBuilder struct{}

// NewQueryBuilder returns a pointer to a new instance of the QueryBuilder struct.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// BuildInsertQuery constructs and returns a MySQL INSERT query from the passed object.
func (builder *QueryBuilder) BuildInsertQuery(object interface{}) string {
	return strings.Join([]string{
		"INSERT INTO",
		builder.getTableName(object),
		"(" + strings.Join(builder.getColumnNames(object), ",") + ")",
		"VALUES",
		"(" + strings.Join(builder.getValues(object), ",") + ")",
	}, " ")
}

// BuildUpdateQuery constructs and returns a MySQL UPDATE query from the passed object.
func (builder *QueryBuilder) BuildUpdateQuery(object interface{}) string {
	return strings.Join([]string{
		"UPDATE",
		builder.getTableName(object),
		"SET",
		builder.buildUpdateValues(object),
		"WHERE",
		builder.getPrimaryKeyColumnName(object) + "=?",
	}, " ")
}

func (builder *QueryBuilder) getTableName(object interface{}) string {
	typ := reflect.TypeOf(object).Elem()
	return inflection.Plural(typ.Name())
}

func (builder *QueryBuilder) getColumnNames(object interface{}) []string {
	typ := reflect.TypeOf(object).Elem()
	columnNames := make([]string, 0)
	for i := 0; i < typ.NumField(); i++ {
		columnNames = append(columnNames, typ.Field(i).Name)
	}
	return columnNames
}

func (builder *QueryBuilder) getValues(object interface{}) []string {
	typ := reflect.TypeOf(object).Elem()
	values := make([]string, 0)
	for i := 0; i < typ.NumField(); i++ {
		values = append(values, "?")
	}
	return values
}

func (builder *QueryBuilder) buildUpdateValues(object interface{}) string {
	typ := reflect.TypeOf(object).Elem()
	values := make([]string, 0)
	for i := 0; i < typ.NumField(); i++ {
		values = append(values, typ.Field(i).Name+"=?")
	}
	return strings.Join(values, ",")
}

func (builder *QueryBuilder) getPrimaryKeyColumnName(object interface{}) string {
	return "Id"
}
