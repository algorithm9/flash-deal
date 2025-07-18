package seckillhttp

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	seckilldto "github.com/algorithm9/flash-deal/internal/module/seckill/dto"
	"github.com/algorithm9/flash-deal/internal/module/seckill/service"
	"github.com/algorithm9/flash-deal/pkg/errorx"
	"github.com/algorithm9/flash-deal/pkg/response"
)

type SeckillHandler struct {
	service service.Service
}

func NewSeckillHandler(s service.Service) *SeckillHandler {
	return &SeckillHandler{service: s}
}

// GetSeckillActivities godoc
//
//	@Summary		Get all valid seckill activities
//	@Description	Get all valid seckill activities
//	@Tags			Seckill
//	@Produce		json
//	@Success		200	{object}	response.Response{data=seckilldto.SeckillActivities}
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/seckill/activities [get]
func (h *SeckillHandler) GetSeckillActivities(c *gin.Context) {
	activities, err := h.service.GetSeckillActivities(c)
	if err != nil {
		return
	}
	response.Success(c, activities)
}

// GetSeckillActivityDetail godoc
//
//	@Summary		Get detail of seckill activity
//	@Description	Get detail of seckill activity
//	@Tags			Seckill
//	@Produce		json
//	@Success		200	{object}	response.Response{data=seckilldto.SeckillSkuDetail}
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/seckill/activities/sku/:id [get]
func (h *SeckillHandler) GetSeckillActivityDetail(c *gin.Context) {
	skuIDStr := c.Param("id")
	skuID, err := strconv.ParseUint(skuIDStr, 10, 64)
	if err != nil {
		response.Fail(c, errorx.New(http.StatusBadRequest, http.StatusBadRequest, "invalid sku_id"))
		return
	}

	detail, err := h.service.GetSeckillActivityDetail(c, skuID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "get detail failed"})
		response.Fail(c, errorx.Wrap(http.StatusInternalServerError, http.StatusInternalServerError, "get detail failed", err))
		return
	}
	response.Success(c, detail)
}

// SeckillRequest godoc
//
//	@Summary		seckill request
//	@Description	seckill request
//	@Tags			Seckill
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.Response{data=string}
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		409	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/seckill/:activity_id/:sku_id/request [post]
func (h *SeckillHandler) SeckillRequest(c *gin.Context) {
	activityIDStr := c.Param("activity_id")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 64)
	if err != nil {
		response.Fail(c, errorx.New(http.StatusBadRequest, http.StatusBadRequest, "invalid activity id"))
		return
	}
	skuIDStr := c.Param("sku_id")
	skuID, err := strconv.ParseUint(skuIDStr, 10, 64)
	if err != nil {
		response.Fail(c, errorx.New(http.StatusBadRequest, http.StatusBadRequest, "invalid sku_id"))
		return
	}
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, errorx.New(http.StatusBadRequest, http.StatusBadRequest, "no user param exists"))
	}
	userID := userIDStr.(uint64)

	if err := h.service.Seckill(c, userID, activityID, skuID); err != nil {
		response.Fail(c, err)
		return
	}
	response.Success(c, seckilldto.Queueing)
}

// SeckillResult godoc
//
//	@Summary		Get result of seckill
//	@Description	Get result of seckill
//	@Tags			Seckill
//	@Produce		json
//	@Success		200	{object}	response.Response{data=string}
//	@Failure		400	{object}	response.Response
//	@Failure		403	{object}	response.Response
//	@Failure		404	{object}	response.Response
//	@Failure		500	{object}	response.Response
//	@Router			/api/v1/seckill/:activity_id/:sku_id/result [get]
func (h *SeckillHandler) SeckillResult(c *gin.Context) {
	activityIDStr := c.Param("activity_id")
	activityID, err := strconv.ParseUint(activityIDStr, 10, 64)
	if err != nil {
		response.Fail(c, errorx.New(http.StatusBadRequest, http.StatusBadRequest, "invalid activity id"))
		return
	}
	skuIDStr := c.Param("sku_id")
	skuID, err := strconv.ParseUint(skuIDStr, 10, 64)
	if err != nil {
		response.Fail(c, errorx.New(http.StatusBadRequest, http.StatusBadRequest, "invalid sku_id"))
		return
	}
	userIDStr, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, errorx.New(http.StatusBadRequest, http.StatusBadRequest, "no user param exists"))
	}
	userID := userIDStr.(uint64)

	result, err := h.service.GetResult(c.Request.Context(), userID, activityID, skuID)
	if err != nil {
		response.Fail(c, errorx.Wrap(http.StatusInternalServerError, http.StatusInternalServerError, "failed to get seckill result", err))
		return
	}
	response.Success(c, result)
}
