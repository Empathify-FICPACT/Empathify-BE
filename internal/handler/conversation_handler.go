package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type ConversationHandler struct {
	convService service.ConversationService
}

func NewConversationHandler(convService service.ConversationService) *ConversationHandler {
	return &ConversationHandler{convService: convService}
}

// @Summary      Ambil semua topik percakapan
// @Tags         Conversation
// @Produce      json
// @Success      200  {object}  response.Response{data=[]response.TopicResponse}
// @Failure      500  {object}  response.Response
// @Security     BearerAuth
// @Router       /conversation/topics [get]
func (h *ConversationHandler) GetTopics(c fiber.Ctx) error {
	topics, err := h.convService.GetTopics(c.Context())
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}
	return response.Success(c, "berhasil mengambil topik", topics)
}

// @Summary      Mulai sesi percakapan
// @Tags         Conversation
// @Accept       json
// @Produce      json
// @Param        body  body      request.StartConversationRequest  true  "Topic ID"
// @Success      201   {object}  response.Response{data=response.SessionResponse}
// @Failure      400   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Security     BearerAuth
// @Router       /conversation/sessions [post]
func (h *ConversationHandler) StartSession(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	req := new(request.StartConversationRequest)
	if err := c.Bind().JSON(req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}
	if req.TopicID == "" {
		return response.BadRequest(c, "topic_id wajib diisi")
	}

	result, err := h.convService.StartSession(c.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTopicNotFound):
			return response.NotFound(c, "topik tidak ditemukan")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Created(c, "sesi dimulai", result)
}

// @Summary      Kirim pesan suara dalam sesi percakapan
// @Tags         Conversation
// @Accept       multipart/form-data
// @Produce      json
// @Param        id     path      string  true   "Session ID"
// @Param        audio  formData  file    true   "File audio (webm/opus)"
// @Success      200    {object}  response.Response{data=response.SendMessageResponse}
// @Failure      400    {object}  response.Response
// @Failure      403    {object}  response.Response
// @Failure      404    {object}  response.Response
// @Failure      500    {object}  response.Response
// @Security     BearerAuth
// @Router       /conversation/sessions/{id}/messages [post]
func (h *ConversationHandler) SendMessage(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	sessionID := c.Params("id")

	// ambil file audio dari form
	fileHeader, err := c.FormFile("audio")
	if err != nil {
		return response.BadRequest(c, "file audio wajib disertakan")
	}

	// validasi ukuran file max 10MB
	if fileHeader.Size > 10*1024*1024 {
		return response.BadRequest(c, "ukuran file maksimal 10MB")
	}

	// baca isi file
	file, err := fileHeader.Open()
	if err != nil {
		return response.InternalServerError(c, "gagal membaca file audio")
	}
	defer file.Close()

	audioBytes := make([]byte, fileHeader.Size)
	if _, err := file.Read(audioBytes); err != nil {
		return response.InternalServerError(c, "gagal membaca file audio")
	}

	result, err := h.convService.SendMessage(c.Context(), userID, sessionID, audioBytes)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSessionNotFound):
			return response.NotFound(c, "sesi tidak ditemukan")
		case errors.Is(err, service.ErrSessionNotOwned):
			return response.Error(c, fiber.StatusForbidden, "akses ditolak")
		case errors.Is(err, service.ErrSessionCompleted):
			return response.BadRequest(c, "sesi sudah selesai")
		case errors.Is(err, service.ErrSTTFailed):
			return response.BadRequest(c, "gagal memproses audio, coba rekam ulang")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Success(c, "pesan terkirim", result)
}

// @Summary      Selesaikan sesi percakapan
// @Tags         Conversation
// @Produce      json
// @Param        id   path      string  true  "Session ID"
// @Success      200  {object}  response.Response{data=response.CompleteSessionResponse}
// @Failure      400  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     BearerAuth
// @Router       /conversation/sessions/{id}/complete [patch]
func (h *ConversationHandler) CompleteSession(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	sessionID := c.Params("id")

	result, err := h.convService.CompleteSession(c.Context(), userID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSessionNotFound):
			return response.NotFound(c, "sesi tidak ditemukan")
		case errors.Is(err, service.ErrSessionNotOwned):
			return response.Error(c, fiber.StatusForbidden, "akses ditolak")
		case errors.Is(err, service.ErrSessionCompleted):
			return response.BadRequest(c, "sesi sudah selesai")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Success(c, "sesi selesai", result)
}