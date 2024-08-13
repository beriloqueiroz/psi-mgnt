package routes_view

import (
	"errors"
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	"html/template"
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
	ownerId := r.PathValue("ownerId")
	input := application.ListSessionsInputDto{
		PageSize: pageSizeInt,
		Page:     pageInt,
		OwnerId:  ownerId,
	}

	output, err := cr.listSessionsUseCase.Execute(r.Context(), input)

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/sessions/%s?error_msg=%s", input.OwnerId, msg), http.StatusMovedPermanently)
		return
	}

	tmpl := template.Must(template.ParseFiles("internal/infra/web/view_routes/templates/sessoes_list.html"))
	err = tmpl.Execute(w, output)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (cr *ListSessionRouteView) HandlerPost(w http.ResponseWriter, r *http.Request) {
	var input application.DeleteSessionInputDTO

	var err error
	input.ID = r.PathValue("id")
	input.OwnerId = r.PathValue("ownerId")

	if input.ID == "" || input.OwnerId == "" {
		err = errors.New("invalid id or ownerId")
	}

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/sessions/%s?error_msg=%s", input.OwnerId, msg), http.StatusMovedPermanently)
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
		http.Redirect(w, r, fmt.Sprintf("/sessions/%s?error_msg=%s", input.OwnerId, msg), http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/sessions/%s?delete_success=true", input.OwnerId), http.StatusMovedPermanently)
	return
}