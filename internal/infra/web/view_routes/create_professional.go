package routes_view

import (
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/internal/application"
	"html/template"
	"net/http"
)

type CreateProfessionalRouteView struct {
	CreateProfessionalUseCase application.CreateProfessionalUseCase
}

func NewCreateProfessionalRouteView(createProfessionalUseCase application.CreateProfessionalUseCase) *CreateProfessionalRouteView {
	return &CreateProfessionalRouteView{createProfessionalUseCase}
}

func (cr *CreateProfessionalRouteView) HandlerGet(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("internal/infra/web/view_routes/templates/profissional_form.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func (cr *CreateProfessionalRouteView) HandlerPost(w http.ResponseWriter, r *http.Request) {
	var input application.CreateProfessionalInputDTO

	var err error

	input.Email = r.FormValue("email")
	input.Document = r.FormValue("documento")
	input.Phone = r.FormValue("telefone")
	input.Name = r.FormValue("nome")

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/professional?error_msg=%s", msg), http.StatusMovedPermanently)
		return
	}

	fmt.Println(input)

	_, err = cr.CreateProfessionalUseCase.Execute(r.Context(), input)

	if err != nil {
		msg := struct {
			Message string `json:"message"`
		}{
			Message: err.Error(),
		}
		http.Redirect(w, r, fmt.Sprintf("/professional?error_msg=%s", msg), http.StatusMovedPermanently)
		return
	}

	http.Redirect(w, r, "/sessions?created_success=true", http.StatusMovedPermanently)
}
