package routes_view

import (
	"errors"
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	"github.com/beriloqueiroz/psi-mgnt/internal/infra/web/view_routes/components"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"
	"net/http"
	"strconv"
)

type ListSessionRouteView struct {
	listSessionsUseCase  application.ListSessionsUseCase
	deleteSessionUseCase application.DeleteSessionUseCase
}

func NewListSessionRouteView(
	listSessionsUseCase application.ListSessionsUseCase,
	deleteSessionUseCase application.DeleteSessionUseCase) *ListSessionRouteView {
	return &ListSessionRouteView{listSessionsUseCase, deleteSessionUseCase}
}

func (cr *ListSessionRouteView) HandlerGet(w http.ResponseWriter, r *http.Request) {
	pageSize := r.URL.Query().Get("pageSize")
	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 50
		err = nil
	}
	page := r.URL.Query().Get("page")
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		pageInt = 1
		err = nil
	}
	professionalId := r.PathValue("professionalId")
	listConfig := helpers.ListConfig{
		PageSize: pageSizeInt,
		Page:     pageInt,
	}
	input := application.ListSessionsInputDto{
		ListConfig:     listConfig,
		ProfessionalId: professionalId,
	}

	output, err := cr.listSessionsUseCase.Execute(r.Context(), input)

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/sessions?error_msg=%s", msg), http.StatusMovedPermanently)
		return
	}

	h := components.SessionsForm(output)
	err = h.Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (cr *ListSessionRouteView) HandlerPost(w http.ResponseWriter, r *http.Request) {
	var input application.DeleteSessionInputDTO

	var err error
	input.ID = r.PathValue("id")

	if input.ID == "" {
		err = errors.New("invalid id or ownerId")
	}

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/sessions?error_msg=%s", msg), http.StatusMovedPermanently)
		return
	}

	_, err = cr.deleteSessionUseCase.Execute(r.Context(), input)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/sessions?error_msg=%s", msg), http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/sessions?delete_success=%s", "true"), http.StatusMovedPermanently)
	return
}
