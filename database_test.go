package database_test

import (
	. "github.com/dtucker2/database"

	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase_Insert(t *testing.T) {
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
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectExec(`INSERT INTO Objects \(Id,Name,CreatedAt,UpdatedAt\) VALUES \(\?,\?,\?,\?\)`).
			WithArgs(0, "Test Object", (*time.Time)(nil), (*time.Time)(nil)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, NewDatabase(db).Insert(&object))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabase_Update(t *testing.T) {
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
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectExec(`UPDATE Objects SET Id=\?,Name=\?,CreatedAt=\?,UpdatedAt=\? WHERE Id=\?`).
			WithArgs(0, "Test Object", (*time.Time)(nil), (*time.Time)(nil), 0).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, NewDatabase(db).Update(&object))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
