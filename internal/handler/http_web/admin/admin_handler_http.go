package admin

import (
	"FGW_WEB/internal/model"
	"html/template"
	"net/http"
)

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/admin.html",
		"templates/performers.html", // подключаем как частичный шаблон
	))

	data := struct {
		PerformerId   string
		PerformerRole string
		PerformerFIO  string
		CurrentPage   string
		// данные для performers
		Title      string
		Performers []model.Performer
		Roles      []model.Role
	}{
		CurrentPage: "welcome",
		// заполняем остальные поля
	}

	tmpl.Execute(w, data)
}

func PerformersHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/admin.html",
		"templates/performers.html",
	))

	data := struct {
		PerformerId   string
		PerformerRole string
		PerformerFIO  string
		CurrentPage   string
		Title         string
		Performers    []model.Performer
		Roles         []model.Role
	}{
		CurrentPage: "performers",
		Title:       "Сотрудники",
		// заполняем Performers и Roles
	}

	tmpl.Execute(w, data)
}
