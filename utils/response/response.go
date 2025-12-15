package response

import (
	"encoding/json"
	"net/http"
)

// Response adalah struktur dasar untuk semua respons API (sukses atau error).
type Response struct {
	Status string `json:"status"`
	// Message string      `json:"message"`
	Data interface{} `json:"data,omitempty"` // omitempty: sembunyikan field jika nil/kosong
}

// JSON menulis respons JSON standar ke http.ResponseWriter.
func JSON(w http.ResponseWriter, statusCode int, status string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	resp := Response{
		Status: status,
		Data:   data,
	}

	// Langsung encode struktur Response menjadi JSON dan tulis ke response writer
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		// Log error jika gagal encode (jarang terjadi)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
