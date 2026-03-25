package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type BadgeHandler struct {
	badgeService service.BadgeService
}

func NewBadgeHandler(badgeService service.BadgeService) *BadgeHandler {
	return &BadgeHandler{badgeService: badgeService}
}

// @Summary      Ambil semua badge user
// @Description  Menampilkan semua badge beserta status unlock
// @Tags         Badge
// @Produce      json
// @Success      200  {object}  response.Response{data=response.BadgeListResponse}
// @Failure      500  {object}  response.Response
// @Security     BearerAuth
// @Router       /badges [get]
func (h *BadgeHandler) GetUserBadges(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.badgeService.GetUserBadges(c.Context(), userID)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}
	return response.Success(c, "berhasil mengambil badge", result)
}