package transaction

import (
	"app/be/internal/utility"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type TransactionHandler struct {
	transService *TransactionService
}

func NewTransactionHandler(transService *TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transService: transService,
	}
}

func (ins *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	utility.MapCor(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReqGetTransactions
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if req.UserId == "" {
		log.Printf("Missing userId in request body")
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	transactions, total, err := ins.transService.GetTransactions(r.Context(), req.UserId)
	if err != nil {
		log.Printf("Failed to get transactions for userId %s: %v", req.UserId, err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transactions": transactions,
		"total":        total,
	})
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (ins *TransactionHandler) PostTransaction(w http.ResponseWriter, r *http.Request) {
	utility.MapCor(w)
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}
	if r.Method != http.MethodPost {
		log.Printf("Invalid method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReqPostTransaction
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request payload"})
		return
	}

	if req.UserId == "" || req.SenderAccountNumber == "" || req.RecipientAccountNumber == "" || req.OTP == "" {
		log.Printf("Missing required fields: %+v", req)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Missing required fields"})
		return
	}

	transId, createTime, err := ins.transService.PostTransaction(r.Context(), req.UserId, req.SenderAccountNumber, req.RecipientAccountNumber, req.Amount, req.Description, req.OTP)
	if err != nil {
		log.Printf("Failed to process transaction: %v", err)
		w.Header().Set("Content-Type", "application/json")
		errorMsg := err.Error()
		switch {
		case strings.Contains(errorMsg, "user not found"):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Sender account not found"})
		case strings.Contains(errorMsg, "recipient user not found"):
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Recipient account not found"})
		case strings.Contains(errorMsg, "sender account does not belong to user"):
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Sender account does not belong to user"})
		case strings.Contains(errorMsg, "2FA not setup for user"):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "2FA not setup for user"})
		case strings.Contains(errorMsg, "invalid OTP code"):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid OTP code"})
		case strings.Contains(errorMsg, "insufficient balance"):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Insufficient balance"})
		case strings.Contains(errorMsg, "invalid amount"):
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid amount"})
		default:
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Internal server error"})
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transactionId":          transId,
		"senderAccountNumber":    req.SenderAccountNumber,
		"recipientAccountNumber": req.RecipientAccountNumber,
		"amount":                 req.Amount,
		"description":            req.Description,
		"date":                   createTime,
	})
}
