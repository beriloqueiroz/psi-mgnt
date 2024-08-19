package routes_view

import (
	"fmt"
	domain "github.com/beriloqueiroz/psi-mgnt/internal/domain/entity"
	"github.com/beriloqueiroz/psi-mgnt/internal/infra/web/view_routes/components"
	"net/http"
	"strconv"
	"time"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
)

type SessionRouteView struct {
	CreateSessionUseCase application.CreateSessionUseCase
}

func NewSessionRouteView(createSessionUseCase application.CreateSessionUseCase) *SessionRouteView {
	return &SessionRouteView{createSessionUseCase}
}

func (cr *SessionRouteView) HandlerGet(w http.ResponseWriter, r *http.Request) {
	session := domain.Session{ID: "", Patient: &domain.Patient{ID: "111", Name: "berilo ttt"},
		Professional: &domain.Professional{ID: "111233", Name: "ric"},
		Duration:     time.Hour,
		Notes:        "aiaiai aiaiai",
		Plan:         "PARTICULAR",
		Price:        50,
		Date:         time.Now(),
	}
	if session.ID == "" {
		session = domain.Session{}
		session.Date = time.Now()
		session.Duration = time.Minute * 50
		session.Price = 50
		session.Patient = &domain.Patient{}
		session.Professional = &domain.Professional{}
	}
	h := components.SessionForm(session)
	err := h.Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (cr *SessionRouteView) HandlerPost(w http.ResponseWriter, r *http.Request) {
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