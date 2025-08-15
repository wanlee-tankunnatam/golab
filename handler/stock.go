package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	config "go-lab/config"
	models "go-lab/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InsertStock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	stock := models.Stock{
		Balance: 100.0,
		Reserve: 0.0,
		OnHand:  100.0,
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
	log.Printf("Updating stock 1111")
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Printf("Updating stock 2222")
	// อ่าน stock ID จาก query
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing stock ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid stock ID", http.StatusBadRequest)
		return
	}

	// อ่าน qty
	qtyStr := r.URL.Query().Get("qty")
	if qtyStr == "" {
		http.Error(w, "Missing qty parameter", http.StatusBadRequest)
		return
	}
	qty, err := strconv.ParseFloat(qtyStr, 64)
	if err != nil || qty < 0 {
		http.Error(w, "Invalid qty value", http.StatusBadRequest)
		return
	}

	// เริ่ม transaction
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var stock models.Stock

	// Lock row เพื่อป้องกัน race condition
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&stock, id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Stock not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database Error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	// คำนวณค่าใหม่
	log.Printf("stock data Stock ID %d before update: Reserve=%f, OnHand=%f", stock.ID, stock.Reserve, stock.OnHand)
	// return
	newReserve := stock.Reserve + qty
	newOnHand := stock.Balance - newReserve

	// Update ทั้ง reserve และ on_hand ใน query เดียว
	if err := tx.Model(&models.Stock{}).Where("id = ?", id).
		Updates(map[string]interface{}{
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

	// อ่านค่าล่าสุดจาก DB
	if err := config.DB.First(&stock, id).Error; err != nil {
		http.Error(w, "Database Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ส่ง response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stock)

	log.Printf("Updated stock ID %d: Reserve=%f, OnHand=%f, Balance=%f", stock.ID, stock.Reserve, stock.OnHand, stock.Balance)
}
