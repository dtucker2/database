package query

import ()

type tableNamer interface {
	GetTableName() string
}
