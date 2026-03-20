package provider

import (
	"context"
	"encoding/json"
	"errors"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/Empathify-FICPACT/Empathify-BE/config"
	"github.com/Empathify-FICPACT/Empathify-BE/internal/domain"
)

var (
	ErrFailedExchangeCode  = errors.New("failed to exchange code")
	ErrFailedGetUserInfo   = errors.New("failed to get user info")
)

type GoogleOAuthProvider struct {
	config *oauth2.Config
}

func NewGoogleOAuthProvider() *GoogleOAuthProvider {
	cfg := &oauth2.Config{
		ClientID:     config.App.GoogleClientID,
		ClientSecret: config.App.GoogleClientSecret,
		RedirectURL:  config.App.GoogleRedirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthProvider{config: cfg}
}

func (g *GoogleOAuthProvider) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (g *GoogleOAuthProvider) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, ErrFailedExchangeCode
	}
	return token, nil
}

func (g *GoogleOAuthProvider) GetUserInfo(ctx context.Context, token *oauth2.Token) (*domain.GoogleUser, error) {
	client := g.config.Client(ctx, token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, ErrFailedGetUserInfo
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, ErrFailedGetUserInfo
	}

	var raw struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Gender string `json:"gender"`
	}

	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, ErrFailedGetUserInfo
	}

	return &domain.GoogleUser{
		ID:     raw.ID,
		Name:   raw.Name,
		Email:  raw.Email,
		Gender: raw.Gender,
	}, nil
}