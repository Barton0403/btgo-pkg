package middleware

import (
	"barton.top/btgo/pkg/common"
	"barton.top/btgo/pkg/http"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Connection = amqp.Connection

func MqMiddleware(source string) http.HandlerFunc {
	conn, err := amqp.Dial(source)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		panic(err)
	}

	return func(c http.Context) {
		c.Set(common.MqKey, conn)
	}
}

func GetDefaultMq(ctx http.Context) *Connection {
	if v, e := ctx.Get(common.MqKey); e {
		return v.(*Connection)
	}
	panic("no mq found")
}
