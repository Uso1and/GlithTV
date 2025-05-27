package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"glitchTV/internal/database"
	"glitchTV/internal/session"
	"html/template"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/register.html")
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	hashedPassword := hashPassword(password)

	_, err := database.DB.Exec(
		"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
		username, email, hashedPassword)
	if err != nil {
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "templates/login.html")
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	var storedHash string
	err := database.DB.QueryRow(
		"SELECT password_hash FROM users WHERE username = $1", username).Scan(&storedHash)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if hashPassword(password) != storedHash {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Создаём сессию
	sess, _ := session.Store.Get(r, "session-name")
	sess.Values["authenticated"] = true
	sess.Values["username"] = username
	sess.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.Store.Get(r, "session-name")
	sess.Values["authenticated"] = false
	delete(sess.Values, "username")
	sess.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess, _ := session.Store.Get(r, "session-name")

		// Разрешаем доступ к публичным страницам без аутентификации
		if r.URL.Path == "/login" || r.URL.Path == "/register" || r.URL.Path == "/login-form" || r.URL.Path == "/register-form" {
			next(w, r)
			return
		}

		if auth, ok := sess.Values["authenticated"].(bool); !ok || !auth {
			if r.URL.Path == "/" {
				next(w, r)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next(w, r)
	}
}

func ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/login.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки формы входа", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ShowRegisterForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/register.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки формы регистрации", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
