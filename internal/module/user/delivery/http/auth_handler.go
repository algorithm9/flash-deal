package userhttp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	userdto "github.com/algorithm9/flash-deal/internal/module/user/dto"
	"github.com/algorithm9/flash-deal/internal/module/user/service"
	"github.com/algorithm9/flash-deal/internal/shared/middleware"
	"github.com/algorithm9/flash-deal/pkg/errorx"
	"github.com/algorithm9/flash-deal/pkg/response"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// GenerateCaptcha godoc
//
//	@Summary		Get graphic verification code
//	@Description	Generate graphic verification code
//	@Tags			User
//	@Produce		json
//	@Success		200	{object}	response.Response{data=userdto.CaptchaResponse}
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/users/captcha [get]
func (auth *AuthHandler) GenerateCaptcha(c *gin.Context) {
	id, b64s, err := auth.authService.GenerateCaptcha(c)
	if err != nil {
		response.Fail(c, errorx.Internal("failed to generate captcha"))
		return
	}
	response.Success(c, userdto.CaptchaResponse{
		CaptchaID: id,
		Image:     b64s,
	},
	)
}

// RequestVerificationCode godoc
//
//	@Summary		request sms code
//	@Description	send sms code
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		userdto.VerificationRequest	true	"request param"
//	@Success		200		{object}	response.Response{data=userdto.SMSResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		409		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/users/sms/send [post]
func (auth *AuthHandler) RequestVerificationCode(c *gin.Context) {
	var req userdto.VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, errorx.BadRequest("Invalid request"))
		return
	}

	// verify graphic code
	if !auth.authService.VerifyCaptcha(c, req.CaptchaID, req.Captcha, true) {
		response.Fail(c, errorx.BadRequest("Invalid captcha code"))
		return
	}

	if err := auth.authService.SendVerificationCode(c, req.Phone); err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, userdto.SMSResponse{Message: "Verification code sent"})
}

// Register godoc
//
//	@Summary		user register
//	@Description	create new user account
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		userdto.UserRegisterRequest	true	"register info"
//	@Success		200		{object}	response.Response{data=userdto.UserRegisterResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/users/register [post]
func (auth *AuthHandler) Register(authMiddleware *middleware.AuthJWT) func(c *gin.Context) {
	return func(c *gin.Context) {
		var req userdto.UserRegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Fail(c, errorx.BadRequest("Invalid request"))
			return
		}

		userID, err := auth.authService.Register(c, req.Phone, req.Password, req.Code)
		if err != nil {
			response.Fail(c, err)
			return
		}

		token, expire, err := authMiddleware.TokenGenerator(&middleware.User{
			UserID: strconv.FormatUint(userID, 10),
			Phone:  req.Phone,
		})
		if err != nil {
			response.Fail(c, errorx.Wrap(
				userdto.FailedGenerateToken.Int(),
				http.StatusInternalServerError,
				"failed to generate token, err:%v",
				err,
			),
			)
			return
		}

		response.Success(c, userdto.UserRegisterResponse{
			UserID:   userID,
			Token:    token,
			ExpireAt: expire.Format(time.RFC3339),
		})
	}
}

// Login godoc
//
//	@Summary		user login
//	@Description	login
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			request	body		userdto.UserLoginRequest	true	"login info"
//	@Success		200		{object}	response.Response{data=userdto.UserRegisterResponse}
//	@Failure		400		{object}	response.Response
//	@Failure		403		{object}	response.Response
//	@Failure		404		{object}	response.Response
//	@Failure		500		{object}	response.Response
//	@Router			/api/v1/users/login [post]
func (auth *AuthHandler) Login(authMiddleware *middleware.AuthJWT) func(c *gin.Context) {
	return authMiddleware.LoginHandler
}

// Logout godoc
//
//	@Summary		user logout
//	@Description	logout
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/users/logout [post]
func (auth *AuthHandler) Logout(authMiddleware *middleware.AuthJWT) func(c *gin.Context) {
	return authMiddleware.LogoutHandler
}

// RefreshToken godoc
//
//	@Summary		user refresh token
//	@Description	refresh token
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=userdto.UserRegisterResponse}
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/users/token/refresh [post]
func (auth *AuthHandler) RefreshToken(authMiddleware *middleware.AuthJWT) func(c *gin.Context) {
	return authMiddleware.RefreshHandler
}
