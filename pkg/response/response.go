package response

import "github.com/gofiber/fiber/v3"

type Meta struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Response struct {
	Meta Meta `json:"meta"`
	Data any  `json:"data"`
}

type PaginationMeta struct {
	Meta
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type PaginationResponse struct {
	Meta PaginationMeta `json:"meta"`
	Data any            `json:"data"`
}

func Success(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(Response{
		Meta: Meta{
			Success: true,
			Message: message,
		},
		Data: data,
	})
}

func Created(c fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusCreated).JSON(Response{
		Meta: Meta{
			Success: true,
			Message: message,
		},
		Data: data,
	})
}

func Error(c fiber.Ctx, statusCode int, message string) error {
	return c.Status(statusCode).JSON(Response{
		Meta: Meta{
			Success: false,
			Message: message,
		},
		Data: nil,
	})
}

func BadRequest(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusBadRequest, message)
}

func Unauthorized(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusUnauthorized, message)
}

func NotFound(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message)
}

func InternalServerError(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusInternalServerError, message)
}

func Conflict(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusConflict, message)
}