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
		assert.Equal(t, `INSERT INTO Objects (Id,Name,CreatedAt,UpdatedAt) VALUES (?,?,?,?)`, builder.BuildInsertQuery(&object))
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
		assert.Equal(t, `UPDATE Objects SET Id=?,Name=?,CreatedAt=?,UpdatedAt=? WHERE Id=?`, builder.BuildUpdateQuery(&object))
	})
}
