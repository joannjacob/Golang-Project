package main

import (
	"fmt"
	"net/http"
	"os"
	"product_management/auth"
	"product_management/logging"
	"product_management/product"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
)

var db *gorm.DB

func setupDB() {

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&auth.Account{}, &product.Product{})

}

func registerRoutes(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(auth.JwtAuthentication)

	pr := product.NewProductRepository(db)
	pu := product.NewProductUsecase(db, pr)
	pHandler := product.NewProductHandler(db, pu)

	ar := auth.NewAuthRepository(db)
	au := auth.NewAuthUsecase(db, ar)
	aHandler := auth.NewAuthHandler(db, au)

	r.Route("/api", func(r chi.Router) {

		r.Post("/signup", aHandler.CreateAccount)
		r.Post("/login", aHandler.Authenticate)
		r.Post("/token/refresh", aHandler.Refresh)

		r.Get("/get_product", pHandler.GetProducts)
		r.Get("/get_product/{id}", pHandler.GetProductById)
		r.Post("/add_product", pHandler.CreateProduct)
		r.Put("/update_product/{id}", pHandler.UpdateProduct)
		r.Delete("/delete_product/{id}", pHandler.DeleteProduct)

	})
	return r
}

func main() {
	err := godotenv.Load("./config/local.env")
	if err != nil {
		log.WithFields(log.Fields{"method": "DB connection init()", "error": err}).Error("Error loading local.env file")
	}

	logging.InitializeLogging()

	setupDB()

	router := registerRoutes(db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8002"
	}

	handler := cors.Default().Handler(router)
	err = http.ListenAndServe(":"+port, handler)
	if err != nil {
		fmt.Print(err)
	}
}
