package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type DashboardHandler struct {
	dashboardService service.DashboardService
}

func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// @Summary      Ambil data dashboard user
// @Description  Menampilkan ringkasan XP, statistik fitur, misi harian, dan badge
// @Tags         Dashboard
// @Produce      json
// @Success      200  {object}  response.Response{data=response.DashboardResponse}
// @Failure      500  {object}  response.Response
// @Security     BearerAuth
// @Router       /dashboard [get]
func (h *DashboardHandler) GetDashboard(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.dashboardService.GetDashboard(c.Context(), userID)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}
	return response.Success(c, "berhasil mengambil data dashboard", result)
}