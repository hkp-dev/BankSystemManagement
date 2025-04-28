package transaction

import (
	"app/be/internal/utility"
	"net/http"
)

type TransactionHandler struct {
	transService *TransactionService
}

func NewTransactionHandler(transService *TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transService: transService,
	}
}

func (ins *TransactionHandler) GetAllTransaction(w http.ResponseWriter, r *http.ResponseWriter) {
	utility.MapCor(w)
	
}
