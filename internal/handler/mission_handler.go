package handler

import (
	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type MissionHandler struct {
	missionService service.MissionService
}

func NewMissionHandler(missionService service.MissionService) *MissionHandler {
	return &MissionHandler{missionService: missionService}
}

// @Summary      Ambil misi harian user
// @Description  Menampilkan 3 misi harian beserta progress user hari ini
// @Tags         Mission
// @Produce      json
// @Success      200  {object}  response.Response{data=response.DailyMissionsResponse}
// @Failure      500  {object}  response.Response
// @Security     BearerAuth
// @Router       /missions/today [get]
func (h *MissionHandler) GetTodayMissions(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	result, err := h.missionService.GetTodayMissions(c.Context(), userID)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}
	return response.Success(c, "berhasil mengambil misi harian", result)
}