package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"secure-communication-ltd/backend/internal/services"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(db *sqlx.DB, pol services.PasswordPolicy) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req RegisterRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
		}

		// INSECURE: No Trim/No validation - to allow POC of SQLi
		if req.Username == "" || req.Email == "" || req.Password == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing fields"})
		}

		// *** Secure: Parametric check for user/email existence to avoid failure before the vulnerable INSERT ***
		var exists int
		if err := db.Get(&exists,
			"SELECT COUNT(*) FROM users WHERE username = ? OR email = ?",
			req.Username, req.Email,
		); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
		}
		if exists > 0 {
			return c.JSON(http.StatusConflict, map[string]string{"error": "username or email already exists"})
		}

		// Prepare salt and hash (it looks "legitimate" to the eye, but the INSERT below will be vulnerable)
		salt, err := services.GenerateSalt16()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "salt error"})
		}
		hashHex, err := services.HashPasswordHMACHex(req.Password, salt)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "hash error"})
		}

		// *** Insecure: INSERT with string interpolation -> can be injected through username/email ***
		qInsert := fmt.Sprintf(
			"INSERT INTO users SET "+
				"username='%s', "+
				"email='%s', "+
				"password_hmac='%s', "+
				"salt=UNHEX('%x'), "+
				"is_active=1, "+
				"is_verified=0",
			req.Username, req.Email, hashHex, salt,
		)
		c.Logger().Printf("qInsert => %s\n", qInsert)

		res, err := db.Exec(qInsert)
		if err != nil {
			// DEBUG ONLY:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "insert error: " + err.Error()})
		}

		uid, _ := res.LastInsertId()

		// Create a verification token and send it by email (remains as in the secure version)
		vTok, err := services.NewVerificationToken(24 * time.Hour)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token error"})
		}
		if _, err := db.Exec(
			"INSERT INTO email_verification_tokens (user_id, token_sha1, expires_at) VALUES (?, ?, ?)",
			uid, vTok.SHA1Hex, vTok.ExpiresAt,
		); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "token save error"})
		}

		mailer, err := services.NewMailerFromEnv()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "mailer error"})
		}
		base := os.Getenv("BACKEND_PUBLIC_URL")
		if base == "" {
			base = "http://localhost:8081"
		}
		link := fmt.Sprintf("%s/api/verify-email?token=%s", strings.TrimRight(base, "/"), vTok.Raw)

		html := fmt.Sprintf(`
			<h2>Verify your email</h2>
			<p>Hi %s, thanks for registering.</p>
			<p>Please click the button below to verify your email address:</p>
			<p><a href="%s" style="display:inline-block;padding:10px 16px;border-radius:8px;background:#4f9cff;color:#fff;text-decoration:none">Verify Email</a></p>
			<p>If the button doesn't work, copy this URL:</p>
			<p><code>%s</code></p>
		`, req.Username, link, link)

		if err := mailer.Send(req.Email, "Verify your email", html); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "send mail error"})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "User registered. Please check your email to verify your account.",
		})
	}
}

// Kept in case you want to use it sometime
func looksLikeEmail(s string) bool {
	return strings.Count(s, "@") == 1 && len(s) >= 6 && strings.Contains(s, ".")
}
