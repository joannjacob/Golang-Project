package router

import (
	"net/http"

	"product_management/auth"
	"product_management/product"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jinzhu/gorm"
)

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
