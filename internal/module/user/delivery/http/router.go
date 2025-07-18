package userhttp

import (
	"github.com/gin-gonic/gin"
	//"ecommerce-mvp/internal/module/user/service"
	//"ecommerce-mvp/internal/shared/middleware"

	"github.com/algorithm9/flash-deal/internal/shared/middleware"
)

// UserRouter 用户模块路由控制器
type UserRouter struct {
	authHandler *AuthHandler
	//userHandler   *UserHandler
	//deviceHandler *DeviceHandler
}

// NewUserRouter 创建用户路由实例
func NewUserRouter(
	authHandle *AuthHandler,
) *UserRouter {
	return &UserRouter{
		authHandler: authHandle,
	}
}

//// NewUserRouter 创建用户路由实例
//func NewUserRouter(
//	authService service.AuthService,
//	userService service.UserService,
//	deviceService service.DeviceService,
//	jwtUtil *middleware.JWTUtil,
//) *UserRouter {
//	return &UserRouter{
//		authHandler:   NewAuthHandler(authService, jwtUtil),
//		userHandler:   NewUserHandler(userService),
//		deviceHandler: NewDeviceHandler(deviceService),
//	}
//}

// RegisterRoutes 注册用户模块所有路由
func (r *UserRouter) RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthJWT) {
	// 不需要认证的路由
	public := router.Group("/v1/users")
	{
		public.GET("/captcha", r.authHandler.GenerateCaptcha)            // 获取captcha
		public.POST("/sms/send", r.authHandler.RequestVerificationCode)  //请求短信验证码
		public.POST("/register", r.authHandler.Register(authMiddleware)) // 用户注册
		public.POST("/login", r.authHandler.Login(authMiddleware))       // 用户登陆
		public.POST("/logout", r.authHandler.Logout(authMiddleware))
		public.POST("/token/refresh", r.authHandler.RefreshToken(authMiddleware))
		//	public.POST("/reset-password", r.authHandler.RequestPasswordReset) // 请求重置密码
		//	public.PUT("/reset-password", r.authHandler.ConfirmPasswordReset)  // 确认重置密码
	}
	//
	//// 需要双因素认证的路由 (特殊中间件)
	//twoFA := router.Group("/user")
	//twoFA.Use(middleware.TempAuthMiddleware())
	//{
	//	twoFA.POST("/verify-2fa", r.authHandler.Verify2FA)            // 验证双因素
	//	twoFA.POST("/recover-account", r.authHandler.AccountRecovery) // 账号恢复
	//}
	//
	//// 需要完整认证的路由
	//authRequired := router.Group("/user")
	//authRequired.Use(middleware.JWTAuthMiddleware())
	//{
	//	// 用户资料管理
	//	authRequired.GET("/profile", r.userHandler.GetProfile)      // 获取个人资料
	//	authRequired.PUT("/profile", r.userHandler.UpdateProfile)   // 更新个人资料
	//	authRequired.PUT("/password", r.userHandler.ChangePassword) // 修改密码
	//	authRequired.POST("/logout", r.authHandler.Logout)          // 用户登出
	//
	//	// 双因素认证管理
	//	authRequired.POST("/2fa/enable", r.authHandler.Enable2FA)           // 启用双因素
	//	authRequired.POST("/2fa/disable", r.authHandler.Disable2FA)         // 禁用双因素
	//	authRequired.GET("/2fa/backup-codes", r.authHandler.GetBackupCodes) // 获取备用码
	//
	//	// 设备管理
	//	authRequired.GET("/devices", r.deviceHandler.ListDevices)         // 设备列表
	//	authRequired.DELETE("/devices/:id", r.deviceHandler.RemoveDevice) // 移除设备
	//	authRequired.POST("/devices/trust", r.deviceHandler.TrustDevice)  // 信任设备
	//}
}
