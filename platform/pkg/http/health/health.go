package healthhttp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/Alexander-Mandzhiev/school_schedule/platform/pkg/logger"
)

type Response struct {
	Status string `json:"status"`
}

// Handler возвращает простой JSON-ответ статуса для healthcheck.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(Response{Status: "ok"}); err != nil {
		logger.Error(r.Context(),
			"❌ [Health] Не удалось отправить ответ healthcheck",
			zap.Error(err),
		)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

// LiveHandler — эндпоинт liveness (приложение запущено).
func LiveHandler(w http.ResponseWriter, r *http.Request) {
	Handler(w, r)
}

// ReadyHandler — эндпоинт readiness (приложение готово обслуживать трафик).
func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	Handler(w, r)
}

// StartHandler — эндпоинт startup (инициализация завершена).
func StartHandler(w http.ResponseWriter, r *http.Request) {
	Handler(w, r)
}
