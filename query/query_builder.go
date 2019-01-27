package query

import (
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/inflection"
	"github.com/pkg/errors"
)

const (
	tagName              = "name"
	tagType              = "type"
	tagKey               = "key"
	tagTypeAutoIncrement = "auto-increment"
	tagTypeCreatedAt     = "created_at"
	tagTypeUpdatedAt     = "updated_at"
)

// QueryBuilder provides methods to generate MySQL queries from pointers to structs.
type QueryBuilder struct{}

// NewQueryBuilder returns a pointer to a new instance of the QueryBuilder struct.
func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

// BuildInsertQuery constructs and returns a MySQL INSERT query and arguments from the passed object (struct pointer).
func (builder *QueryBuilder) BuildInsertQuery(object interface{}) (string, []interface{}) {
	columnNames, args := builder.getColumnNamesAndValues(object, true)
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
func (builder *QueryBuilder) BuildUpdateQuery(object interface{}) (string, []interface{}, error) {
	keyName, keyValue, err := builder.getPrimaryKeyNameAndValue(object)
	if err != nil {
		return "", nil, err
	}
	columnNames, args := builder.getColumnNamesAndValues(object, false)
	return strings.Join([]string{
		"UPDATE",
		builder.getTableName(object),
		"SET",
		builder.buildUpdateValues(columnNames),
		"WHERE",
		keyName + "=?",
	}, " "), append(args, keyValue), nil
}

// BuildDeleteQuery constructs and returns a MySQL DELETE query and arguments from the passed object (struct pointer).
func (builder *QueryBuilder) BuildDeleteQuery(object interface{}) (string, []interface{}, error) {
	keyName, keyValue, err := builder.getPrimaryKeyNameAndValue(object)
	if err != nil {
		return "", nil, err
	}
	return strings.Join([]string{
		"DELETE",
		"FROM",
		builder.getTableName(object),
		"WHERE",
		keyName + "=?",
	}, " "), []interface{}{keyValue}, nil
}

// BuildSelectQuery constructs and returns a MySQL SELECT query and arguments from the passed object (struct pointer).
func (builder *QueryBuilder) BuildSelectQuery(object interface{}) (string, []interface{}, error) {
	keyName, keyValue, err := builder.getPrimaryKeyNameAndValue(object)
	if err != nil {
		return "", nil, err
	}
	return strings.Join([]string{
		"SELECT",
		"(" + strings.Join(builder.getColumnNames(object), ",") + ")",
		"FROM",
		builder.getTableName(object),
		"WHERE",
		keyName + "=?",
	}, " "), []interface{}{keyValue}, nil
}

func (builder *QueryBuilder) getTableName(object interface{}) string {
	if namer, ok := object.(tableNamer); ok {
		return namer.GetTableName()
	}
	typ := reflect.TypeOf(object).Elem()
	return inflection.Plural(typ.Name())
}

func (builder *QueryBuilder) getColumnNamesAndValues(object interface{}, insertion bool) ([]string, []interface{}) {
	typ := reflect.TypeOf(object).Elem()
	val := reflect.ValueOf(object).Elem()
	columnNames := make([]string, 0)
	values := make([]interface{}, 0)
	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		if builder.fieldShouldBeInserted(structField, insertion) {
			columnNames = append(columnNames, builder.getFieldName(structField))
			values = append(values, builder.getFieldValue(structField, val.Field(i)))
		}
	}
	return columnNames, values
}

func (builder *QueryBuilder) fieldShouldBeInserted(field reflect.StructField, insertion bool) bool {
	switch field.Tag.Get(tagType) {
	case tagTypeAutoIncrement:
		return false
	case tagTypeCreatedAt:
		return insertion
	case tagTypeUpdatedAt:
		return !insertion
	}
	return true
}

func (builder *QueryBuilder) getFieldName(structField reflect.StructField) string {
	name := structField.Tag.Get(tagName)
	if name != "" {
		return name
	}
	return structField.Name
}

func (builder *QueryBuilder) getFieldValue(structField reflect.StructField, value reflect.Value) interface{} {
	switch structField.Tag.Get(tagType) {
	case tagTypeCreatedAt, tagTypeUpdatedAt:
		return time.Now()
	}
	return value.Interface()
}

func (builder *QueryBuilder) getColumnNames(object interface{}) []string {
	typ := reflect.TypeOf(object).Elem()
	columnNames := make([]string, 0)
	for i := 0; i < typ.NumField(); i++ {
		columnNames = append(columnNames, builder.getFieldName(typ.Field(i)))
	}
	return columnNames
}

func (builder *QueryBuilder) buildUpdateValues(columnNames []string) string {
	str := ""
	for _, columnName := range columnNames {
		str += columnName + "=?,"
	}
	return str[:len(str)-1]
}

func (builder *QueryBuilder) getPrimaryKeyNameAndValue(object interface{}) (string, interface{}, error) {
	typ := reflect.TypeOf(object).Elem()
	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		if builder.hasPrimaryKeyTag(structField) {
			return builder.getFieldName(structField), reflect.ValueOf(object).Elem().Field(i).Interface(), nil
		}
	}
	// Attempt to default to any field with a name of 'Id' or a name tag of 'id'.
	for i := 0; i < typ.NumField(); i++ {
		structField := typ.Field(i)
		name := builder.getFieldName(structField)
		switch name {
		case "Id", "id":
			return name, reflect.ValueOf(object).Elem().Field(i).Interface(), nil
		}
	}
	return "", nil, errors.Errorf("Unable to identify primary key (struct is missing a '%s:\"true\"' tag).", tagKey)
}

func (builder *QueryBuilder) hasPrimaryKeyTag(structField reflect.StructField) bool {
	return structField.Tag.Get(tagKey) == "true"
}
