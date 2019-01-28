# Database
[![GoDoc](https://godoc.org/github.com/dtucker2/database?status.svg)](https://godoc.org/github.com/dtucker2/database) [![Build Status](https://travis-ci.org/dtucker2/database.svg)](https://travis-ci.org/dtucker2/database) [![Go Report Card](https://goreportcard.com/badge/github.com/dtucker2/database)](https://goreportcard.com/report/github.com/dtucker2/database) [![codecov.io](https://codecov.io/github/dtucker2/database/branch/master/graph/badge.svg)](https://codecov.io/github/dtucker2/database)

Removes the need for building queries!
## Usage
### setup.sql
``` sql
DROP DATABASE IF EXISTS test;

CREATE DATABASE test;
USE test;

CREATE TABLE people (
	name varchar(256) not null,
	age smallint not null,
	created_at TIMESTAMP not null,
	updated_at TIMESTAMP null,
	primary key (name)
);
```
### main.go
``` go
package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/dtucker2/database"
	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	Name      string     `name:"name" key:"true"`
	Age       int        `name:"age"`
	CreatedAt *time.Time `name:"created_at" type:"created_at"`
	UpdatedAt *time.Time `name:"updated_at" type:"updated_at"`
}

func (person *Person) GetTableName() string {
	return "people"
}

func main() {

	sqlDB, _ := sql.Open("mysql", "user:password@/test?parseTime=true")
	db := database.NewDatabase(sqlDB)

	person := Person{Name: "Frank", Age: 32}

	db.Insert(&person)
	/* Result:
	+-------+-----+---------------------+------------+
	| name  | age | created_at          | updated_at |
	+-------+-----+---------------------+------------+
	| Frank |  32 | 2019-01-28 23:20:00 | NULL       |
	+-------+-----+---------------------+------------+
	*/

	person.Age = 33
	db.Update(&person)
	/* Result:
	+-------+-----+---------------------+---------------------+
	| name  | age | created_at          | updated_at          |
	+-------+-----+---------------------+---------------------+
	| Frank |  33 | 2019-01-28 23:20:00 | 2019-01-28 23:20:00 |
	+-------+-----+---------------------+---------------------+
	*/

	person = Person{Name: "Frank"} // Note the lack of any age.
	db.Select(&person)
	fmt.Printf("%+v\n", person)
	/* Output:
	{Name:Frank Age:33 CreatedAt:2019-01-28 23:20:00 +0000 UTC UpdatedAt:2019-01-28 23:20:01 +0000 UTC}
	*/

	db.Delete(&person)
	/* Result:
	Empty set (0.00 sec)
	*/

}
```
