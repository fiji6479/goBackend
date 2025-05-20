package handler

import (
	"encoding/json"
	"goBackend/internal/model"
	"goBackend/internal/service"
	"net/http"
)

// Handler содержит ссылку на Evaluator
type Handler struct {
	evaluator *service.Evaluator
}

// NewHandler создает новый экземпляр Handler с инициализированным Evaluator
func NewHandler() *Handler {
	return &Handler{
		evaluator: service.NewEvaluator(),
	}
}

// Calculate обрабатывает POST-запрос на /calculate
func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	var instructions []model.Instruction

	// Попытка распарсить JSON из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&instructions); err != nil {
		http.Error(w, "invalid JSON input", http.StatusBadRequest)
		return
	}

	// Выполняем инструкции через Evaluator
	result, err := h.evaluator.EvalInstructions(instructions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.Output{Items: result})
}
