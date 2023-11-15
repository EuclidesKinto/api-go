package main

import (
	"api-go/configs"
	"api-go/internal/entity"
	"api-go/internal/infra/database"
	"api-go/internal/infra/webserver/handlers"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Falha ao abrir o configs: %v", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao abrir o banco de dados: %v", err)
	}

	if errDb := db.AutoMigrate(&entity.Product{}, &entity.User{}); errDb != nil {
		log.Fatalf("Falha na migração do banco de dados: %v", errDb)
	}

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, config.TokenAuth, config.JWTExpiresIn)

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	route := chi.NewRouter()
	route.Use(middleware.Logger)

	route.Post("/users", userHandler.Create)
	route.Post("/users/jwt", userHandler.GetJWT)

	route.Route("/products", func(route chi.Router) {
		route.Use(jwtauth.Verifier(config.TokenAuth))
		route.Use(jwtauth.Authenticator)
		route.Post("/", productHandler.CreateProduct)
		route.Get("/", productHandler.GetProducts)
		route.Get("/{id}", productHandler.GetProduct)
		route.Put("/{id}", productHandler.UpdateProduct)
		route.Delete("/{id}", productHandler.Delete)
	})

	fmt.Println("Servidor iniciando na porta 8000...")
	if errHttp := http.ListenAndServe(":8000", route); errHttp != nil {
		log.Fatalf("Falha ao iniciar o servidor HTTP: %v", errHttp)
	}

}
