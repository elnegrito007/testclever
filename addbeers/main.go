package main

import (
	"encoding/json"
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
	_ = http.ListenAndServe(":"+os.Getenv("PORT_CREATE_BEER"), routerCreate())
}

type Beers struct {
	Id       int32   `json:"id"`
	Name     string  `json:"Name"`
	Brewery  string  `json:"Brewery"`
	Country  string  `json:"Country"`
	Price    float32 `json:"Price"`
	Currency string  `json:"Currency"`
}

func routerCreate() http.Handler {
	r := chi.NewRouter()
	cors2 := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"POST"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	r.Use(cors2.Handler)
	r.Post("/beers/{beerID}", func(w http.ResponseWriter, r *http.Request) {
		beerID := chi.URLParam(r, "beerID")
		decoder := json.NewDecoder(r.Body)
		var t Beers
		err := decoder.Decode(&t)
		if err != nil {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`Request invalida`))
		}
		out, err := json.Marshal(t)
		if err != nil {
			w.WriteHeader(400)
			_, _ = w.Write([]byte(`Request invalida`))
		}
		exit, _ := conect.Get("Beer_" + string(beerID)).Result()
		if exit == "" {
			_ = conect.Set("Beer_"+string(beerID), string(out), 0).Err()
			w.WriteHeader(201)
			_, _ = w.Write([]byte(`Cerveza creada`))
		} else {
			w.WriteHeader(409)
			_, _ = w.Write([]byte(`El ID de la cerveza ya existe`))
		}
	})
	return r
}
