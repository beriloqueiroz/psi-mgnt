package routes

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
)

type DeleteSessionRoute struct {
	deleteSessionUseCase application.DeleteSessionUseCase
}

func NewDeleteSessionRoute(deleteSessionUseCase application.DeleteSessionUseCase) *DeleteSessionRoute {
	return &DeleteSessionRoute{deleteSessionUseCase}
}

func (cr *DeleteSessionRoute) Handler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	input := application.DeleteSessionInputDTO{
		ID: id,
	}

	output, err := cr.deleteSessionUseCase.Execute(r.Context(), input)

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
