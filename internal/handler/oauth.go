package handler

import (
	"net/http"

	"github.com/ali/sso-server/internal/model"
	"github.com/ali/sso-server/pkg/logger"
	"github.com/labstack/echo/v4"
)

type OAuthHandler struct {
	// TODO: add oauth service
}

func NewOAuthHandler() *OAuthHandler {
	return &OAuthHandler{}
}

// Authorize godoc
// @Summary OAuth2 authorization endpoint
// @Tags oauth
// @Param client_id query string true "Client ID"
// @Param redirect_uri query string true "Redirect URI"
// @Param response_type query string true "Response type (code)"
// @Param scope query string false "Requested scope"
// @Param state query string true "State parameter"
// @Success 302
// @Failure 400 {object} ErrorResponse
// @Router /oauth/authorize [get]
func (h *OAuthHandler) Authorize(c echo.Context) error {
	clientID := c.QueryParam("client_id")
	redirectURI := c.QueryParam("redirect_uri")
	responseType := c.QueryParam("response_type")
	scope := c.QueryParam("scope")
	state := c.QueryParam("state")

	if clientID == "" || redirectURI == "" || responseType == "" || state == "" {
		return badRequest(c, "missing required parameters")
	}

	if responseType != "code" {
		return badRequest(c, "unsupported response_type")
	}

	// TODO: validate client_id
	// TODO: validate redirect_uri
	// TODO: check if user is logged in
	// TODO: show consent page or redirect with code

	logger.Info("oauth authorize request",
		"client_id", clientID,
		"redirect_uri", redirectURI,
		"scope", scope,
	)

	// Placeholder: redirect with authorization code
	code := "authorization_code_placeholder"
	return c.Redirect(http.StatusFound, redirectURI+"?code="+code+"&state="+state)
}

// Token godoc
// @Summary OAuth2 token endpoint
// @Tags oauth
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "Grant type"
// @Param code formData string false "Authorization code"
// @Param redirect_uri formData string false "Redirect URI"
// @Param client_id formData string true "Client ID"
// @Param client_secret formData string true "Client secret"
// @Param refresh_token formData string false "Refresh token"
// @Success 200 {object} model.TokenResponse
// @Failure 400 {object} OAuthErrorResponse
// @Failure 401 {object} OAuthErrorResponse
// @Router /oauth/token [post]
func (h *OAuthHandler) Token(c echo.Context) error {
	grantType := c.FormValue("grant_type")
	clientID := c.FormValue("client_id")
	clientSecret := c.FormValue("client_secret")

	if clientID == "" || clientSecret == "" {
		return oauthError(c, "invalid_client", "client credentials required")
	}

	// TODO: validate client credentials

	switch grantType {
	case "authorization_code":
		return h.handleAuthorizationCode(c)
	case "refresh_token":
		return h.handleRefreshToken(c)
	default:
		return oauthError(c, "unsupported_grant_type", "grant type not supported")
	}
}

func (h *OAuthHandler) handleAuthorizationCode(c echo.Context) error {
	code := c.FormValue("code")
	redirectURI := c.FormValue("redirect_uri")

	if code == "" || redirectURI == "" {
		return oauthError(c, "invalid_request", "code and redirect_uri required")
	}

	// TODO: validate authorization code
	// TODO: validate redirect_uri matches
	// TODO: generate tokens
	// TODO: delete authorization code

	logger.Info("oauth token exchange", "grant_type", "authorization_code")

	return c.JSON(http.StatusOK, model.TokenResponse{
		AccessToken:  "access_token_placeholder",
		RefreshToken: "refresh_token_placeholder",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

func (h *OAuthHandler) handleRefreshToken(c echo.Context) error {
	refreshToken := c.FormValue("refresh_token")

	if refreshToken == "" {
		return oauthError(c, "invalid_request", "refresh_token required")
	}

	// TODO: validate refresh token
	// TODO: generate new tokens

	logger.Info("oauth token refresh", "grant_type", "refresh_token")

	return c.JSON(http.StatusOK, model.TokenResponse{
		AccessToken:  "new_access_token_placeholder",
		RefreshToken: "new_refresh_token_placeholder",
		TokenType:    "Bearer",
		ExpiresIn:    3600,
	})
}

// Revoke godoc
// @Summary Revoke OAuth2 token
// @Tags oauth
// @Accept application/x-www-form-urlencoded
// @Param token formData string true "Token to revoke"
// @Param token_type_hint formData string false "Token type hint (access_token or refresh_token)"
// @Success 200
// @Router /oauth/revoke [post]
func (h *OAuthHandler) Revoke(c echo.Context) error {
	token := c.FormValue("token")
	tokenTypeHint := c.FormValue("token_type_hint")

	if token == "" {
		return oauthError(c, "invalid_request", "token required")
	}

	// TODO: revoke token

	logger.Info("oauth token revoked", "token_type_hint", tokenTypeHint)

	return c.NoContent(http.StatusOK)
}

// UserInfo godoc
// @Summary Get user info (OpenID Connect)
// @Tags oauth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} UserInfoResponse
// @Failure 401 {object} OAuthErrorResponse
// @Router /oauth/userinfo [get]
func (h *OAuthHandler) UserInfo(c echo.Context) error {
	// TODO: get user from access token

	logger.Debug("oauth userinfo request")

	return c.JSON(http.StatusOK, UserInfoResponse{
		Sub:   "user_id_placeholder",
		Email: "user@example.com",
		Name:  "John Doe",
	})
}

type OAuthErrorResponse struct {
	Error       string `json:"error"`
	Description string `json:"error_description,omitempty"`
}

type UserInfoResponse struct {
	Sub   string `json:"sub"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
}

func oauthError(c echo.Context, err, description string) error {
	return c.JSON(http.StatusBadRequest, OAuthErrorResponse{
		Error:       err,
		Description: description,
	})
}
