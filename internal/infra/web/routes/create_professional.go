package routes

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
)

type CreateProfessionalRoute struct {
	CreateProfessionalUseCase application.CreateProfessionalUseCase
}

func NewCreateProfessionalRoute(createProfessionalUseCase application.CreateProfessionalUseCase) *CreateProfessionalRoute {
	return &CreateProfessionalRoute{createProfessionalUseCase}
}

func (cr *CreateProfessionalRoute) Handler(w http.ResponseWriter, r *http.Request) {

	var input application.CreateProfessionalInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(msg)
		return
	}

	output, err := cr.CreateProfessionalUseCase.Execute(r.Context(), input)

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
