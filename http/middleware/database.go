package middleware

import (
	"barton.top/btgo/pkg/common"
	"barton.top/btgo/pkg/http"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DB = sql.DB

func DbMiddleware(dataSourceName string) http.HandlerFunc {
	// Opening a driver typically will not attempt to connect to the database.
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		panic(err)
	}
	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	return func(c http.Context) {
		c.Set(common.DatabaseKey, db)
	}
}

func GetDefaultDb(ctx http.Context) *DB {
	if v, e := ctx.Get(common.DatabaseKey); e {
		return v.(*DB)
	}
	panic("no db found")
}
