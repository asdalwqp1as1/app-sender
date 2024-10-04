package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Структура для хранения информации о последнем запросе
type RequestInfo struct {
	Status    string
	Response  string
	Timestamp time.Time
}

var lastRequest RequestInfo

func main() {
	url := "https://sdalas.onrender.com"

	// Запуск HTTP-сервера в отдельной горутине
	go startHTTPServer(":8080")

	// Создание тикера для периодических запросов каждые 5 минут
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	// Первоначальный запрос при запуске
	sendGetRequest(url)

	// Обработка периодических запросов
	for range ticker.C {
		sendGetRequest(url)
	}
}

// Функция для отправки GET-запроса
func sendGetRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка при отправке GET-запроса: %v\n", err)
		updateLastRequest("Ошибка при отправке запроса", "")
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v\n", err)
		updateLastRequest("Ошибка при чтении ответа", "")
		return
	}

	log.Printf("Статус: %s\n", resp.Status)
	if len(body) > 100 {
		log.Printf("Ответ (первые 100 байт): %s...\n", body[:100])
		updateLastRequest(resp.Status, string(body[:100])+"...")
	} else {
		log.Printf("Ответ: %s\n", body)
		updateLastRequest(resp.Status, string(body))
	}
}

// Функция для обновления информации о последнем запросе
func updateLastRequest(status, response string) {
	lastRequest = RequestInfo{
		Status:    status,
		Response:  response,
		Timestamp: time.Now(),
	}
}

// Функция для запуска HTTP-сервера
func startHTTPServer(addr string) {
	http.HandleFunc("/", handler)
	log.Printf("Запуск HTTP-сервера на %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Ошибка запуска HTTP-сервера: %v\n", err)
	}
}

// Обработчик HTTP-запросов
func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := fmt.Sprintf(`{
	"last_request": {
		"status": "%s",
		"response": "%s",
		"timestamp": "%s"
	}
}`, lastRequest.Status, lastRequest.Response, lastRequest.Timestamp.Format(time.RFC3339))
	w.Write([]byte(response))
}
