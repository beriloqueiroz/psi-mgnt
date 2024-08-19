package routes_view

import (
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/internal/infra/web/view_routes/components"
	"github.com/beriloqueiroz/psi-mgnt/pkg/helpers"
	"net/http"
	"strconv"
	"time"

	"github.com/beriloqueiroz/psi-mgnt/internal/application"
)

type SessionRouteView struct {
	CreateSessionUseCase application.CreateSessionUseCase
	UpdateSessionUseCase application.UpdateSessionUseCase
	FindSessionUseCase   application.FindSessionUseCase
}

func NewSessionRouteView(createSessionUseCase application.CreateSessionUseCase,
	updateSessionUseCase application.UpdateSessionUseCase,
	findSessionUseCase application.FindSessionUseCase) *SessionRouteView {
	return &SessionRouteView{createSessionUseCase, updateSessionUseCase, findSessionUseCase}
}

func (cr *SessionRouteView) HandlerGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var session application.FindSessionOutputDto
	var err error
	if id != "" {
		session, err = cr.FindSessionUseCase.Execute(r.Context(), application.FindSessionInputDto{
			ID: id,
		})
	} else {
		session.Date = time.Now()
		session.Duration = time.Minute * 50
		session.Price = 50
	}
	h := components.SessionForm(session)
	err = h.Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (cr *SessionRouteView) HandlerPost(w http.ResponseWriter, r *http.Request) {
	var input application.CreateSessionInputDTO
	var err error
	id := r.FormValue("session_id")
	input.Notes = r.FormValue("notas")
	if id != "" {
		inputUpdate := application.UpdateSessionInputDTO{
			ID:    id,
			Notes: input.Notes,
		}
		_, err = cr.UpdateSessionUseCase.Execute(r.Context(), inputUpdate)
	} else {
		input.PatientName = r.FormValue("paciente_nome")
		input.PatientId = r.FormValue("paciente_id")
		year, month, day, hour, mi, _, err := helpers.DecomposeStringDate(r.FormValue("data_hora") + ":00")
		input.Date = time.Date(year, time.Month(month), day, hour, mi, 0, 0, time.UTC)
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
		_, err = cr.CreateSessionUseCase.Execute(r.Context(), input)
	}

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
