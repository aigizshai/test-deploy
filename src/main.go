package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	log.Println("Сервер запущен на http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	resp := response{Message: "its works"}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}
