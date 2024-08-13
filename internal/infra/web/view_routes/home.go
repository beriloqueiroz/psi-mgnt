package routes_view

import (
	"fmt"
	"github.com/beriloqueiroz/psi-mgnt/internal/infra/web/view_routes/components"
	"net/http"
)

type HomeRouteView struct {
}

func NewHomeRouteView() *HomeRouteView {
	return &HomeRouteView{}
}

func (cr *HomeRouteView) HandlerGet(w http.ResponseWriter, r *http.Request) {
	homePage := components.Index()
	err := homePage.Render(r.Context(), w)
	if err != nil {
		fmt.Println(err)
		return
	}
}
