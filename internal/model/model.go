package model

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type TransactionRequest struct {
	ValletID      int64   `json:"valletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}

func UpdateBalance(tx *sql.Tx, req TransactionRequest) error {
	query := "Update vallets Set balance = balance + $1 Where vallet = $2"
	amount := req.Amount
	if req.OperationType == "WITHDRAW" {
		amount = -amount
	}

	_, err := tx.ExecContext(context.Background(), query, amount, req.ValletID)
	if err != nil {
		return err
	}
	return nil
}

func GetBalance(db *sql.DB, valletID int64) (float64, error) {
	var balance float64
	query := "Select balance from vallets where vallet = $1"
	err := db.QueryRow(query, valletID).Scan(&balance)
	if err != nil {
		log.Printf("Error fetching balance for valletID %d: %v", valletID, err)
		return 0, fmt.Errorf("wallet not found")
	}
	return balance, nil
}
