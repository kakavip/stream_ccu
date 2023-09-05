package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

var logger Logging = NewLogger()
var rdb *redis.Client = redis.NewClient(&redis.Options{
	Addr:     os.Getenv("REDIS_HOST"),
	Password: "", // no password set
	DB:       0,  // use default DB
})
var rCtx = context.Background()

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	port := 8080

	limitCcuApi := LimitCCUAPI{}
	r.Mount("/tenants/{tenant_slug}/accounts/{account_id}/csl/{action}", limitCcuApi.Routers())

	logger.info(fmt.Sprintf("Server is running at port %d", port))
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
