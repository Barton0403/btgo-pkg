package middleware

import (
	"barton.top/btgo/pkg/http"
	"barton.top/btgo/pkg/http/engine"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
)

var handle http.HandlerFunc

func init() {
	gin.SetMode(gin.ReleaseMode)
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	before := m.HeapAlloc
	handle = DbMiddleware("root:root@tcp(192.168.137.3:3306)/sso")
	runtime.ReadMemStats(m)
	after := m.HeapAlloc
	println(after-before, "Bytes")
	println()
}

func BenchmarkDbMiddleware(b *testing.B) {
	router := gin.New()
	router.Use(engine.MakeGinHandlerFunc(handle))
	router.GET("/", engine.MakeGinHandlerFunc(func(ctx http.Context) {
		db := GetDefaultDb(ctx)
		var paswword string
		db.QueryRow("select password from user where username = ?", "test").Scan(&paswword)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		router.ServeHTTP(w, req)
		wg.Done()
	}
	wg.Wait()
}

func BenchmarkDbMiddleware2(b *testing.B) {
	router := gin.New()
	router.Use(engine.MakeGinHandlerFunc(DbMiddleware("root:root@tcp(192.168.137.3:3306)/sso")))
	router.GET("/", engine.MakeGinHandlerFunc(func(ctx http.Context) {
		db := GetDefaultDb(ctx)
		var paswword string
		db.QueryRow("select password from user where username = ?", "test").Scan(&paswword)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		router.ServeHTTP(w, req)
		wg.Done()
	}
	wg.Wait()
}

func BenchmarkCoDbMiddleware2(b *testing.B) {
	router := gin.New()
	router.Use(engine.MakeGinHandlerFunc(DbMiddleware("root:root@tcp(192.168.137.3:3306)/sso")))
	router.GET("/", engine.MakeGinHandlerFunc(func(ctx http.Context) {
		db := GetDefaultDb(ctx)
		var paswword string
		db.QueryRow("select password from user where username = ?", "test").Scan(&paswword)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			router.ServeHTTP(w, req)
			wg.Done()
		}()
	}
	wg.Wait()
}

func BenchmarkCoDbMiddleware(b *testing.B) {
	router := gin.New()
	router.Use(engine.MakeGinHandlerFunc(handle))
	router.GET("/", engine.MakeGinHandlerFunc(func(ctx http.Context) {
		db := GetDefaultDb(ctx)
		var paswword string
		db.QueryRow("select password from user where username = ?", "test").Scan(&paswword)
	}))

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()

	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			router.ServeHTTP(w, req)
			wg.Done()
		}()
	}
	wg.Wait()
}
