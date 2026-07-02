package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type Event struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Code      string    `json:"code"`
	OwnerID   int64     `json:"owner_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	// Исправляем Баг 4: Читаем JSON вместо Формы
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	if input.Title == "" {
		http.Error(w, "Название не может быть пустым", http.StatusBadRequest)
		return
	}

	ownerID := h.GetCurrentUserID(r)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
        INSERT INTO events (title, owner_id) 
        VALUES ($1, $2) 
        RETURNING code
    `
	
	// Исправляем Баг 1: Объявляем переменную eventCode перед сканированием
	var eventCode string
	err := h.DB.QueryRow(ctx, query, input.Title, ownerID).Scan(&eventCode)
	if err != nil {
		h.Logger.Printf("Ошибка вставки ивента в БД: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Исправляем Баг 5: Отдаем JSON с кодом вместо редиректа
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"code": eventCode, "status": "created"})
}

func (h *Handler) ListEvents(w http.ResponseWriter, r *http.Request) {
	currentOwnerID := h.GetCurrentUserID(r)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	query := `
        SELECT id, title, code, owner_id, is_active, created_at 
        FROM events 
        WHERE owner_id = $1 AND is_active = TRUE
        ORDER BY created_at DESC 
    `
	
	// Исправляем Баг 2: Передаем реальный query и аргумент currentOwnerID
	rows, err := h.DB.Query(ctx, query, currentOwnerID)
	if err != nil {
		h.Logger.Printf("Ошибка БД: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Внутренняя ошибка сервера"})
		return
	}
	defer rows.Close()
	
	var events []Event
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.ID, &e.Title, &e.Code, &e.OwnerID, &e.IsActive, &e.CreatedAt)
		if err != nil {
			h.Logger.Printf("Ошибка сканирования строки: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		events = append(events, e)
	}

	if err = rows.Err(); err != nil {
		h.Logger.Printf("Ошибка после итерации по rows: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if events == nil {
		events = []Event{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(events)
}

// Исправляем Баг 3: Меняем возвращаемый тип с int на int64, чтобы подходил к BIGINT
func (h *Handler) GetCurrentUserID(r *http.Request) int64 {
	return 1 
}