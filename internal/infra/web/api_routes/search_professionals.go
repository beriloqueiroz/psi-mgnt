package api_routes

import (
	"encoding/json"
	"net/http"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
)

type SearchProfessionalsRoute struct {
	searchProfessionalsUseCase application.SearchProfessionalsUseCase
}

func NewSearchProfessionalsRoute(searchProfessionalsUseCase application.SearchProfessionalsUseCase) *SearchProfessionalsRoute {
	return &SearchProfessionalsRoute{searchProfessionalsUseCase}
}

func (cr *SearchProfessionalsRoute) Handler(w http.ResponseWriter, r *http.Request) {
	term := r.URL.Query().Get("term")
	input := application.SearchProfessionalsInputDTO{
		Term: term,
	}
	output, err := cr.searchProfessionalsUseCase.Execute(r.Context(), input)

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
