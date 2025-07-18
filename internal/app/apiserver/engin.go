package apiserver

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"

	"github.com/algorithm9/flash-deal/internal/model"
	producthttp "github.com/algorithm9/flash-deal/internal/module/product/delivery/http"
	seckillhttp "github.com/algorithm9/flash-deal/internal/module/seckill/delivery/http"
	userhttp "github.com/algorithm9/flash-deal/internal/module/user/delivery/http"
	"github.com/algorithm9/flash-deal/internal/shared/constant/profile"
	"github.com/algorithm9/flash-deal/internal/shared/middleware"

	ginSwagger "github.com/swaggo/gin-swagger"

	swagger "github.com/algorithm9/flash-deal/api"
)

func NewGinEngine(
	cfg *model.ServerConfig,
	userRouter *userhttp.UserRouter,
	productRouter *producthttp.GoodsRouter,
	seckillRouter *seckillhttp.SeckillRouter,
	authMiddleware *middleware.AuthJWT,
) *gin.Engine {
	gin.SetMode(cfg.Env)
	engine := gin.New()
	// 注册中间件
	engine.Use(
		middleware.TraceAndLogger(),
		middleware.RecoveryWithLog(),
		middleware.CORS,
	)

	// 注册路由
	api := engine.Group("/api")
	userRouter.RegisterRoutes(api, authMiddleware)
	productRouter.RegisterRoutes(api, authMiddleware)
	seckillRouter.RegisterRoutes(api, authMiddleware)

	// use ginSwagger middleware to serve the API docs
	if cfg.Env != profile.Prod {
		swagger.SwaggerInfo.BasePath = "/api"
		api.GET("/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		pprofRoutes := api.Group("/debug/pprof")
		{
			pprofRoutes.GET("/", gin.WrapH(http.DefaultServeMux))                      // 首页展示信息
			pprofRoutes.GET("/heap", gin.WrapH(pprof.Handler("heap")))                 // 堆内存分析
			pprofRoutes.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))       // Goroutine 状态
			pprofRoutes.GET("/block", gin.WrapH(pprof.Handler("block")))               // 阻塞分析
			pprofRoutes.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate"))) // 线程创建分析
			pprofRoutes.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))               // mutex分析
		}
	}

	return engine
}
