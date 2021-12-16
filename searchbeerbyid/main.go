package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"net/http"
	"os"
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
	_ = http.ListenAndServe(":"+os.Getenv("PORT_SEARCH_BY_ID_BEER"), routerCreate())
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
	r.Get("/beers/{beerID}", func(w http.ResponseWriter, r *http.Request) {
		beerID := chi.URLParam(r, "beerID")
		result, _ := conect.Get("Beer_" + beerID).Result()
		if result == "" {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`El Id de la cerveza no existe`))
		} else {
			_, _ = w.Write([]byte(result))
		}
	})
	return r
}
