package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type StoryHandler struct {
	storyService service.StoryService
}

func NewStoryHandler(storyService service.StoryService) *StoryHandler {
	return &StoryHandler{storyService: storyService}
}

// @Summary      Mulai sesi cerita interaktif
// @Tags         Story
// @Accept       json
// @Produce      json
// @Param        body  body      request.StartStoryRequest  false  "Jumlah skenario (5-8, default 5)"
// @Success      201   {object}  response.Response{data=response.StorySessionResponse}
// @Failure      500   {object}  response.Response
// @Security     BearerAuth
// @Router       /story/sessions [post]
func (h *StoryHandler) StartSession(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	req := new(request.StartStoryRequest)
	if err := c.Bind().JSON(req); err != nil {
		req = &request.StartStoryRequest{TotalScenarios: 5}
	}

	result, err := h.storyService.StartSession(c.Context(), userID, req)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}
	return response.Created(c, "sesi dimulai", result)
}

// @Summary      Submit jawaban skenario cerita
// @Tags         Story
// @Accept       json
// @Produce      json
// @Param        id    path      string                          true  "Session ID"
// @Param        body  body      request.SubmitStoryAnswerRequest  true  "Jawaban user"
// @Success      200   {object}  response.Response{data=response.StoryAnswerResponse}
// @Failure      400   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Security     BearerAuth
// @Router       /story/sessions/{id}/answers [post]
func (h *StoryHandler) SubmitAnswer(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	sessionID := c.Params("id")

	req := new(request.SubmitStoryAnswerRequest)
	if err := c.Bind().JSON(req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.ScenarioID == "" || req.ChosenAnswer == "" {
		return response.BadRequest(c, "scenario_id dan chosen_answer wajib diisi")
	}

	result, err := h.storyService.SubmitAnswer(c.Context(), userID, sessionID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStorySessionNotFound):
			return response.NotFound(c, "sesi tidak ditemukan")
		case errors.Is(err, service.ErrStorySessionNotOwned):
			return response.Error(c, fiber.StatusForbidden, "akses ditolak")
		case errors.Is(err, service.ErrStorySessionCompleted):
			return response.BadRequest(c, "sesi sudah selesai")
		case errors.Is(err, service.ErrStoryScenarioNotFound):
			return response.NotFound(c, "skenario tidak ditemukan")
		case errors.Is(err, service.ErrStoryAlreadyAnswered):
			return response.BadRequest(c, "semua skenario sudah dijawab")
		case errors.Is(err, service.ErrInvalidStoryAnswer):
			return response.BadRequest(c, "jawaban harus a, b, c, atau d")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Success(c, "jawaban berhasil disimpan", result)
}

// @Summary      Selesaikan sesi cerita interaktif
// @Tags         Story
// @Produce      json
// @Param        id   path      string  true  "Session ID"
// @Success      200  {object}  response.Response{data=response.CompleteStoryResponse}
// @Failure      400  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     BearerAuth
// @Router       /story/sessions/{id}/complete [patch]
func (h *StoryHandler) CompleteSession(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	sessionID := c.Params("id")

	result, err := h.storyService.CompleteSession(c.Context(), userID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrStorySessionNotFound):
			return response.NotFound(c, "sesi tidak ditemukan")
		case errors.Is(err, service.ErrStorySessionNotOwned):
			return response.Error(c, fiber.StatusForbidden, "akses ditolak")
		case errors.Is(err, service.ErrStorySessionCompleted):
			return response.BadRequest(c, "sesi sudah selesai")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Success(c, "sesi selesai", result)
}