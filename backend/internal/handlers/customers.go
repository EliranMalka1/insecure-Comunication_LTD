// internal/handlers/customers_create.go (vulnerable)
package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type CreateCustomerRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone,omitempty"`
	Notes string `json:"notes,omitempty"`
}

func CreateCustomer(db *sqlx.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req CreateCustomerRequest
		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
		}
		if req.Name == "" || req.Email == "" {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "name and email are required"})
		}

		q := fmt.Sprintf(
			"INSERT INTO customers SET "+
				"name='%s', "+
				"email='%s', "+
				"phone='%s', "+
				"notes='%s'",
			req.Name, req.Email, req.Phone, req.Notes,
		)

		log.Printf("INSERT query:\n%s\n", q)

		if _, err := db.Exec(q); err != nil {
			log.Printf("INSERT error: %v", err)

			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "insert error: " + err.Error(),
				"q":     q,
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"message": "customer created",
		})
	}
}
