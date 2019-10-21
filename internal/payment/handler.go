package payment

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/logger"
)

func AllPaymentsForClient(s Service) func(http.ResponseWriter,
	*http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		log := logger.FromContext(ctx)
		level.Debug(log).Log("msg", "get all payments request")
		forAccount := request.URL.Query().Get("forAccount")
		if len(forAccount) == 0 {
			jsonResponse(writer, map[string]string{"message": "forAccount is required parameter"}, http.StatusBadRequest)
			return
		}
		payments, err := s.GetAllForClient(ctx, forAccount)
		if err != nil {
			jsonResponse(writer, map[string]string{"message": "something went wrong"}, http.StatusBadRequest)
			return
		}
		resultList := make([]ResponseItem, len(payments))
		for i, p := range payments {
			resultList[i] = p.ToResponseItem(forAccount)
		}
		jsonResponse(writer, resultList, http.StatusOK)
	}
}

func CreatePayment(s Service) func(http.ResponseWriter,
	*http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		log := logger.FromContext(ctx)
		level.Debug(log).Log("msg", "create payment request")
		var req *CreatePaymentRequest
		err := json.NewDecoder(request.Body).Decode(&req)
		if err != nil {
			jsonResponse(writer, map[string]string{"message": "wrong parameters"}, http.StatusBadRequest)
			return
		}
		if req.Amount <= 0 {
			jsonResponse(writer, map[string]string{"message": "amount should be more than 0"}, http.StatusBadRequest)
			return
		}
		if len(req.AccountFrom) == 0 {
			jsonResponse(writer, map[string]string{"message": "accountFrom is required"}, http.StatusBadRequest)
			return
		}
		if len(req.AccountTo) == 0 {
			jsonResponse(writer, map[string]string{"message": "accountTo is required"}, http.StatusBadRequest)
			return
		}
		if len(req.Currency) == 0 {
			jsonResponse(writer, map[string]string{"message": "currency is required"}, http.StatusBadRequest)
			return
		}
		if req.AccountTo == req.AccountFrom {
			jsonResponse(writer, map[string]string{"message": "accountTo and accountFrom must be different"}, http.StatusBadRequest)
			return
		}
		err = s.CreatePayment(ctx, req.AccountFrom, req.AccountTo, req.Amount, req.Currency)
		if err != nil {
			jsonResponse(writer, map[string]string{"message": "something went wrong"}, http.StatusBadRequest)
			return
		}
		jsonResponse(writer, nil, http.StatusOK)
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, string(jsonData))
	}
}
