package query_test

import (
	. "github.com/dtucker2/database/query"

	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type object struct {
	Id        int
	Name      string
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type objectWithTags struct {
	Id        int        `name:"id" type:"auto-increment"`
	Name      string     `name:"name"`
	CreatedAt *time.Time `name:"created_at" type:"created_at"`
	UpdatedAt *time.Time `name:"updated_at" type:"updated_at"`
}

func (obj *objectWithTags) GetTableName() string {
	return "objects"
}

func TestQueryBuilder_BuildInsertQuery(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		obj := object{
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args := builder.BuildInsertQuery(&obj)
		assert.Equal(t, `INSERT INTO objects (Id,Name,CreatedAt,UpdatedAt) VALUES (?,?,?,?)`, query)
		assert.Equal(t, []interface{}{0, "Test Object", (*time.Time)(nil), (*time.Time)(nil)}, args)
	})
	t.Run("tags", func(t *testing.T) {
		obj := objectWithTags{
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args := builder.BuildInsertQuery(&obj)
		assert.Equal(t, `INSERT INTO objects (name,created_at) VALUES (?,?)`, query)
		if assert.Len(t, args, 2) {
			assert.Equal(t, args[0], "Test Object")
			assert.NotNil(t, args[1])
			assert.IsType(t, time.Time{}, args[1])
		}
	})
}

func TestQueryBuilder_BuildUpdateQuery(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		obj := object{
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args, err := builder.BuildUpdateQuery(&obj)
		require.NoError(t, err)
		assert.Equal(t, `UPDATE objects SET Id=?,Name=?,CreatedAt=?,UpdatedAt=? WHERE Id=?`, query)
		assert.Equal(t, []interface{}{0, "Test Object", (*time.Time)(nil), (*time.Time)(nil), 0}, args)
	})
	t.Run("tags", func(t *testing.T) {
		obj := objectWithTags{
			Id:   1,
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args, err := builder.BuildUpdateQuery(&obj)
		require.NoError(t, err)
		assert.Equal(t, `UPDATE objects SET name=?,updated_at=? WHERE id=?`, query)
		if assert.Len(t, args, 3) {
			assert.Equal(t, args[0], "Test Object")
			assert.NotNil(t, args[1])
			assert.IsType(t, time.Time{}, args[1])
			assert.Equal(t, args[2], 1)
		}
	})
}

func TestQueryBuilder_BuildDeleteQuery(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		obj := object{
			Id: 1,
		}
		builder := NewQueryBuilder()
		query, args, err := builder.BuildDeleteQuery(&obj)
		require.NoError(t, err)
		assert.Equal(t, `DELETE FROM objects WHERE Id=?`, query)
		assert.Equal(t, []interface{}{1}, args)
	})
	t.Run("tags", func(t *testing.T) {
		obj := objectWithTags{
			Id: 1,
		}
		builder := NewQueryBuilder()
		query, args, err := builder.BuildDeleteQuery(&obj)
		require.NoError(t, err)
		assert.Equal(t, `DELETE FROM objects WHERE id=?`, query)
		assert.Equal(t, []interface{}{1}, args)
	})
}

func TestQueryBuilder_BuildSelectQuery(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		obj := object{
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args, err := builder.BuildSelectQuery(&obj)
		require.NoError(t, err)
		assert.Equal(t, `SELECT (Id,Name,CreatedAt,UpdatedAt) FROM objects WHERE Id=?`, query)
		assert.Equal(t, []interface{}{0}, args)
	})
	t.Run("tags", func(t *testing.T) {
		obj := objectWithTags{
			Id:   1,
			Name: "Test Object",
		}
		builder := NewQueryBuilder()
		query, args, err := builder.BuildSelectQuery(&obj)
		require.NoError(t, err)
		assert.Equal(t, `SELECT (id,name,created_at,updated_at) FROM objects WHERE id=?`, query)
		assert.Equal(t, []interface{}{1}, args)
	})
}
