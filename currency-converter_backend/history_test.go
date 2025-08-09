package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHistoryHandler(t *testing.T) {
	history = []Conversion{}

	req := httptest.NewRequest(http.MethodGet, "/history", nil)
	w := httptest.NewRecorder()

	historyHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидается статус 200, но получен %d", resp.StatusCode)
	}

	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("Ожидался Content-Type application/json, но получен %s", ct)
	}

	var result []Conversion
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Errorf("Ошибка при разборе JSON: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Ожидалась пустая история, но получено %d элементов", len(result))
	}
}
