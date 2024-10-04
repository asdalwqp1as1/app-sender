package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	url := "https://sdalas.onrender.com"

	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	sendGetRequest(url)

	for range ticker.C {
		sendGetRequest(url)
	}
}

func sendGetRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка при отправке GET-запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Ошибка при чтении ответа: %v\n", err)
		return
	}

	fmt.Printf("Статус: %s\n", resp.Status)
	if len(body) > 100 {
		fmt.Printf("Ответ (первые 100 байт): %s...\n", body[:100])
	} else {
		fmt.Printf("Ответ: %s\n", body)
	}
}
