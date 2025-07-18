package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"github.com/algorithm9/flash-deal/internal/model"
	userdto "github.com/algorithm9/flash-deal/internal/module/user/dto"
	"github.com/algorithm9/flash-deal/internal/module/user/repository"
	"github.com/algorithm9/flash-deal/pkg/errorx"
	"github.com/algorithm9/flash-deal/pkg/logger"
	"github.com/algorithm9/flash-deal/pkg/response"
)

type AuthJWT struct {
	IdentityKey string
	*jwt.GinJWTMiddleware
}

func NewAuthJWTMiddleware(cfg *model.JWT, userRepository repository.UserRepository) (*AuthJWT, error) {
	timeFunc := func() time.Time {
		return time.Now().UTC()
	}

	unauthorized := func(c *gin.Context, code int, message string) {
		c.JSON(code, gin.H{
			"code":    "unauthorized",
			"message": message,
		})
	}

	mw, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:             cfg.Realm,
		SigningAlgorithm:  cfg.SigningAlgorithm,
		Key:               []byte(cfg.Key),
		Timeout:           time.Duration(cfg.Timeout) * time.Hour,
		MaxRefresh:        time.Duration(cfg.MaxRefresh) * time.Hour,
		TimeFunc:          timeFunc,
		SendCookie:        true,
		SendAuthorization: true,
		IdentityHandler:   identityHandler(cfg.IdentityKey),
		Authenticator:     authenticator(userRepository),
		Authorizator:      authorizator(),
		PayloadFunc:       payloadFunc(cfg.IdentityKey),
		Unauthorized:      unauthorized,
		LoginResponse:     loginResponse,
		LogoutResponse:    logoutResponse,
		RefreshResponse:   refreshResponse,
		TokenLookup:       "header: Authorization, query: token, cookie: jwt",
		TokenHeadName:     "Bearer",
	})

	if err != nil {
		return nil, errorx.Wrap(
			http.StatusInternalServerError,
			http.StatusInternalServerError,
			"failed to create auth jwt middleware",
			err,
		)
	}
	return &AuthJWT{
		IdentityKey:      cfg.IdentityKey,
		GinJWTMiddleware: mw,
	}, nil
}

type User struct {
	UserID string
	Phone  string
}

// authenticator verify the user credentials(i.e. password matches hashed password for a given user email, and any other authentication logic).
// Then the authenticator should return a struct or map that contains the user data that will be embedded in the jwt token.
// the data returned from the authenticator is passed into the payloadFunc.
// If an error is returned, the unauthorized function is used.
func authenticator(repository repository.UserRepository) func(c *gin.Context) (interface{}, error) {
	return func(c *gin.Context) (interface{}, error) {
		var payload userdto.UserLoginRequest
		if err := c.ShouldBindJSON(&payload); err != nil {
			return nil, errorx.Wrap(
				http.StatusBadRequest,
				http.StatusBadRequest,
				"invalid payload",
				err,
			)
		}

		userByPhone, err := repository.GetByPhone(c, payload.Phone)
		if err != nil {
			return nil, errorx.Wrap(
				http.StatusInternalServerError,
				http.StatusInternalServerError, "failed to query user",
				err,
			)
		}

		err = bcrypt.CompareHashAndPassword([]byte(userByPhone.PasswordHash), []byte(payload.Password))
		if err != nil {
			return 0, errorx.New(
				userdto.WrongPassword.Int(),
				http.StatusUnauthorized, "Wrong password",
			)
		}

		return &User{
			UserID: strconv.FormatUint(userByPhone.ID, 10),
			Phone:  userByPhone.Phone,
		}, nil
	}
}

// payloadFunc convert data from authenticator into MapClaims.
// The elements of MapClaims returned in payloadFunc will be embedded within the jwt token (as token claims)
func payloadFunc(identityKey string) func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				identityKey: v.UserID,
				"phone":     v.Phone,
			}
		}
		return jwt.MapClaims{}
	}
}

// identityHandler fetch the user identity from claims embedded within the jwt token,
// and pass this identity value to authorizator.
func identityHandler(identityKey string) func(c *gin.Context) interface{} {
	return func(c *gin.Context) interface{} {
		claims := jwt.ExtractClaims(c)
		fmt.Println(claims)
		return &User{
			UserID: claims[identityKey].(string),
			Phone:  claims["phone"].(string),
		}
	}
}

// authorizator check if the user is authorized to be reaching this endpoint
func authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {
		user, ok := data.(*User)
		if ok {
			logger.L().Trace().Msgf("current user, id: %s, phone: %s", user.UserID, user.Phone)
			userID, err := strconv.ParseUint(user.UserID, 10, 64)
			if err != nil {
				logger.L().Err(err).Msgf("failed to parse user id id: %s", user.UserID)
				return false
			}
			c.Set("user_id", userID)
			c.Set("phone", user.Phone)
			return true
		}
		return false
	}
}

func loginResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "Login successful",
		Data: userdto.UserLoginResponse{
			Token:    token,
			ExpireAt: expire.Format(time.RFC3339),
		},
	})
}

func logoutResponse(c *gin.Context, code int) {
	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "Logout successful",
	})
}

func refreshResponse(c *gin.Context, code int, token string, expire time.Time) {
	c.JSON(http.StatusOK, response.Response{
		Code:    0,
		Message: "Token refreshed",
		Data: userdto.UserLoginResponse{
			Token:    token,
			ExpireAt: expire.Format(time.RFC3339),
		},
	})
}
