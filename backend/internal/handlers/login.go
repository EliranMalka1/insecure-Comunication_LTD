// internal/handlers/login.go
package handlers

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"secure-communication-ltd/backend/internal/services"
)

type LoginRequest struct {
	ID       string `json:"id"` // email or username (לא נעשה Trim כדי לא לפגוע ב-POC)
	Password string `json:"password"`
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

func Login(db *sqlx.DB, pol services.PasswordPolicy) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
		}

		if req.ID == "" || len(req.Password) == 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing fields"})
		}

		var salt []byte
		qSalt := fmt.Sprintf(
			"SELECT salt FROM users WHERE email = '%s' OR username = '%s' LIMIT 1",
			req.ID, req.ID,
		)
		if err := db.Get(&salt, qSalt); err != nil {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials1"})
		}

		hexHMAC, err := services.HashPasswordHMACHex(req.Password, salt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "hash error"})
		}

		var u userRow
		qUser := fmt.Sprintf(
			"SELECT id, username, email, password_hmac, salt, is_active, is_verified "+
				"FROM users "+
				"WHERE email = '%s' OR username = '%s' AND password_hmac = '%s' "+
				"LIMIT 1",
			req.ID, req.ID, hexHMAC,
		)
		err = db.Get(&u, qUser)
		if err != nil || !u.IsActive || !u.IsVerified {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
		}

		mailer, err := services.NewMailerFromEnv()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "mailer error"})
		}

		ttl := 10
		maxAtt := 5
		if v := os.Getenv("MFA_OTP_TTL_MINUTES"); v != "" {
			if n, e := strconv.Atoi(v); e == nil && n > 0 && n <= 60 {
				ttl = n
			}
		}
		if v := os.Getenv("MFA_OTP_MAX_ATTEMPTS"); v != "" {
			if n, e := strconv.Atoi(v); e == nil && n >= 1 && n <= 10 {
				maxAtt = n
			}
		}
		cfg := services.OTPConfig{TTLMinutes: ttl, MaxAttempts: maxAtt}
		if err := services.StartEmailOTP(db, mailer, u.ID, u.Email, cfg); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "otp start error"})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"mfa_required": true,
			"method":       "email_otp",
			"expires_in":   ttl,
		})
	}
}

func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
