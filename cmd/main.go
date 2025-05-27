package main

import (
	"glitchTV/internal/auth"
	"glitchTV/internal/database"
	"glitchTV/internal/handler"
	"glitchTV/internal/session"
	"net/http"
)

func main() {
	// Инициализация базы данных
	database.InitDB()

	// Инициализация хранилища сессий
	session.InitSessionStore("your-very-secret-key") // Замените на реальный секретный ключ

	// Статические файлы
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Роуты
	http.HandleFunc("/", auth.AuthMiddleware(handler.Home))
	http.HandleFunc("/login", auth.LoginHandler)
	http.HandleFunc("/register", auth.RegisterHandler)
	http.HandleFunc("/logout", auth.LogoutHandler)
	http.HandleFunc("/profile", auth.AuthMiddleware(handler.Profile))

	// Новые маршруты для форм
	http.HandleFunc("/login-form", auth.ShowLoginForm)
	http.HandleFunc("/register-form", auth.ShowRegisterForm)

	// Запуск сервера
	println("Сервер запущен на http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
