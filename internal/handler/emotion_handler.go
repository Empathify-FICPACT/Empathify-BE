package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type EmotionHandler struct {
	emotionService service.EmotionService
}

func NewEmotionHandler(emotionService service.EmotionService) *EmotionHandler {
	return &EmotionHandler{emotionService: emotionService}
}

// @Summary      Mulai sesi memahami emosi
// @Tags         Emotion
// @Accept       json
// @Produce      json
// @Param        body  body      request.StartEmotionRequest  false  "Jumlah soal (5-8, default 5)"
// @Success      201   {object}  response.Response{data=response.EmotionSessionResponse}
// @Failure      500   {object}  response.Response
// @Security     BearerAuth
// @Router       /emotion/sessions [post]
func (h *EmotionHandler) StartSession(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	req := new(request.StartEmotionRequest)
	if err := c.Bind().JSON(req); err != nil {
		req = &request.StartEmotionRequest{TotalQuestions: 5}
	}

	result, err := h.emotionService.StartSession(c.Context(), userID, req)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}
	return response.Created(c, "sesi dimulai", result)
}

// @Summary      Submit jawaban soal emosi
// @Tags         Emotion
// @Accept       json
// @Produce      json
// @Param        id    path      string                           true  "Session ID"
// @Param        body  body      request.SubmitEmotionAnswerRequest  true  "Jawaban user"
// @Success      200   {object}  response.Response{data=response.EmotionAnswerResponse}
// @Failure      400   {object}  response.Response
// @Failure      403   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Security     BearerAuth
// @Router       /emotion/sessions/{id}/answers [post]
func (h *EmotionHandler) SubmitAnswer(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	sessionID := c.Params("id")

	req := new(request.SubmitEmotionAnswerRequest)
	if err := c.Bind().JSON(req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.QuestionID == "" || req.ChosenAnswer == "" {
		return response.BadRequest(c, "question_id dan chosen_answer wajib diisi")
	}

	result, err := h.emotionService.SubmitAnswer(c.Context(), userID, sessionID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmotionSessionNotFound):
			return response.NotFound(c, "sesi tidak ditemukan")
		case errors.Is(err, service.ErrEmotionSessionNotOwned):
			return response.Error(c, fiber.StatusForbidden, "akses ditolak")
		case errors.Is(err, service.ErrEmotionSessionCompleted):
			return response.BadRequest(c, "sesi sudah selesai")
		case errors.Is(err, service.ErrEmotionQuestionNotFound):
			return response.NotFound(c, "soal tidak ditemukan")
		case errors.Is(err, service.ErrEmotionAlreadyAnswered):
			return response.BadRequest(c, "semua soal sudah dijawab")
		case errors.Is(err, service.ErrInvalidAnswer):
			return response.BadRequest(c, "jawaban harus a, b, c, atau d")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Success(c, "jawaban berhasil disimpan", result)
}

// @Summary      Selesaikan sesi memahami emosi
// @Tags         Emotion
// @Produce      json
// @Param        id   path      string  true  "Session ID"
// @Success      200  {object}  response.Response{data=response.CompleteEmotionResponse}
// @Failure      400  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     BearerAuth
// @Router       /emotion/sessions/{id}/complete [patch]
func (h *EmotionHandler) CompleteSession(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	sessionID := c.Params("id")

	result, err := h.emotionService.CompleteSession(c.Context(), userID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmotionSessionNotFound):
			return response.NotFound(c, "sesi tidak ditemukan")
		case errors.Is(err, service.ErrEmotionSessionNotOwned):
			return response.Error(c, fiber.StatusForbidden, "akses ditolak")
		case errors.Is(err, service.ErrEmotionSessionCompleted):
			return response.BadRequest(c, "sesi sudah selesai")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Success(c, "sesi selesai", result)
}