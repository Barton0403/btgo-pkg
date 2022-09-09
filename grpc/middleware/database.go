package middleware

import (
	"barton.top/btgo/pkg/common"
	"context"
	"database/sql"
	"google.golang.org/grpc"
)

type DB = sql.DB

func DBUnaryServerInterceptor(dataSourceName string) grpc.UnaryServerInterceptor {
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

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(context.WithValue(ctx, common.DatabaseKey, db), req)
	}
}

func GetDefaultDb(ctx context.Context) *DB {
	if v := ctx.Value(common.DatabaseKey); v != nil {
		return v.(*DB)
	}
	panic("no db found")
}
