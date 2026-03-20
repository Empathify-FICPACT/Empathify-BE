package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/provider"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/response"
)

type AuthHandler struct {
	authService  service.AuthService
	googleOAuth  *provider.GoogleOAuthProvider
}

func NewAuthHandler(authService service.AuthService, googleOAuth *provider.GoogleOAuthProvider) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		googleOAuth: googleOAuth,
	}
}

func (h *AuthHandler) RegisterEmail(c fiber.Ctx) error {
	req := new(request.RegisterRequest)

	if err := c.Bind().JSON(req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if req.Email == "" || req.Password == "" || req.Gender == "" || req.AvatarID == 0 {
		return response.BadRequest(c, "email, password, gender, dan avatar_id wajib diisi")
	}

	if req.AvatarID < 1 || req.AvatarID > 4 {
		return response.BadRequest(c, "avatar_id harus antara 1 sampai 4")
	}

	result, err := h.authService.RegisterEmail(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			return response.Conflict(c, "email sudah digunakan")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}

	return response.Created(c, "registrasi berhasil", result)
}

func (h *AuthHandler) LoginEmail(c fiber.Ctx) error {
	req := new(request.LoginRequest)

	if err := c.Bind().JSON(req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return response.BadRequest(c, "email dan password wajib diisi")
	}

	result, err := h.authService.LoginEmail(c.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return response.Unauthorized(c, "email atau password salah")
		default:
			return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
		}
	}

	return response.Success(c, "login berhasil", result)
}

// step 1 — redirect user ke halaman login Google
func (h *AuthHandler) GoogleLogin(c fiber.Ctx) error {
	// generate random state untuk CSRF protection
	state, err := generateState()
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}

	// simpan state di cookie sementara
	c.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   60 * 10, // 10 menit
		HTTPOnly: true,
		Secure:   false, // ganti true saat production (HTTPS)
	})

	authURL := h.googleOAuth.GetAuthURL(state)
	return c.Redirect().To(authURL)
}

// step 2 — Google redirect balik ke sini dengan code
func (h *AuthHandler) GoogleCallback(c fiber.Ctx) error {
	// validasi state (CSRF check)
	cookieState := c.Cookies("oauth_state")
	queryState := c.Query("state")

	if cookieState == "" || cookieState != queryState {
		return response.Unauthorized(c, "invalid oauth state")
	}

	// hapus cookie state
	c.Cookie(&fiber.Cookie{
		Name:   "oauth_state",
		Value:  "",
		MaxAge: -1,
	})

	code := c.Query("code")
	if code == "" {
		return response.BadRequest(c, "authorization code tidak ditemukan")
	}

	// tukar code dengan token
	token, err := h.googleOAuth.ExchangeCode(c.Context(), code)
	if err != nil {
		return response.Unauthorized(c, "gagal verifikasi akun Google")
	}

	// ambil info user dari Google
	googleUser, err := h.googleOAuth.GetUserInfo(c.Context(), token)
	if err != nil {
		return response.InternalServerError(c, "gagal mengambil data akun Google")
	}

	// login atau register otomatis
	result, err := h.authService.LoginGoogle(c.Context(), googleUser)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}

	return response.Success(c, "login berhasil", result)
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}