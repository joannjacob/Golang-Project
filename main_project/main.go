package main

import (
	"main_project/auth"
	"main_project/cronjob"
	"main_project/handlers"
	"main_project/logging"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/robfig/cron"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var wg sync.WaitGroup

func runCronJobs() {
	c := cron.New()
	// To import product data from csv
	c.AddFunc("@every 1h", cronjob.ImportData)
	c.Start()
	time.Sleep(10 * time.Second)
	wg.Done()
}

func registerRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(auth.JwtAuthentication)

	r.Route("/api", func(r chi.Router) {
		//Authentication APIs
		r.Post("/signup", handlers.CreateAccount)
		r.Post("/login", handlers.Authenticate)
		r.Post("/token/refresh", handlers.Refresh)
		// Product APIs
		r.Get("/get_product", handlers.GetProducts)
		r.Get("/get_product/{id}", handlers.GetProductById)
		r.Post("/add_product", handlers.CreateProduct)
		r.Put("/update_product/{id}", handlers.UpdateProduct)
		r.Delete("/delete_product/{id}", handlers.DeleteProduct)

	})
	return r
}

func main() {
	logging.InitializeLogging()
	wg.Add(1)
	go runCronJobs()
	router := registerRoutes()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	handler := cors.Default().Handler(router)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Errorf("Connection could not be established", err)
	}
	wg.Wait()
}
