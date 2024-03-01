package handler

import (
	"app/api/models"
	"app/config"
	"app/pkg/helper"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create market godoc
// @ID create_market
// @Router /v1/market [POST]
// @Summary Create Market
// @Description Create Market
// @Tags Market
// @Accept json
// @Procedure json
// @Param market body models.CreateMarket true "CreateMarketRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateMarket(c *gin.Context) {

	var createMarket models.CreateMarket
	err := c.ShouldBindJSON(&createMarket)
	if err != nil {
		h.handlerResponse(c, "error market should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Market().Create(c.Request.Context(), &createMarket)
	if err != nil {
		h.handlerResponse(c, "storage.market.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Market().GetByID(c.Request.Context(), &models.MarketPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.market.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create market resposne", http.StatusCreated, resp)
}

// GetByID market godoc
// @ID get_by_id_market
// @Router /v1/market/{id} [GET]
// @Summary Get By ID Market
// @Description Get By ID Market
// @Tags Market
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdMarket(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Market().GetByID(c.Request.Context(), &models.MarketPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.market.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id market resposne", http.StatusOK, resp)
}

// @Security ApiKeyAuth
// GetList market godoc
// @ID get_list_market
// @Router /v1/market [GET]
// @Summary Get List Market
// @Description Get List Market
// @Tags Market
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListMarket(c *gin.Context) {

	info, exists := c.Get("Auth")
	if !exists {
		h.handlerResponse(c, "invalid info token", http.StatusBadRequest, "invalid info token")
		return
	}

	user := info.(helper.TokenInfo)
	fmt.Println(user)

	if user.ClientType == config.ClientTypeSuper {
		user.UserID = ""
	}

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list market offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list market limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Market().GetList(c.Request.Context(), &models.MarketGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
		// UserId: userId,
	})
	if err != nil {
		h.handlerResponse(c, "storage.market.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list market resposne", http.StatusOK, resp)
}

// Update market godoc
// @ID update_market
// @Router /v1/market/{id} [PUT]
// @Summary Update Market
// @Description Update Market
// @Tags Market
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param market body models.UpdateMarket true "UpdateMarketRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateMarket(c *gin.Context) {

	var (
		id           string = c.Param("id")
		updateMarket models.UpdateMarket
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateMarket)
	if err != nil {
		h.handlerResponse(c, "error market should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateMarket.Id = id
	rowsAffected, err := h.strg.Market().Update(c.Request.Context(), &updateMarket)
	if err != nil {
		h.handlerResponse(c, "storage.market.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.market.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Market().GetByID(c.Request.Context(), &models.MarketPrimaryKey{Id: updateMarket.Id})
	if err != nil {
		h.handlerResponse(c, "storage.market.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create market resposne", http.StatusAccepted, resp)
}

// Patch market godoc
// @ID patch_market
// @Router /v1/market/{id} [PATCH]
// @Summary Patch Market
// @Description Patch Market
// @Tags Market
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param market body models.PatchRequest true "PatchMarketRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) PatchMarket(c *gin.Context) {

	var (
		id          string = c.Param("id")
		patchMarket models.PatchRequest
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&patchMarket)
	if err != nil {
		h.handlerResponse(c, "error market should bind json", http.StatusBadRequest, err.Error())
		return
	}

	patchMarket.ID = id
	rowsAffected, err := h.strg.Market().Patch(c.Request.Context(), &patchMarket)
	if err != nil {
		h.handlerResponse(c, "storage.market.patch", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.market.patch", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Market().GetByID(c.Request.Context(), &models.MarketPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.market.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create market resposne", http.StatusAccepted, resp)
}

// Delete market godoc
// @ID delete_market
// @Router /v1/market/{id} [DELETE]
// @Summary Delete Market
// @Description Delete Market
// @Tags Market
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteMarket(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Market().Delete(c.Request.Context(), &models.MarketPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.market.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create market resposne", http.StatusNoContent, nil)
}
