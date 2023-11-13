package main

import (
	"api-go/configs"
	"api-go/internal/entity"
	"api-go/internal/infra/database"
	"api-go/internal/infra/webserver/handlers/produtc_handlers"
	"fmt"
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

	productDB := database.NewProduct(db)
	productHandler := produtc_handlers.NewProductHandler(productDB)
	http.HandleFunc("/products", productHandler.CreateProduct)
	fmt.Println("Servidor iniciando na porta 8000...")
	if errHttp := http.ListenAndServe(":8000", nil); errHttp != nil {
		log.Fatalf("Falha ao iniciar o servidor HTTP: %v", errHttp)
	}

}
