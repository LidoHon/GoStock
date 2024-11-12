package controller

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/LidoHon/GoStock/database"
	"github.com/LidoHon/GoStock/models"
	"github.com/gorilla/mux"
)

// Response struct to send responses with status and message
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// CreateStock creates a new stock record
func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	// Decode request body into stock model
	err := json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		http.Error(w, "Unable to decode the request body", http.StatusBadRequest)
		log.Println("Unable to decode the request body:", err)
		return
	}

	// Insert the stock into the database
	insertID := database.InsertStock(stock)

	// Create a response with the stock ID
	res := response{
		ID:      insertID,
		Message: "Stock created successfully",
	}

	// Encode response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// GetStock retrieves a stock by ID
func GetStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Unable to get stock ID", http.StatusBadRequest)
		log.Println("Unable to get stock ID:", err)
		return
	}

	stock, err := database.GetStock(int64(id))
	if err != nil {
		http.Error(w, "Unable to retrieve stock", http.StatusInternalServerError)
		log.Println("Unable to retrieve stock:", err)
		return
	}

	// Encode the stock details as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stock)
}

// GetAllStock retrieves all stocks
func GetAllStock(w http.ResponseWriter, r *http.Request) {
	stocks, err := database.GetAllStock()
	if err != nil {
		http.Error(w, "Unable to retrieve stocks", http.StatusInternalServerError)
		log.Println("Unable to retrieve stocks:", err)
		return
	}

	// Encode the list of stocks as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stocks)
}

// DeleteStock deletes a stock by ID
func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Unable to get stock ID", http.StatusBadRequest)
		log.Println("Unable to get stock ID:", err)
		return
	}

	deletedRows := database.DeleteStock(int64(id))
	if deletedRows == 0 {
		http.Error(w, "No stock found to delete", http.StatusNotFound)
		log.Println("No stock found to delete")
		return
	}

	msg := fmt.Sprintf("Stock deleted successfully. Total rows affected: %v", deletedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// Encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

// UpdateStock updates an existing stock by ID
func UpdateStock(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "Unable to get stock ID", http.StatusBadRequest)
		log.Println("Unable to get stock ID:", err)
		return
	}

	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		http.Error(w, "Unable to decode the request body", http.StatusBadRequest)
		log.Println("Unable to decode the request body:", err)
		return
	}

	// Update the stock in the database
	updatedRows := database.UpdateStock(int64(id), stock)
	if updatedRows == 0 {
		http.Error(w, "No stock found to update", http.StatusNotFound)
		log.Println("No stock found to update")
		return
	}

	msg := fmt.Sprintf("Stock updated successfully. Total rows affected: %v", updatedRows)
	res := response{
		ID:      int64(id),
		Message: msg,
	}

	// Encode the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}
