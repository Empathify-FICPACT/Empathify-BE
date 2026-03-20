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

// @Summary      Onboarding user
// @Description  Isi gender dan avatar setelah register email atau login Google pertama kali
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body  body      request.OnboardingRequest  true  "Onboarding payload"
// @Success      200   {object}  response.Response{data=response.UserResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Security     BearerAuth
// @Router       /user/onboarding [patch]
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