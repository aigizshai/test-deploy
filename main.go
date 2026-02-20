package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// response структура для формирования JSON-ответа
type response struct {
	Message string `json:"message"`
}

func main() {
	// Регистрируем обработчик для пути /hello
	http.HandleFunc("/hello", helloHandler)

	// Запускаем сервер на порту 8080
	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

// helloHandler обрабатывает запросы к /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, что метод запроса — GET (опционально, для демонстрации)
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Устанавливаем заголовок Content-Type для JSON
	w.Header().Set("Content-Type", "application/json")

	// Создаём экземпляр ответа
	resp := response{Message: "its works"}

	// Кодируем структуру в JSON и отправляем клиенту
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Если произошла ошибка при кодировании, логируем её и отправляем статус 500
		log.Println("Ошибка кодирования JSON:", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}
