package handlers

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"secure-communication-ltd/backend/internal/services"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type forgotReq struct {
	Email string `json:"email"`
}

// Always return a generic response — even if the email does not exist.
func PasswordForgot(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req forgotReq
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
		}
		email := strings.TrimSpace(req.Email)
		if email == "" || !strings.Contains(email, "@") {

			return c.JSON(http.StatusOK, map[string]string{"message": "If this email exists, a reset link has been sent."})
		}

		var userID int64
		err := db.Get(&userID, `SELECT id FROM users WHERE email = ? AND is_verified = TRUE LIMIT 1`, email)
		if err != nil {

			return c.JSON(http.StatusOK, map[string]string{"message": "If this email exists, a reset link has been sent."})
		}

		_, _ = db.Exec(`UPDATE password_reset_tokens SET used_at = NOW() WHERE user_id = ? AND used_at IS NULL`, userID)

		raw, err := services.NewRandomBase64URL(32)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token error"})
		}
		h := sha1.Sum([]byte(raw))
		sha := hex.EncodeToString(h[:])

		exp := time.Now().Add(30 * time.Minute)
		if _, err := db.Exec(`
			INSERT INTO password_reset_tokens (user_id, token_sha1, expires_at)
			VALUES (?, ?, ?)`, userID, sha, exp); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
		}

		// Send email
		mailer, err := services.NewMailerFromEnv()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "mailer error"})
		}
		base := os.Getenv("BACKEND_PUBLIC_URL")
		if base == "" {
			base = "http://localhost:8081"
		}

		link := fmt.Sprintf("%s/api/password/reset?token=%s", strings.TrimRight(base, "/"), raw)

		html := fmt.Sprintf(`
			<h2>Password Reset</h2>
			<p>We received a request to reset your password.</p>
			<p><a href="%s" style="display:inline-block;padding:10px 16px;border-radius:8px;background:#4f9cff;color:#fff;text-decoration:none">Set a new password</a></p>
			<p>If the button doesn't work, copy this URL:</p>
			<p><code>%s</code></p>
			<p>This link expires in 30 minutes.</p>
		`, link, link)

		_ = mailer.Send(email, "Reset your password", html)

		return c.JSON(http.StatusOK, map[string]string{"message": "If this email exists, a reset link has been sent."})
	}
}
