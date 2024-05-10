package routes

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
)

type SearchPatientsRoute struct {
	searchPatientsUseCase application.SearchPatientsUseCase
}

func NewSearchPatientsRoute(searchPatientsUseCase application.SearchPatientsUseCase) *SearchPatientsRoute {
	return &SearchPatientsRoute{searchPatientsUseCase}
}

func (cr *SearchPatientsRoute) Handler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	input := application.SearchPatientsInputDTO{
		Term: term,
	}

	output, err := cr.searchPatientsUseCase.Execute(r.Context(), input)

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
