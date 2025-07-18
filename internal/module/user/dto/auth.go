package userdto

type CaptchaResponse struct {
	CaptchaID string `json:"captchaId"`
	Image     string `json:"image"`
}

// VerificationRequest 验证码请求
type VerificationRequest struct {
	Phone     string `json:"phone" binding:"required,e164"`
	CaptchaID string `json:"captcha_id" binding:"required"`
	Captcha   string `json:"captcha_code" binding:"required"`
}

type SMSResponse struct {
	Message string `json:"message"`
}

// UserRegisterRequest 用户注册请求
type UserRegisterRequest struct {
	Password string `json:"password" binding:"required,min=6,max=30" example:"P@ssw0rd"` // 密码
	Phone    string `json:"phone" binding:"required,e164" example:"+8613812345678"`      // 手机号
	Code     string `json:"code" binding:"required" example:"123456"`                    // 手机验证码
}

// UserRegisterResponse 用户注册响应
type UserRegisterResponse struct {
	UserID   uint64 `json:"user_id" example:"123"`
	Token    string `json:"token"`
	ExpireAt string `json:"expire_at"`
}

type UserLoginRequest struct {
	Phone    string `json:"phone" binding:"required,e164"`
	Password string `json:"password" binding:"required,min=6,max=30" example:"123456"`
}

type UserLoginResponse struct {
	Token    string `json:"token"`
	ExpireAt string `json:"expire_at"`
}

type UserModuleResponseCode int

const (
	QueryUserError      UserModuleResponseCode = 10001
	UserAlreadyExist    UserModuleResponseCode = 10002
	FailedSendSMS       UserModuleResponseCode = 10003
	FailedFindSMSCode   UserModuleResponseCode = 10004
	InvalidSMSCode      UserModuleResponseCode = 10005
	FailedGenerateHash  UserModuleResponseCode = 10006
	FailedCreateUser    UserModuleResponseCode = 10007
	FailedGenerateToken UserModuleResponseCode = 10008
	WrongPassword       UserModuleResponseCode = 10009
)

func (r UserModuleResponseCode) Int() int {
	return int(r)
}
