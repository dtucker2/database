package query_test

import (
	. "github.com/dtucker2/database/query"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestQueryBuilder_BuildInsertQuery(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type Object struct {
			Id        int
			Name      string
			CreatedAt *time.Time
			UpdatedAt *time.Time
		}
		object := Object{
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args := builder.BuildInsertQuery(&object)
		assert.Equal(t, `INSERT INTO Objects (Id,Name,CreatedAt,UpdatedAt) VALUES (?,?,?,?)`, query)
		assert.Equal(t, []interface{}{0, "Test Object", (*time.Time)(nil), (*time.Time)(nil)}, args)
	})
}

func TestQueryBuilder_BuildUpdateQuery(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type Object struct {
			Id        int
			Name      string
			CreatedAt *time.Time
			UpdatedAt *time.Time
		}
		object := Object{
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args := builder.BuildUpdateQuery(&object)
		assert.Equal(t, `UPDATE Objects SET Id=?,Name=?,CreatedAt=?,UpdatedAt=? WHERE Id=?`, query)
		assert.Equal(t, []interface{}{0, "Test Object", (*time.Time)(nil), (*time.Time)(nil), 0}, args)
	})
}

func TestQueryBuilder_BuildDeleteQuery(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type Object struct {
			Id        int
			Name      string
			CreatedAt *time.Time
			DeletedAt *time.Time
		}
		object := Object{
			Id: 1,
		}
		builder := NewQueryBuilder()
		query, args := builder.BuildDeleteQuery(&object)
		assert.Equal(t, `DELETE FROM Objects WHERE Id=?`, query)
		assert.Equal(t, []interface{}{1}, args)
	})
}

func TestQueryBuilder_BuildSelectQuery(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		type Object struct {
			Id        int
			Name      string
			CreatedAt *time.Time
			UpdatedAt *time.Time
		}
		object := Object{
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args := builder.BuildSelectQuery(&object)
		assert.Equal(t, `SELECT (Id,Name,CreatedAt,UpdatedAt) FROM Objects WHERE Id=?`, query)
		assert.Equal(t, []interface{}{0}, args)
	})
}
