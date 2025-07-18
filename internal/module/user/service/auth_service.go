package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	userdto "github.com/algorithm9/flash-deal/internal/module/user/dto"
	"github.com/algorithm9/flash-deal/internal/module/user/repository"
	"github.com/algorithm9/flash-deal/internal/shared/ent/gen"
	"github.com/algorithm9/flash-deal/internal/shared/redisclient"
	"github.com/algorithm9/flash-deal/internal/shared/sms"
	"github.com/algorithm9/flash-deal/pkg/errorx"
)

type AuthService interface {
	GenerateCaptcha(ctx context.Context) (string, string, error)
	VerifyCaptcha(ctx context.Context, id, captcha string, clear bool) bool
	SendVerificationCode(ctx context.Context, phone string) error
	Register(ctx context.Context, phone, password, code string) (uint64, error)
}

type authServiceImpl struct {
	userRepo    repository.UserRepository
	captcha     *Captcha
	redisClient *redisclient.Client
	smsClient   *sms.SMS
}

func NewAuthService(userRep repository.UserRepository, captcha *Captcha, client *redisclient.Client) AuthService {
	return &authServiceImpl{userRepo: userRep, captcha: captcha, redisClient: client}
}

func (auth *authServiceImpl) GenerateCaptcha(ctx context.Context) (string, string, error) {
	id, b64s, _, err := auth.captcha.Generate()
	if err != nil {
		return "", "", err
	}
	return id, b64s, nil
}

func (auth *authServiceImpl) VerifyCaptcha(ctx context.Context, id, captcha string, clear bool) bool {
	return auth.captcha.Verify(id, captcha, clear)
}

func (auth *authServiceImpl) SendVerificationCode(ctx context.Context, phone string) error {
	// Limit to one SMS per minute
	keyLimit := fmt.Sprintf("sms_limit:%s", phone)
	if _, err := auth.redisClient.Client.Get(ctx, keyLimit).Result(); err == nil {
		return errorx.New(
			errorx.CodeLimitRequest.Int(),
			http.StatusTooManyRequests,
			"Too many requests",
		)
	}

	// Check whether the user is already registered
	user, err := auth.userRepo.GetByPhone(ctx, phone)
	if err != nil && !gen.IsNotFound(err) {
		return errorx.Wrap(
			userdto.QueryUserError.Int(),
			http.StatusInternalServerError,
			"Error in querying user", err,
		)
	}

	if user != nil {
		return errorx.New(
			userdto.UserAlreadyExist.Int(),
			http.StatusConflict,
			"User already exists",
		)
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	key := fmt.Sprintf("sms:%s", phone)
	auth.redisClient.Client.Set(ctx, key, code, 5*time.Minute)
	auth.redisClient.Client.Set(ctx, keyLimit, 1, time.Minute)

	if err := auth.smsClient.Send(phone, code); err != nil {
		return errorx.Wrap(
			userdto.FailedSendSMS.Int(),
			http.StatusInternalServerError,
			"Failed to send SMS",
			err,
		)
	}

	return nil
}

func (auth *authServiceImpl) Register(ctx context.Context, phone, password, smsCode string) (uint64, error) {

	key := fmt.Sprintf("sms:%s", phone)
	code, err := auth.redisClient.Client.Get(ctx, key).Result()
	if err != nil {
		return 0, errorx.New(
			userdto.FailedFindSMSCode.Int(),
			http.StatusBadRequest,
			"Failed to find code",
		)
	}

	if code != smsCode {
		return 0, errorx.Wrap(
			userdto.InvalidSMSCode.Int(),
			http.StatusBadRequest,
			"Invalid SMS code",
			err,
		)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, errorx.Wrap(
			userdto.FailedGenerateHash.Int(),
			http.StatusInternalServerError,
			"Failed to generate hash",
			err,
		)
	}

	user, err := auth.userRepo.Create(ctx, phone, string(hash))
	if err != nil {
		return 0, errorx.Wrap(
			userdto.FailedCreateUser.Int(),
			http.StatusInternalServerError,
			"Failed to create user",
			err,
		)
	}

	auth.redisClient.Client.Del(ctx, key)

	return user.ID, nil
}
