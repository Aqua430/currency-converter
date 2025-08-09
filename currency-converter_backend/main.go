package main

import (
	"currency-converter_backend/cache"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Conversion struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
	Result float64 `json:"result"`
}

var history []Conversion

const maxHistorySize = 100

func saveHistory() {
	file, err := os.Create("history.json")
	if err != nil {
		log.Println("Ошибка при создании файла:", err)
		return
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(history)
	if err != nil {
		log.Println("Ошибка при сохранении истории:", err)
	}
}

func historyLoad() {
	file, err := os.Open("history.json")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&history)
	if err != nil {
		log.Println("Ошибка при разборе JSON:", err)
	}
}

func convertHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	from := req.URL.Query().Get("from")
	to := req.URL.Query().Get("to")
	amountStr := req.URL.Query().Get("amount")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount < 0 {
		http.Error(w, "Неккоректное значение amount", http.StatusBadRequest)
		return
	}
	if from == "" || to == "" || amountStr == "" {
		http.Error(w, "Missing parameters", http.StatusBadRequest)
		return
	}

	key := fmt.Sprintf("rate:%s:%s", from, to)
	cachedRateStr, err := cache.Get(key)
	if err == nil && cachedRateStr != "" {
		rate, _ := strconv.ParseFloat(cachedRateStr, 64)
		result := rate * amount

		json.NewEncoder(w).Encode(map[string]interface{}{
			"from":   from,
			"to":     to,
			"amount": amount,
			"result": result,
			"cached": true,
		})
		return
	}

	url := fmt.Sprintf("https://open.er-api.com/v6/latest/%s", from)

	resp, err := http.Get(url)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Error(w, "Не удалось отправить запрос", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data struct {
		Rates map[string]float64 `json:"rates"`
	}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Ошибка при разборе ответа", http.StatusInternalServerError)
		return
	}

	rate, ok := data.Rates[to]
	if !ok {
		http.Error(w, "Не найдена валюта "+to, http.StatusBadRequest)
		return
	}
	result := rate * amount

	cache.Set(key, fmt.Sprintf("%f", rate), 10*time.Minute)

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
		"result": result,
	})
	conversion := Conversion{
		From:   from,
		To:     to,
		Amount: amount,
		Result: result,
	}
	history = append(history, conversion)
	if len(history) > maxHistorySize {
		history = history[len(history)-maxHistorySize:]
	}
	saveHistory()
}

func historyHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func main() {
	cache.InitRedis()
	historyLoad()
	http.HandleFunc("/history", historyHandler)
	http.HandleFunc("/convert", convertHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Currency Converter API is running!")
	})
	http.ListenAndServe(":8080", nil)
}
