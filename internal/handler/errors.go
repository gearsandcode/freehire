package handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrorHandler is the single place every error returned by a handler becomes an
// HTTP response. It is wired into fiber.New so the error envelope (`{"error":
// ...}`, mirroring the `{"data": ...}` success shape) and the status mapping
// live in one place instead of being hand-rolled per handler:
//
//   - a *fiber.Error (from fiber.NewError) keeps its code and message — this is
//     how handlers declare a specific HTTP meaning (e.g. 400 "invalid job id");
//   - a not-found from the DB layer (pgx.ErrNoRows) maps to 404, so read
//     handlers can just `return err`;
//   - a foreign-key violation (a write referencing a missing parent row, e.g.
//     applying to a non-existent job id) also maps to 404 — the referenced
//     resource doesn't exist;
//   - anything else is an unexpected failure: 500 with a generic message, never
//     leaking internals.
func ErrorHandler(c *fiber.Ctx, err error) error {
	status := fiber.StatusInternalServerError
	msg := "internal server error"

	var fe *fiber.Error
	switch {
	case errors.As(err, &fe):
		status, msg = fe.Code, fe.Message
	case errors.Is(err, pgx.ErrNoRows), isForeignKeyViolation(err):
		status, msg = fiber.StatusNotFound, "not found"
	}

	return c.Status(status).JSON(fiber.Map{"error": msg})
}

// isForeignKeyViolation reports whether err is a Postgres foreign-key violation
// (SQLSTATE 23503) — a write that references a row that doesn't exist.
func isForeignKeyViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23503"
}
