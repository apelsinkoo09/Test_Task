package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"test_task/internal/model"
)

func UpdateBalanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			log.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var updateBalance model.TransactionRequest

		if err := json.NewDecoder(r.Body).Decode(&updateBalance); err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		if updateBalance.Amount <= 0 {
			http.Error(w, "Invalid amount", http.StatusBadRequest)
			return
		}
		if strings.ToUpper(updateBalance.OperationType) != "DEPOSIT" && strings.ToUpper(updateBalance.OperationType) != "WITHDRAW" {

			http.Error(w, "Invalid operation", http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
			return
		}
		defer tx.Rollback() // rollback transaction

		// Check balance
		if strings.ToUpper(updateBalance.OperationType) == "WITHDRAW" {
			currentBalance, err := model.GetBalance(db, updateBalance.ValletID)
			if err != nil {
				http.Error(w, "Wallet not found", http.StatusNotFound)
				return
			}
			if currentBalance < updateBalance.Amount {
				http.Error(w, "Шnsufficient funds", http.StatusBadRequest)
				return
			}
		}

		if err := model.UpdateBalance(tx, updateBalance); err != nil {
			http.Error(w, "Failed to update balance", http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Balance updated successfully"))
	}
}

func GetBalanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			log.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idParam := r.URL.Query().Get("valletId")
		valletId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		balance, err := model.GetBalance(db, valletId)
		if err != nil {
			http.Error(w, "Invalid id", http.StatusBadRequest)
			return
		}

		resp := map[string]interface{}{
			"valletID": valletId,
			"balance":  balance,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		//w.Write([]byte(fmt.Sprintf("%f", balance)))
	}
}
func GetUUIDBalanceHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idParam := strings.TrimPrefix(r.URL.Path, "/api/v1/wallets/")
		valletId, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			http.Error(w, "Invalid uuid", http.StatusBadRequest)
			return
		}

		balance, err := model.GetBalance(db, valletId)
		if err != nil {

			http.Error(w, "Wallet not found", http.StatusNotFound)
			return
		}

		resp := map[string]interface{}{
			"uuid":    valletId,
			"balance": balance,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
