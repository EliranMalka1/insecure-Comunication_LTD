// internal/handlers/login.go
package handlers

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"secure-communication-ltd/backend/internal/services"
)

type LoginRequest struct {
	ID       string `json:"id"`       // email or username (no Trim to not break PoC)
	Password string `json:"password"` // sent as-is
}

type userRow struct {
	ID         int64  `db:"id"`
	Username   string `db:"username"`
	Email      string `db:"email"`
	PassHMAC   string `db:"password_hmac"`
	Salt       []byte `db:"salt"`
	IsActive   bool   `db:"is_active"`
	IsVerified bool   `db:"is_verified"`
}

func Login(db *sqlx.DB, _ services.PasswordPolicy) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
		}
		if req.ID == "" || len(req.Password) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing fields"})
		}

		// ---- INSECURE #1: fetch salt via string concatenation (SQLi possible)
		var salt []byte
		qSalt := fmt.Sprintf(
			"SELECT salt FROM users WHERE email = '%s' OR username = '%s' LIMIT 1",
			req.ID, req.ID,
		)
		if err := db.Get(&salt, qSalt); err != nil {
			// generic message (matches your UI expectation)
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}

		// legit-looking HMAC calc (still allows SQLi bypass later)
		hexHMAC, err := services.HashPasswordHMACHex(req.Password, salt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "hash error"})
		}

		// ---- INSECURE #2: logic/precedence bug + string concat (SQLi possible)

		var u userRow
		qUser := fmt.Sprintf(
			"SELECT id, username, email, password_hmac, salt, is_active, is_verified "+
				"FROM users "+
				"WHERE email = '%s' OR username = '%s' AND password_hmac = '%s' "+
				"LIMIT 1",
			req.ID, req.ID, hexHMAC,
		)
		if err := db.Get(&u, qUser); err != nil || !u.IsActive || !u.IsVerified {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}

		// ===== NO MFA (vulnerable): issue JWT cookie immediately =====
		token, err := services.CreateJWT(u.ID, u.Username, 24*time.Hour)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token error"})
		}

		cookie := &http.Cookie{
			Name:     services.CookieName,
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Secure:   false,
		}
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{"message": "ok"})
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
