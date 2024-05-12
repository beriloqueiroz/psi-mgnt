package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
)

type ListSessionsRoute struct {
	listSessionsUseCase application.ListSessionsUseCase
}

func NewListSessionsRoute(listSessionsUseCase application.ListSessionsUseCase) *ListSessionsRoute {
	return &ListSessionsRoute{listSessionsUseCase}
}

func (cr *ListSessionsRoute) Handler(w http.ResponseWriter, r *http.Request) {
	pageSize := r.URL.Query().Get("pageSize")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 50
	}
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	ownerId := r.URL.Query().Get("ownerId")
	input := application.ListSessionsInputDto{
		PageSize: pageSizeInt,
		Page:     pageInt,
		OwnerId:  ownerId,
	}

	output, err := cr.listSessionsUseCase.Execute(r.Context(), input)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(msg)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
