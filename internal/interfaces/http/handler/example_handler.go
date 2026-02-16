package handler

import (
	"strconv"

	"go-ddd-scaffold/internal/application/dto"
	"go-ddd-scaffold/internal/application/service"
	"go-ddd-scaffold/pkg/response"

	"github.com/gin-gonic/gin"
)

// ExampleHandler handles example CRUD endpoints
type ExampleHandler struct {
	svc *service.ExampleAppService
}

// NewExampleHandler creates a new handler
func NewExampleHandler(svc *service.ExampleAppService) *ExampleHandler {
	return &ExampleHandler{svc: svc}
}

// List returns example list
// @Summary  List examples
// @Tags     Example
// @Security Bearer
// @Param    page      query int    false "page"      default(1)
// @Param    page_size query int    false "page size"  default(10)
// @Param    keyword   query string false "search keyword"
// @Param    status    query string false "status filter" Enums(active, inactive)
// @Success  200 {object} response.Response{data=response.PageData}
// @Router   /examples [get]
func (h *ExampleHandler) List(c *gin.Context) {
	var req dto.QueryExampleRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ParamError(c, "invalid parameters")
		return
	}

	items, total, err := h.svc.List(&req)
	if err != nil {
		response.ServerError(c, "query failed")
		return
	}

	response.SuccessPage(c, items, total, req.Page, req.PageSize)
}

// Create creates an example
// @Summary  Create example
// @Tags     Example
// @Security Bearer
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateExampleRequest true "create parameters"
// @Success  200  {object} response.Response{data=dto.ExampleResponse}
// @Router   /examples [post]
func (h *ExampleHandler) Create(c *gin.Context) {
	var req dto.CreateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "invalid parameters: "+err.Error())
		return
	}

	item, err := h.svc.Create(&req)
	if err != nil {
		response.ServerError(c, "create failed: "+err.Error())
		return
	}

	response.Success(c, item)
}

// Get returns example details
// @Summary  Get example by ID
// @Tags     Example
// @Security Bearer
// @Param    id path int true "ID"
// @Success  200 {object} response.Response{data=dto.ExampleResponse}
// @Router   /examples/{id} [get]
func (h *ExampleHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid ID")
		return
	}

	item, err := h.svc.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "record not found")
		return
	}

	response.Success(c, item)
}

// Update updates an example
// @Summary  Update example
// @Tags     Example
// @Security Bearer
// @Accept   json
// @Param    id   path int                       true "ID"
// @Param    body body dto.UpdateExampleRequest   true "update parameters"
// @Success  200  {object} response.Response{data=dto.ExampleResponse}
// @Router   /examples/{id} [put]
func (h *ExampleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid ID")
		return
	}

	var req dto.UpdateExampleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "invalid parameters: "+err.Error())
		return
	}

	item, err := h.svc.Update(uint(id), &req)
	if err != nil {
		response.ServerError(c, "update failed: "+err.Error())
		return
	}

	response.Success(c, item)
}

// Delete deletes an example
// @Summary  Delete example
// @Tags     Example
// @Security Bearer
// @Param    id path int true "ID"
// @Success  200 {object} response.Response
// @Router   /examples/{id} [delete]
func (h *ExampleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.ParamError(c, "invalid ID")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		response.ServerError(c, "delete failed: "+err.Error())
		return
	}

	response.OK(c)
}
