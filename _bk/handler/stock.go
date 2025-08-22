package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	config "go-lab/config"
	models "go-lab/model"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetStockById(w http.ResponseWriter, r *http.Request) {
	// ตรวจสอบว่าเป็น GET request
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// ดึง id จาก path parameter (/stock/{id})
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing stock ID", http.StatusBadRequest)
		return
	}

	// แปลง id เป็น uint
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	// ดึง stock จาก DB
	var stock models.Stock
	if result := config.DB.First(&stock, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			http.Error(w, "Stock not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	// ส่ง response เป็น JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

func InsertStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	stock := models.Stock{
		Balance: 1000000.0,
		Reserve: 0.0,
		OnHand:  1000000.0,
	}

	log.Printf("Inserting stock: %+v", stock)
	if result := config.DB.Create(&stock); result.Error != nil {
		log.Printf("DB error: %v", result.Error)
		http.Error(w, "Database Error: "+result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(stock)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {
	log.Printf("Updating stock request received")

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// อ่าน query params
	idStr := r.URL.Query().Get("id")
	qtyStr := r.URL.Query().Get("qty")
	if idStr == "" || qtyStr == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	qty, err := strconv.ParseFloat(qtyStr, 64)
	if err != nil || qty < 0 {
		http.Error(w, "Invalid qty value", http.StatusBadRequest)
		return
	}

	// เริ่ม transaction
	tx := config.DB.Begin()
	if tx.Error != nil {
		http.Error(w, "Database transaction error: "+tx.Error.Error(), http.StatusInternalServerError)
		return
	}

	// defer rollback ป้องกัน panic / error
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Panic recovered, transaction rolled back: %v", r)
		}
	}()

	var stock models.Stock
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&stock, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Stock not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database Error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	newReserve := stock.Reserve + qty
	newOnHand := stock.Balance - newReserve

	if err := tx.Model(&models.Stock{}).Where("id = ?", id).Updates(map[string]interface{}{
		"reserve": newReserve,
		"on_hand": newOnHand,
	}).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Database Update Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit().Error; err != nil {
		http.Error(w, "Database Commit Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// อ่านค่าล่าสุด
	if err := config.DB.First(&stock, id).Error; err != nil {
		http.Error(w, "Database Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(stock); err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}

	log.Printf("Updated stock ID %d: Reserve=%f, OnHand=%f, Balance=%f", stock.ID, stock.Reserve, stock.OnHand, stock.Balance)
}
