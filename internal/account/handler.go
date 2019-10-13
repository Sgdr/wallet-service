package account

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log/level"
	"github.com/sgdr/wallet-service/internal/logger"
)

func AllAccountsHandler(s Service) func(http.ResponseWriter,
	*http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		ctx := request.Context()
		log := logger.FromContext(ctx)
		level.Debug(log).Log("msg", "get all accounts request")
		accounts, err := s.GetAll(ctx)
		if err != nil {
			jsonResponse(writer, map[string]string{"message": "something went wrong"}, http.StatusBadRequest)
			return
		}
		resultList := make([]ResponseItem, len(accounts))
		for i, a := range accounts {
			resultList[i] = a.ToResponseItem()
		}
		jsonResponse(writer, resultList, http.StatusOK)
	}
}

func jsonResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, string(jsonData))
}
