package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	// 1. Загружаем .env файл (только для локальной разработки)
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, используем системные переменные")
	}

	// 2. Достаем строку подключения из окружения
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal("Переменная DB_URL не задана")
	}

	// 3. Подключаемся к базе
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer db.Close()

	// 4. Проверяем соединение
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Успешно подключились через секреты!")

	// 5. Получаем список таблиц из схемы public
	rows, err := db.Query(`
			SELECT table_name 
			FROM information_schema.tables 
			WHERE table_schema = 'public' 
			ORDER BY table_name
		`)
	if err != nil {
		log.Fatalf("Ошибка запроса списка таблиц: %v", err)
	}
	defer rows.Close()

	// 6. Выводим названия таблиц
	log.Println("Список таблиц в базе pioneer (схема public):")
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		log.Println(" -", tableName)
	}

	// 7. Проверяем на ошибки после итерации
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

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
