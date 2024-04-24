package handlers

import (
	"html/template"
	"net/http"
	"path/filepath"
)

func MainPageHandler(w http.ResponseWriter, r *http.Request) {
	//указываем путь к нужному файлу
	path := filepath.Join("webPages", "main.html")
	//создаем html-шаблон
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	//выводим шаблон клиенту в браузер
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}
