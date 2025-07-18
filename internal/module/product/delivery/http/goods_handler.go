package producthttp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/algorithm9/flash-deal/internal/module/product/service"
	"github.com/algorithm9/flash-deal/pkg/errorx"
	"github.com/algorithm9/flash-deal/pkg/response"
)

type GoodsHandler struct {
	svc service.GoodsService
}

func NewGoodsHandler(svc service.GoodsService) *GoodsHandler {
	return &GoodsHandler{svc: svc}
}

func (h *GoodsHandler) GetSKUDetail(c *gin.Context) {
	skuIDStr := c.Param("sku_id")
	skuID, err := strconv.ParseUint(skuIDStr, 10, 64)
	if err != nil {
		response.Fail(c, errorx.New(http.StatusBadRequest, http.StatusBadRequest, "invalid sku_id"))
		return
	}

	detail, err := h.svc.GetSKUDetail(c.Request.Context(), skuID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get detail failed"})
		response.Fail(c, errorx.Wrap(http.StatusInternalServerError, http.StatusInternalServerError, "get detail failed", err))
		return
	}
	response.Success(c, detail)
}
