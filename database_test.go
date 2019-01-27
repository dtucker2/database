package database_test

import (
	. "github.com/dtucker2/database"

	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
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
	Id        int        `name:"id" type:"auto-increment" key:"true"`
	Name      string     `name:"name"`
	CreatedAt *time.Time `name:"created_at" type:"created_at"`
	UpdatedAt *time.Time `name:"updated_at" type:"updated_at"`
}

func (obj *objectWithTags) GetTableName() string {
	return "objects"
}

type anyTime struct{}

// Match satisfies sqlmock.Argument interface.
func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestDatabase_Insert(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		obj := object{
			Name: "Test Object",
		}
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectExec(`INSERT INTO objects \(Id,Name,CreatedAt,UpdatedAt\) VALUES \(\?,\?,\?,\?\)`).
			WithArgs(0, "Test Object", (*time.Time)(nil), (*time.Time)(nil)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, NewDatabase(db).Insert(&obj))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("tags", func(t *testing.T) {
		obj := objectWithTags{
			Name: "Test Object",
		}
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectExec(`INSERT INTO objects \(name,created_at\) VALUES \(\?,\?\)`).
			WithArgs("Test Object", anyTime{}).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, NewDatabase(db).Insert(&obj))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabase_Update(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		obj := object{
			Name: "Test Object",
		}
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectExec(`UPDATE objects SET Id=\?,Name=\?,CreatedAt=\?,UpdatedAt=\? WHERE Id=\?`).
			WithArgs(0, "Test Object", (*time.Time)(nil), (*time.Time)(nil), 0).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, NewDatabase(db).Update(&obj))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("tags", func(t *testing.T) {
		obj := objectWithTags{
			Id:   1,
			Name: "Test Object",
		}
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectExec(`UPDATE objects SET name=\?,updated_at=\? WHERE id=\?`).
			WithArgs("Test Object", anyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, NewDatabase(db).Update(&obj))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabase_Delete(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		obj := object{
			Id: 1,
		}
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectExec(`DELETE FROM objects WHERE Id=\?`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, NewDatabase(db).Delete(&obj))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("tags", func(t *testing.T) {
		obj := objectWithTags{
			Id: 1,
		}
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectExec(`DELETE FROM objects WHERE id=\?`).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		require.NoError(t, NewDatabase(db).Delete(&obj))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestDatabase_Select(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		obj := object{
			Id: 1,
		}
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectQuery(`SELECT \(Id,Name,CreatedAt,UpdatedAt\) FROM objects WHERE Id=\?`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"Id", "Name", "CreatedAt", "UpdatedAt"}).
				AddRow(1, "Test Object", (*time.Time)(nil), (*time.Time)(nil)))
		require.NoError(t, NewDatabase(db).Select(&obj))
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, object{Id: 1, Name: "Test Object", CreatedAt: nil, UpdatedAt: nil}, obj)
	})
	t.Run("tags", func(t *testing.T) {
		obj := objectWithTags{
			Id: 1,
		}
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectQuery(`SELECT \(id,name,created_at,updated_at\) FROM objects WHERE id=\?`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at", "updated_at"}).
				AddRow(1, "Test Object", (*time.Time)(nil), (*time.Time)(nil)))
		require.NoError(t, NewDatabase(db).Select(&obj))
		assert.NoError(t, mock.ExpectationsWereMet())
		assert.Equal(t, objectWithTags{Id: 1, Name: "Test Object", CreatedAt: nil, UpdatedAt: nil}, obj)
	})
}
