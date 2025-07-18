package producthttp

import (
	"github.com/gin-gonic/gin"

	"github.com/algorithm9/flash-deal/internal/shared/middleware"
)

type GoodsRouter struct {
	goodsHandler *GoodsHandler
}

func NewGoodsRouter(goodsHandler *GoodsHandler) *GoodsRouter {
	return &GoodsRouter{goodsHandler: goodsHandler}
}

func (r *GoodsRouter) RegisterRoutes(router *gin.RouterGroup, authMiddleware *middleware.AuthJWT) {
	// 不需要认证的路由
	public := router.Group("/v1/goods")
	{
		public.GET("/:sku_id", r.goodsHandler.GetSKUDetail) // 获取商品详情
	}
}
