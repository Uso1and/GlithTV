package handler

import (
	"glitchTV/internal/session"
	"html/template"
	"net/http"
	"strings"
)

// Кастомные функции для шаблонов
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"substr": func(s string, start, length int) string {
			if start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"upper": strings.ToUpper,
		"firstChar": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return strings.ToUpper(string(s[0]))
		},
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("").Funcs(templateFuncs())
	tmpl, err := tmpl.ParseFiles(
		"templates/base.html",
		"templates/header.html",
		"templates/index.html",
	)
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем данные сессии
	sess, _ := session.Store.Get(r, "session-name")
	authenticated := false
	username := ""

	if auth, ok := sess.Values["authenticated"].(bool); ok {
		authenticated = auth
	}
	if name, ok := sess.Values["username"].(string); ok {
		username = name
	}

	data := struct {
		Authenticated bool
		Username      string
	}{
		Authenticated: authenticated,
		Username:      username,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {
	tmpl := template.New("").Funcs(templateFuncs())
	tmpl, err := tmpl.ParseFiles(
		"templates/base.html",
		"templates/header.html",
		"templates/profile.html",
	)
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблонов: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем данные сессии
	sess, _ := session.Store.Get(r, "session-name")
	authenticated := false
	username := ""

	if auth, ok := sess.Values["authenticated"].(bool); ok {
		authenticated = auth
	}
	if name, ok := sess.Values["username"].(string); ok {
		username = name
	}

	data := struct {
		Authenticated bool
		Username      string
	}{
		Authenticated: authenticated,
		Username:      username,
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		http.Error(w, "Ошибка рендеринга шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}
