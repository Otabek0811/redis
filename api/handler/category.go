package handler

import (
	"app/api/models"
	"app/pkg/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Create category godoc
// @ID create_category
// @Router /category [POST]
// @Summary Create Category
// @Description Create Category
// @Tags Category
// @Accept json
// @Procedure json
// @Param category body models.CreateCategory true "CreateCategoryRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) CreateCategory(c *gin.Context) {

	var createCategory models.CreateCategory
	err := c.ShouldBindJSON(&createCategory)
	if err != nil {
		h.handlerResponse(c, "error category should bind json", http.StatusBadRequest, err.Error())
		return
	}

	id, err := h.strg.Category().Create(c.Request.Context(), &createCategory)
	if err != nil {
		h.handlerResponse(c, "storage.category.create", http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := h.strg.Category().GetByID(c.Request.Context(), &models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.category.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create category resposne", http.StatusCreated, resp)
}

// GetByID category godoc
// @ID get_by_id_category
// @Router /category/{id} [GET]
// @Summary Get By ID Category
// @Description Get By ID Category
// @Tags Category
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetByIdCategory(c *gin.Context) {

	var id string = c.Param("id")
	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	resp, err := h.strg.Category().GetByID(c.Request.Context(), &models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.category.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get by id category resposne", http.StatusOK, resp)
}

// GetList category godoc
// @ID get_list_category
// @Router /category [GET]
// @Summary Get List Category
// @Description Get List Category
// @Tags Category
// @Accept json
// @Procedure json
// @Param offset query string false "offset"
// @Param limit query string false "limit"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) GetListCategory(c *gin.Context) {

	offset, err := h.getOffsetQuery(c.Query("offset"))
	if err != nil {
		h.handlerResponse(c, "get list category offset", http.StatusBadRequest, "invalid offset")
		return
	}

	limit, err := h.getLimitQuery(c.Query("limit"))
	if err != nil {
		h.handlerResponse(c, "get list category limit", http.StatusBadRequest, "invalid limit")
		return
	}

	resp, err := h.strg.Category().GetList(c.Request.Context(), &models.CategoryGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: c.Query("search"),
	})
	if err != nil {
		h.handlerResponse(c, "storage.category.get_list", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "get list category resposne", http.StatusOK, resp)
}

// Update category godoc
// @ID update_category
// @Router /category/{id} [PUT]
// @Summary Update Category
// @Description Update Category
// @Tags Category
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Param category body models.UpdateCategory true "UpdateCategoryRequest"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) UpdateCategory(c *gin.Context) {

	var (
		id             string = c.Param("id")
		updateCategory models.UpdateCategory
	)

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := c.ShouldBindJSON(&updateCategory)
	if err != nil {
		h.handlerResponse(c, "error category should bind json", http.StatusBadRequest, err.Error())
		return
	}

	updateCategory.Id = id
	rowsAffected, err := h.strg.Category().Update(c.Request.Context(), &updateCategory)
	if err != nil {
		h.handlerResponse(c, "storage.category.update", http.StatusInternalServerError, err.Error())
		return
	}

	if rowsAffected <= 0 {
		h.handlerResponse(c, "storage.category.update", http.StatusBadRequest, "now rows affected")
		return
	}

	resp, err := h.strg.Category().GetByID(c.Request.Context(), &models.CategoryPrimaryKey{Id: updateCategory.Id})
	if err != nil {
		h.handlerResponse(c, "storage.category.getById", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create category resposne", http.StatusAccepted, resp)
}

// Delete category godoc
// @ID delete_category
// @Router /category/{id} [DELETE]
// @Summary Delete Category
// @Description Delete Category
// @Tags Category
// @Accept json
// @Procedure json
// @Param id path string true "id"
// @Success 200 {object} Response{data=string} "Success Request"
// @Response 400 {object} Response{data=string} "Bad Request"
// @Failure 500 {object} Response{data=string} "Server error"
func (h *handler) DeleteCategory(c *gin.Context) {

	var id string = c.Param("id")

	if !helper.IsValidUUID(id) {
		h.handlerResponse(c, "is valid uuid", http.StatusBadRequest, "invalid id")
		return
	}

	err := h.strg.Category().Delete(c.Request.Context(), &models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(c, "storage.category.delete", http.StatusInternalServerError, err.Error())
		return
	}

	h.handlerResponse(c, "create category resposne", http.StatusNoContent, nil)
}
