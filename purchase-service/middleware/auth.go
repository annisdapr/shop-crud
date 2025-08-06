// File: middleware/auth.go
package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

var (
	// ErrMissingAuthHeader if missing header auth or wrong format
	ErrMissingAuthHeader = echo.NewHTTPError(http.StatusUnauthorized, "Missing or malformed JWT")
	// ErrInvalidJWT if invalid token or expired
	ErrInvalidJWT        = echo.NewHTTPError(http.StatusUnauthorized, "Invalid or expired JWT")
)

// JWTAuthMiddleware create an instance middleware Echo for JWT token validation .
func JWTAuthMiddleware(jwtSecret string) echo.MiddlewareFunc {
	// Middleware di Echo adalah sebuah fungsi yang mengembalikan fungsi lain.
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		// Fungsi inilah yang akan dieksekusi untuk setiap request.
		return func(c echo.Context) error {
			// 1. Ambil nilai dari header "Authorization".
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return ErrMissingAuthHeader
			}

			// 2. Cek format header. Harus "Bearer <token>".
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				return ErrMissingAuthHeader
			}
			tokenString := parts[1]

			// 3. Parse dan validasi token.
			// jwt.Parse akan memeriksa tanda tangan, waktu kedaluwarsa (exp), dan waktu penerbitan (iat).
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				// Pastikan metode signing yang digunakan adalah HMAC (seperti HS256),
				// ini untuk mencegah serangan downgrade algoritma.
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				// Kembalikan secret key untuk verifikasi.
				return []byte(jwtSecret), nil
			})

			// Jika ada error saat parsing atau token tidak valid, kembalikan error.
			if err != nil || !token.Valid {
				return ErrInvalidJWT
			}

			// 4. Ambil 'claims' (data/payload) dari token.
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				return ErrInvalidJWT
			}

			// 5. Simpan claims ke dalam context Echo dengan kunci "user".
			// Ini adalah langkah krusial yang memungkinkan handler selanjutnya
			// untuk mengakses informasi pengguna (seperti user_id, email, dll).
			c.Set("user", claims)

			// Jika semua validasi berhasil, lanjutkan ke handler berikutnya dalam chain.
			return next(c)
		}
	}
}

// GetUserFromContext adalah fungsi helper untuk mengambil klaim JWT dari context Echo.
// Ini membantu agar kode di handler menjadi lebih bersih.
func GetUserFromContext(c echo.Context) (jwt.MapClaims, bool) {
	// Ambil data dari context dengan kunci "user".
	user := c.Get("user")
	if user == nil {
		return nil, false
	}
	// Lakukan type assertion untuk mengubahnya kembali menjadi jwt.MapClaims.
	claims, ok := user.(jwt.MapClaims)
	return claims, ok
}
