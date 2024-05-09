package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nicjtutorial/data"
	"strconv"

	"github.com/gorilla/mux"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(rw http.ResponseWriter, r *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("Handle Post")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(rw http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	stringId := vars["id"]
	id, err := strconv.Atoi(stringId)
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}
	p.l.Println("Handle Put", id)
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrProductNotFound {
		http.Error(rw, "ProductNotFound", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(rw, "ProductNotFound", http.StatusInternalServerError)
		return
	}
}

type KeyProduct struct {}

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		prod := &data.Product{}
		err := prod.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		}

		// validate product
		err = prod.Validate()
		if err != nil {
			p.l.Println("Error validating product")
			http.Error(
				rw, 
				fmt.Sprintf("Error validating product: %s", err), 
				http.StatusBadRequest,
			)
			return
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, *prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(rw, req)
	})
}