package middleware

import "net/http"

// CORSMiddleware menambahkan header CORS yang diperlukan ke setiap respons HTTP.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Izinkan semua origin. Sesuaikan '*' dengan domain spesifik Anda
		// (misalnya, "http://localhost:3000") untuk keamanan yang lebih baik.
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Izinkan metode HTTP yang umum digunakan.
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		// Izinkan header kustom yang mungkin dikirim oleh klien (misalnya, untuk otentikasi).
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Izinkan kredensial (seperti cookies atau header otentikasi) disertakan dalam permintaan.
		// Jika disetel ke true, Access-Control-Allow-Origin tidak bisa '*' (harus domain spesifik).
		// w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Preflight request (permintaan OPTIONS)
		// Browser akan mengirim OPTIONS terlebih dahulu untuk memeriksa izin.
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Lanjutkan ke handler yang sebenarnya
		next.ServeHTTP(w, r)
	})
}
