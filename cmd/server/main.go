package main

import (
	"api-go/configs"
	"api-go/internal/entity"
	"api-go/internal/infra/database"
	"api-go/internal/infra/webserver/handlers"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	if _, err := configs.LoadConfig("."); err != nil {
		log.Fatalf("Falha ao carregar a configuração: %v", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Falha ao abrir o banco de dados: %v", err)
	}

	if errDb := db.AutoMigrate(&entity.Product{}, &entity.User{}); errDb != nil {
		log.Fatalf("Falha na migração do banco de dados: %v", errDb)
	}

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB)

	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	route := chi.NewRouter()
	route.Use(middleware.Logger)

	route.Post("/users", userHandler.Create)

	route.Post("/products", productHandler.CreateProduct)
	route.Get("/products", productHandler.GetProducts)
	route.Get("/products/{id}", productHandler.GetProduct)
	route.Put("/products/{id}", productHandler.UpdateProduct)
	route.Delete("/products/{id}", productHandler.Delete)

	fmt.Println("Servidor iniciando na porta 8000...")
	if errHttp := http.ListenAndServe(":8000", route); errHttp != nil {
		log.Fatalf("Falha ao iniciar o servidor HTTP: %v", errHttp)
	}

}
