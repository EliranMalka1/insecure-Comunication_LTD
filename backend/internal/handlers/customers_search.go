// internal/handlers/customers_search.go
package handlers

import (
	"net/http"

	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type CustomerDTO struct {
	ID        int64   `db:"id" json:"id"`
	Name      string  `db:"name" json:"name"`
	Email     string  `db:"email" json:"email"`
	Phone     *string `db:"phone" json:"phone,omitempty"`
	Notes     *string `db:"notes" json:"notes,omitempty"`
	CreatedAt string  `db:"created_at" json:"created_at"`
}

func SearchCustomers(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {

		q := c.QueryParam("q")

		if len(q) == 0 {
			return c.JSON(http.StatusOK, map[string]any{
				"items": []CustomerDTO{},
				"total": 0,
			})
		}

		// *** Insecure: SQLi-vulnerable query construction with string interpolation ***
		query := fmt.Sprintf(`
			SELECT id, name, email, phone, notes, created_at
			FROM customers
			WHERE name LIKE '%%%s%%'
			ORDER BY created_at DESC
			LIMIT 50
		`, q)

		var rows []CustomerDTO
		if err := db.Select(&rows, query); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"items": rows,
			"q":     q,
		})
	}
}
