package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/JMKobayashi/GoExpert-ClearArchitecture/internal/usecase"
)

type ListOrdersHandler struct {
	ListOrdersUseCase *usecase.ListOrdersUseCase
}

func NewListOrdersHandler(listOrdersUseCase *usecase.ListOrdersUseCase) *ListOrdersHandler {
	return &ListOrdersHandler{
		ListOrdersUseCase: listOrdersUseCase,
	}
}

func (h *ListOrdersHandler) Handle(w http.ResponseWriter, r *http.Request) {
	output, err := h.ListOrdersUseCase.Execute()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
