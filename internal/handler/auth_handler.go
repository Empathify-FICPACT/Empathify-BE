package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"

	"github.com/gofiber/fiber/v3"

	"github.com/Empathify-FICPACT/Empathify-BE/internal/provider"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/service"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/dto/request"
	"github.com/Empathify-FICPACT/Empathify-BE/pkg/jwt"
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

// @Summary      Register dengan email
// @Description  Daftarkan akun baru menggunakan email dan password. Gender dan avatar diisi saat onboarding.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      request.RegisterRequest  true  "Register payload"
// @Success      201   {object}  response.Response{data=response.AuthResponse}
// @Failure      400   {object}  response.Response
// @Failure      409   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterEmail(c fiber.Ctx) error {
	req := new(request.RegisterRequest)

	if err := c.Bind().JSON(req); err != nil {
		return response.BadRequest(c, "invalid request body")
	}

	if req.Email == "" || req.Password == "" {
		return response.BadRequest(c, "email dan password wajib diisi")
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

// @Summary      Login dengan email
// @Description  Login menggunakan email dan password yang sudah terdaftar
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      request.LoginRequest  true  "Login payload"
// @Success      200   {object}  response.Response{data=response.AuthResponse}
// @Failure      400   {object}  response.Response
// @Failure      401   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /auth/login [post]
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

// @Summary      Login dengan Google
// @Description  Redirect user ke halaman login Google OAuth
// @Tags         Auth
// @Success      302
// @Failure      500  {object}  response.Response
// @Router       /auth/google [get]
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
	// return c.Redirect().Status(fiber.StatusTemporaryRedirect).To(authURL)
}

// @Summary      Google OAuth callback
// @Description  Callback dari Google setelah user berhasil login. Otomatis register jika akun belum ada.
// @Tags         Auth
// @Produce      json
// @Param        code   query     string  true  "Authorization code dari Google"
// @Param        state  query     string  true  "State untuk CSRF protection"
// @Success      200    {object}  response.Response{data=response.AuthResponse}
// @Failure      400    {object}  response.Response
// @Failure      401    {object}  response.Response
// @Failure      500    {object}  response.Response
// @Router       /auth/google/callback [get]
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
		log.Println("ERROR ExchangeCode:", err)
		return response.Unauthorized(c, "gagal verifikasi akun Google")
	}
	log.Println("Token OK:", token.AccessToken[:10])

	// ambil info user dari Google
	googleUser, err := h.googleOAuth.GetUserInfo(c.Context(), token)
	if err != nil {
		log.Println("ERROR GetUserInfo:", err)
		return response.InternalServerError(c, "gagal mengambil data akun Google")
	}
	log.Println("GoogleUser:", googleUser)

	// login atau register otomatis
	result, err := h.authService.LoginGoogle(c.Context(), googleUser)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    result.AccessToken,
		HTTPOnly: true,   
		Secure:   false, 
		SameSite: "Lax", 
		Path:     "/",
		MaxAge:   60 * 60 * 24, 
	})

	return c.Redirect().To("https://empathify.vercel.app/dashboard/beranda")
}

func generateState() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// @Summary      Ambil token setelah Google OAuth
// @Description  Ambil access token dari cookie setelah redirect Google OAuth
// @Tags         Auth
// @Produce      json
// @Success      200  {object}  response.Response{data=response.AuthResponse}
// @Failure      401  {object}  response.Response
// @Router       /auth/me [get]
func (h *AuthHandler) Me(c fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		return response.Unauthorized(c, "tidak ada token")
	}

	// parse token untuk ambil user_id
	claims, err := jwt.ParseToken(token)
	if err != nil {
		return response.Unauthorized(c, "token tidak valid")
	}

	// ambil data user
	user, err := h.authService.GetUserByID(c.Context(), claims.UserID)
	if err != nil {
		return response.InternalServerError(c, "terjadi kesalahan, coba lagi")
	}

	// hapus cookie setelah diambil
	c.Cookie(&fiber.Cookie{
		Name:   "access_token",
		Value:  "",
		MaxAge: -1,
	})

	return response.Success(c, "berhasil", fiber.Map{
		"access_token": token,
		"user":         user,
	})
}