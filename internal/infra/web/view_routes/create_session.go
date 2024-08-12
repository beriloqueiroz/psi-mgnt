package routes_view

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
)

type CreateSessionRouteView struct {
	CreateSessionUseCase application.CreateSessionUseCase
}

func NewCreateSessionRouteView(createSessionUseCase application.CreateSessionUseCase) *CreateSessionRouteView {
	return &CreateSessionRouteView{createSessionUseCase}
}

func (cr *CreateSessionRouteView) HandlerGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("internal/infra/web/view_routes/templates/sessao_form.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func (cr *CreateSessionRouteView) HandlerPost(w http.ResponseWriter, r *http.Request) {
	var input application.CreateSessionInputDTO

	var err error

	fmt.Println(r.FormValue("data_hora"))
	fmt.Println(r.FormValue("duracao"))

	input.PatientName = r.FormValue("paciente_nome")
	input.Date, err = time.Parse("2006-01-02T15:04", r.FormValue("data_hora"))
	input.Notes = r.FormValue("notas")
	input.Duration, err = time.ParseDuration(r.FormValue("duracao"))
	input.Price, err = strconv.ParseFloat(r.FormValue("preco"), 64)
	input.OwnerId = r.FormValue("dono")

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/sessions/%s?error_msg=%s", input.OwnerId, msg), http.StatusMovedPermanently)
		return
	}

	fmt.Println(input)

	_, err = cr.CreateSessionUseCase.Execute(r.Context(), input)

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/sessions/%s?error_msg=%s", input.OwnerId, msg), http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, "/sessions/"+input.OwnerId+"?created_success=true", http.StatusMovedPermanently)
}
