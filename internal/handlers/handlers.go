package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/strelov1/hire/internal/db"
)

// Register wires all routes onto the application.
func Register(app *fiber.App, pool *pgxpool.Pool) {
	h := &Handler{pool: pool, q: db.New(pool)}

	app.Get("/health", h.Health)

	api := app.Group("/api/v1")
	api.Get("/jobs", h.ListJobs)
	api.Get("/jobs/:id", h.GetJob)
}

// Handler holds dependencies shared across HTTP handlers.
type Handler struct {
	pool *pgxpool.Pool
	q    *db.Queries
}
