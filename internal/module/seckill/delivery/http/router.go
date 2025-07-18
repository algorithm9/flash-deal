package seckillhttp

import (
	"github.com/gin-gonic/gin"

	"github.com/algorithm9/flash-deal/internal/shared/middleware"
)

type SeckillRouter struct {
	seckillHandler *SeckillHandler
}

func NewSeckillRouter(seckillHandler *SeckillHandler) *SeckillRouter {
	return &SeckillRouter{seckillHandler: seckillHandler}
}

func (r *SeckillRouter) RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthJWT) {
	api := router.Group("/v1/seckill")
	{
		api.GET("/activities", r.seckillHandler.GetSeckillActivities)
		api.GET("/activities/sku/:id", r.seckillHandler.GetSeckillActivityDetail)
		api.POST("/:activity_id/:sku_id/request", authMiddleware.MiddlewareFunc(), r.seckillHandler.SeckillRequest)
		api.GET("/:activity_id/:sku_id/result", authMiddleware.MiddlewareFunc(), r.seckillHandler.SeckillResult)
	}
}
