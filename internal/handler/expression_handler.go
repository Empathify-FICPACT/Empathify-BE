package handler

import (
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/middleware"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type ExpressionHandler struct {
	exprService service.ExpressionService
}

func NewExpressionHandler(exprService service.ExpressionService) *ExpressionHandler {
	return &ExpressionHandler{exprService: exprService}
}

// @Summary      Mulai sesi mengenal ekspresi
// @Tags         Expression
// @Accept       json
// @Produce      json
// @Param        body  body      request.StartExpressionRequest  false  "Jumlah ekspresi (3-5, default 3)"
// @Success      201   {object}  response.Response{data=response.ExpressionSessionResponse}
// @Failure      500   {object}  response.Response
// @Security     BearerAuth
// @Router       /expression/sessions [post]
func (h *ExpressionHandler) StartSession(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)

	req := new(request.StartExpressionRequest)
	if err := c.Bind().JSON(req); err != nil {
		req = &request.StartExpressionRequest{TotalExpressions: 3}
	}

	result, err := h.exprService.StartSession(c.Context(), userID, req)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}
	return response.Created(c, "sesi dimulai", result)
}

// @Summary      Submit foto ekspresi
// @Tags         Expression
// @Accept       multipart/form-data
// @Produce      json
// @Param        id            path      string  true  "Session ID"
// @Param        reference_id  formData  string  true  "Reference ID ekspresi yang ditiru"
// @Param        photo         formData  file    true  "Foto ekspresi user (jpg/png)"
// @Success      200           {object}  response.Response{data=response.ExpressionAttemptResponse}
// @Failure      400           {object}  response.Response
// @Failure      403           {object}  response.Response
// @Failure      404           {object}  response.Response
// @Failure      500           {object}  response.Response
// @Security     BearerAuth
// @Router       /expression/sessions/{id}/attempts [post]
func (h *ExpressionHandler) SubmitAttempt(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	sessionID := c.Params("id")
	referenceID := c.FormValue("reference_id")

	if referenceID == "" {
		return response.BadRequest(c, "reference_id wajib diisi")
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil {
		return response.BadRequest(c, "foto wajib disertakan")
	}

	if fileHeader.Size > 5*1024*1024 {
		return response.BadRequest(c, "ukuran foto maksimal 5MB")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return response.InternalServerError(c, "gagal membaca foto")
	}
	defer file.Close()

	photoBytes := make([]byte, fileHeader.Size)
	if _, err := file.Read(photoBytes); err != nil {
		return response.InternalServerError(c, "gagal membaca foto")
	}

	result, err := h.exprService.SubmitAttempt(c.Context(), userID, sessionID, referenceID, photoBytes)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrExpressionSessionNotFound):
			return response.NotFound(c, "sesi tidak ditemukan")
		case errors.Is(err, service.ErrExpressionSessionNotOwned):
			return response.Error(c, fiber.StatusForbidden, "akses ditolak")
		case errors.Is(err, service.ErrExpressionSessionCompleted):
			return response.BadRequest(c, "sesi sudah selesai")
		case errors.Is(err, service.ErrReferenceNotFound):
			return response.NotFound(c, "referensi ekspresi tidak ditemukan")
		case errors.Is(err, service.ErrAllAttemptsSubmitted):
			return response.BadRequest(c, "semua ekspresi sudah disubmit")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Success(c, "ekspresi berhasil dianalisis", result)
}

// @Summary      Selesaikan sesi ekspresi
// @Tags         Expression
// @Produce      json
// @Param        id   path      string  true  "Session ID"
// @Success      200  {object}  response.Response{data=response.CompleteExpressionResponse}
// @Failure      400  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     BearerAuth
// @Router       /expression/sessions/{id}/complete [patch]
func (h *ExpressionHandler) CompleteSession(c fiber.Ctx) error {
	userID := middleware.GetUserID(c)
	sessionID := c.Params("id")

	result, err := h.exprService.CompleteSession(c.Context(), userID, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrExpressionSessionNotFound):
			return response.NotFound(c, "sesi tidak ditemukan")
		case errors.Is(err, service.ErrExpressionSessionNotOwned):
			return response.Error(c, fiber.StatusForbidden, "akses ditolak")
		case errors.Is(err, service.ErrExpressionSessionCompleted):
			return response.BadRequest(c, "sesi sudah selesai")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}
	return response.Success(c, "sesi selesai", result)
}