package routes_view

import (
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/internal/infra/web/view_routes/components"
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
	h := components.SessionForm()
	err := h.Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (cr *CreateSessionRouteView) HandlerPost(w http.ResponseWriter, r *http.Request) {
	var input application.CreateSessionInputDTO

	var err error

	input.PatientName = r.FormValue("paciente_nome")
	input.PatientId = r.FormValue("paciente_id")
	input.Date, err = time.Parse("2006-01-02T15:04", r.FormValue("data_hora"))
	input.Notes = r.FormValue("notas")
	input.Plan = r.FormValue("plano")
	input.Duration, err = time.ParseDuration(r.FormValue("duracao"))
	input.Price, err = strconv.ParseFloat(r.FormValue("preco"), 64)
	input.ProfessionalId = r.FormValue("profissional_id")

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/session?error_msg=%s", msg), http.StatusMovedPermanently)
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
		http.Redirect(w, r, fmt.Sprintf("/session?error_msg=%s", msg), http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, "/sessions?created_success=true", http.StatusMovedPermanently)
}
