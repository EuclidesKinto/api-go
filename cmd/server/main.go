package main

import (
	"api-go/configs"
	"api-go/internal/dto"
	"api-go/internal/entity"
	"api-go/internal/infra/database"
	"encoding/json"
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
	productHandler := NewProductHandler(productDB)
	http.HandleFunc("/products", productHandler.CreateProduct)
	fmt.Println("Servidor iniciando na porta 8000...")
	if errHttp := http.ListenAndServe(":8000", nil); errHttp != nil {
		log.Fatalf("Falha ao iniciar o servidor HTTP: %v", errHttp)
	}

}

type ProductHandler struct {
	ProductDB database.ProductInterface
}

func NewProductHandler(db database.ProductInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = h.ProductDB.Create(p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
