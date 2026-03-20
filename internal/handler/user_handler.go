package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Onboarding(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	req := new(request.OnboardingRequest)
	if err := c.Bind().JSON(req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if req.Gender == "" || req.AvatarID == 0 {
		return response.BadRequest(c, "gender dan avatar_id wajib diisi")
	}

	if req.AvatarID < 1 || req.AvatarID > 4 {
		return response.BadRequest(c, "avatar_id harus antara 1 sampai 4")
	}

	result, err := h.userService.Onboarding(c.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			return response.NotFound(c, "user tidak ditemukan")
		case errors.Is(err, service.ErrAlreadyOnboarded):
			return response.BadRequest(c, "user sudah melakukan onboarding")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}

	return response.Success(c, "onboarding berhasil", result)
}