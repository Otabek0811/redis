package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Create product godoc
// @ID create_product
// @Router /product [POST]
// @Summary Create Product
// @Description Create Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param product body models.CreateProduct true "CreateProductRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateProduct(c *gin.Context) {

	var createProduct models.CreateProduct
	err := c.ShouldBindJSON(&createProduct)
	if err != nil {
		h.handlerResponse(c, "error product should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Product().Create(c.Request.Context(), &createProduct)
	if err != nil {
		h.handlerResponse(c, "storage.product.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Product().GetByID(c.Request.Context(), &models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.product.getById", http.StatusInternalServerError, err.Error())
		return
	}

	err = h.cache.Product().Delete()
	if err != nil {
		h.handlerResponse(c, "cache.Product.Delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create product resposne", http.StatusCreated, resp)
}

// GetByID product godoc
// @ID get_by_id_product
// @Router /product/{id} [GET]
// @Summary Get By ID Product
// @Description Get By ID Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdProduct(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Product().GetByID(c.Request.Context(), &models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.product.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id product resposne", http.StatusOK, resp)
}

// GetList product godoc
// @ID get_list_product
// @Router /product [GET]
// @Summary Get List Product
// @Description Get List Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListProduct(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list product offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list product limit", http.StatusBadRequest, "invalid limit")
		return
	}

	exists, err := h.cache.Product().Exists()
	if err != nil {
		h.handlerResponse(c, "exists check products", http.StatusInternalServerError, err.Error())
		return
	}

	if !exists {
		fmt.Println("POSTGRES")

		pgTime := time.Now()
		resp, err := h.strg.Product().GetList(c.Request.Context(), &models.ProductGetListRequest{
			Offset: offset,
			Limit:  limit,
			Search: c.Query("search"),
		})
		if err != nil {
			h.handlerResponse(c, "storage.product.get_list", http.StatusInternalServerError, err.Error())
			return
		}

		fmt.Println("TIME:", time.Since(pgTime))

		writeTime := time.Now()
		err = h.cache.Product().CreateList(resp)
		if err != nil {
			h.handlerResponse(c, "cache.Product.CreateList", http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Println("WRITE TIME:", time.Since(writeTime))

		h.handlerResponse(c, "get list product resposne", http.StatusOK, resp)
		return
	}

	fmt.Println("REDIS")
	rdTime := time.Now()
	resp, err := h.cache.Product().GetList()
	if err != nil {
		h.handlerResponse(c, "cache.Product.GetList", http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("TIME:", time.Since(rdTime))

	h.handlerResponse(c, "get list product resposne", http.StatusOK, resp)
	return
}

// Update product godoc
// @ID update_product
// @Router /product/{id} [PUT]
// @Summary Update Product
// @Description Update Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param product body models.UpdateProduct true "UpdateProductRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateProduct(c *gin.Context) {

	var (
		id            string = c.Param("id")
		updateProduct models.UpdateProduct
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateProduct)
	if err != nil {
		h.handlerResponse(c, "error product should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateProduct.Id = id
	rowsAffected, err := h.strg.Product().Update(c.Request.Context(), &updateProduct)
	if err != nil {
		h.handlerResponse(c, "storage.product.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.product.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Product().GetByID(c.Request.Context(), &models.ProductPrimaryKey{Id: updateProduct.Id})
	if err != nil {
		h.handlerResponse(c, "storage.product.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create product resposne", http.StatusAccepted, resp)
}

// Patch product godoc
// @ID patch_product
// @Router /product/{id} [PATCH]
// @Summary Patch Product
// @Description Patch Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param product body models.PatchRequest true "PatchProductRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) PatchProduct(c *gin.Context) {

	var (
		id           string = c.Param("id")
		patchProduct models.PatchRequest
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&patchProduct)
	if err != nil {
		h.handlerResponse(c, "error product should bind json", http.StatusBadRequest, err.Error())
		return
	}

	patchProduct.ID = id
	rowsAffected, err := h.strg.Product().Patch(c.Request.Context(), &patchProduct)
	if err != nil {
		h.handlerResponse(c, "storage.product.patch", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.product.patch", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Product().GetByID(c.Request.Context(), &models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.product.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create product resposne", http.StatusAccepted, resp)
}

// Delete product godoc
// @ID delete_product
// @Router /product/{id} [DELETE]
// @Summary Delete Product
// @Description Delete Product
// @Tags Product
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteProduct(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Product().Delete(c.Request.Context(), &models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.product.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create product resposne", http.StatusNoContent, nil)
}
