package main

import (
	"fmt"
	// "html/template"
	"net/http"
	"wfi/internal/handlers"
)

func main() {
	handlers.RegisterHandlers()
	handlers.DbStart()

    // Запуск сервера на порту 8080
    fmt.Println("Сервер запущен на http://localhost:8080")
    http.ListenAndServe(":8080", nil)

}
