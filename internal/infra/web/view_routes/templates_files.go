package routes_view

import "html/template"

const path = "internal/infra/web/view_routes/templates/"

func GetBaseFormTemplates(fileName string) (*template.Template, error) {
	return template.New("").ParseFiles(path+fileName, path+"base.html")
}
