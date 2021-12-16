package main

import (
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	_ = http.ListenAndServe(":"+os.Getenv("PORT_SEARCH_BOX_PRICE_BY_ID_BEER"), routerCreate())
}

func convertion(currencyPay string, currencyBeer string, quantity float32, price float32) float32 {
	if strings.ToUpper(currencyPay) == currencyBeer {
		amount := quantity * price
		return amount
	} else {
		req, _ := http.NewRequest("GET", "http://api.currencylayer.com/live?access_key=dcd11bc02ed484cf2a99a2341488245d&format=1", nil)
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		inUSDpay, _ := jsonparser.GetFloat(body, "quotes", "USD"+currencyPay)
		inUSDsell, _ := jsonparser.GetFloat(body, "quotes", "USD"+currencyBeer)
		if inUSDpay > 0 && inUSDsell > 0 {
			amount := ((price / float32(inUSDsell)) * float32(inUSDpay)) * quantity
			return amount
		}
	}
	return 0
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
	r.Get("/beers/{beerID}/boxprice", func(w http.ResponseWriter, r *http.Request) {
		beerID := chi.URLParam(r, "beerID")
		currency := r.URL.Query().Get("currency")
		quantity, _ := strconv.ParseFloat(r.URL.Query().Get("quantity"), 64)
		result, _ := conect.Get("Beer_" + beerID).Result()
		if result == "" {
			w.WriteHeader(404)
			_, _ = w.Write([]byte(`El Id de la cerveza no existe`))
		} else {
			price, _ := jsonparser.GetFloat([]byte(result), "Price") //https://github.com/buger/jsonparser
			currencyBeer, _ := jsonparser.GetString([]byte(result), "Currency")
			amount := convertion(currency, currencyBeer, float32(quantity), float32(price))
			_, _ = w.Write([]byte(`{"Price Total":` + fmt.Sprintf("%f", amount) + `}`))
		}
	})
	return r
}
