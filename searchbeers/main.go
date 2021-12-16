package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strings"
)

func DbConect() (cliente *redis.Client) {
	_ = godotenv.Load(".env")
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("HOST_REDIS"),
		Password: os.Getenv("PASS_REDIS"),
		DB:       0,
	})
	return client
}

var conect = DbConect()

func main() {
	_ = godotenv.Load(".env")
	_ = http.ListenAndServe(":"+os.Getenv("PORT_SEARCH_BEER"), routerCreate())
}

func handle(w http.ResponseWriter, r *http.Request) {
	valor, _ := conect.Keys("Beer_*").Result()
	if len(valor) == 0 {
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`[]`))
	} else {
		result, _ := conect.MGet(valor...).Result()
		stringedIDs := fmt.Sprintf("%v", result)
		stringedIDs = stringedIDs[1 : len(stringedIDs)-1]
		stringedIDs = strings.ReplaceAll(stringedIDs, "} {", "},{")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`[` + stringedIDs + `]`))
	}
}

func routerCreate() http.Handler {
	r := chi.NewRouter()
	cors2 := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors2.Handler)
	r.Get("/beers", handle)
	return r
}
