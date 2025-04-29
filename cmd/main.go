package main

import (
	"log"
	"net/http"
	"voucher-api/internal/config"
	"voucher-api/internal/database"
	"voucher-api/internal/handler"
	"voucher-api/internal/repository"
	"voucher-api/internal/service"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()
	db := database.InitDB(cfg)
	defer db.Close()

	repo := repository.NewRepository(db)
	svc := service.NewService(repo)
	h := handler.NewHandler(svc)

	router := mux.NewRouter()
	router.HandleFunc("/brand", h.CreateBrand).Methods("POST")
	router.HandleFunc("/voucher", h.CreateVoucher).Methods("POST")
	router.HandleFunc("/voucher", h.GetVoucher).Methods("GET")
	router.HandleFunc("/voucher/brand", h.GetVouchersByBrand).Methods("GET")
	router.HandleFunc("/transaction/redemption", h.MakeRedemption).Methods("POST")
	router.HandleFunc("/transaction/redemption", h.GetTransactionDetail).Methods("GET")
	router.HandleFunc("/customer", h.CreateCustomer).Methods("POST")

	log.Printf("Server starting on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
