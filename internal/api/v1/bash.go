package v1

import (
	"fmt"
	"net/http"
	"pg-sh-scripts/internal/api"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/msg"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/internal/usecase"
	"pg-sh-scripts/pkg/sql/pagination"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	groupBashPath       = "/bash"
	getBashByIdPath     = "/:id"
	getBashFileByIdPath = "/:id/file"
	getBashListPath     = "/list"
	createBashPath      = ""
	execBashListPath    = "/execute/list"
	removeBashPath      = "/:id"
)

type (
	IBashHandler interface {
		GetBashById(c *gin.Context)
		GetBashFileById(c *gin.Context)
		GetBashList(c *gin.Context)
		CreateBash(c *gin.Context)
		ExecBash(c *gin.Context)
		RemoveBashById(c *gin.Context)
	}

	BashHandler struct {
		useCase    usecase.IBashUseCase
		httpErrors *config.HTTPErrors
	}
)

func (h *BashHandler) Register(rg *gin.RouterGroup) {
	group := rg.Group(groupBashPath)
	{
		group.GET(getBashByIdPath, h.GetBashById)
		group.GET(getBashFileByIdPath, h.GetBashFileById)
		group.GET(getBashListPath, h.GetBashList)
		group.POST(createBashPath, h.CreateBash)
		group.POST(execBashListPath, h.ExecBashList)
		group.DELETE(removeBashPath, h.RemoveBashById)
	}
}

// GetBashById
// @Summary Get by id
// @Tags Bash
// @Description Get bash script by id
// @Produce json
// @Success 200 {object} model.Bash
// @Failure 500 {object} schema.HTTPError
// @Param id path string true "ID of bash script"
// @Router /bash/{id} [get]
func (h *BashHandler) GetBashById(c *gin.Context) {
	bashId, err := uuid.FromString(c.Param("id"))
	if err != nil {
		api.RaiseError(c, h.httpErrors.BashId)
		return
	}

	bash, err := h.useCase.GetBashById(bashId)
	if err != nil {
		api.RaiseError(c, err)
		return
	}

	c.JSON(http.StatusOK, bash)
}

// GetBashFileById
// @Summary Get file by id
// @Tags Bash
// @Description Get bash script file by id
// @Produce x-www-form-urlencoded
// @Success 200 {file} binary
// @Failure 500 {object} schema.HTTPError
// @Param id path string true "ID of bash script"
// @Router /bash/{id}/file [get]
func (h *BashHandler) GetBashFileById(c *gin.Context) {
	bashId, err := uuid.FromString(c.Param("id"))
	if err != nil {
		api.RaiseError(c, h.httpErrors.BashId)
		return
	}

	bashFileBuffer, bashTitle, err := h.useCase.GetBashFileBufferById(bashId)
	if err != nil {
		api.RaiseError(c, err)
		return
	}

	extraHeaders := map[string]string{"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s.sh\"", bashTitle)}
	c.DataFromReader(http.StatusOK, int64(bashFileBuffer.Len()), "application/x-www-form-urlencoded", bashFileBuffer, extraHeaders)
}

// GetBashList
// @Summary Get list
// @Tags Bash
// @Description Get list of bash scripts
// @Produce json
// @Success 200 {object} schema.BashPaginationPage
// @Failure 500 {object} schema.HTTPError
// @Param limit query int true "Limit param of pagination" default(20)
// @Param offset query int true "Offset param of pagination" default(0)
// @Router /bash/list [get]
func (h *BashHandler) GetBashList(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		api.RaiseError(c, h.httpErrors.PaginationLimitParamMustBeInt)
		return
	}
	if limit < 0 {
		api.RaiseError(c, h.httpErrors.PaginationLimitParamGTEZero)
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		api.RaiseError(c, h.httpErrors.PaginationOffsetParamMustBeInt)
		return
	}
	if offset < 0 {
		api.RaiseError(c, h.httpErrors.PaginationOffsetParamGTEZero)
		return
	}

	paginationParams := pagination.LimitOffsetParams{
		Limit:  limit,
		Offset: offset,
	}

	bashList, err := h.useCase.GetBashPaginationPage(paginationParams)
	if err != nil {
		api.RaiseError(c, err)
		return
	}

	c.JSON(http.StatusOK, bashList)
}

// CreateBash
// @Summary Create
// @Tags Bash
// @Description Create bash script
// @Accept mpfd
// @Produce json
// @Success 200 {object} model.Bash
// @Failure 500 {object} schema.HTTPError
// @Param file formData file true "Bash script file"
// @Router /bash [post]
func (h *BashHandler) CreateBash(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		api.RaiseError(c, h.httpErrors.BashFileUpload)
		return
	}

	bash, err := h.useCase.CreateBash(file)
	if err != nil {
		api.RaiseError(c, err)
		return
	}

	c.JSON(http.StatusOK, bash)
}

// ExecBashList
// @Summary Execute List
// @Tags Bash
// @Description Execute list of bash scripts
// @Accept json
// @Produce json
// @Success 200 {object} schema.Message
// @Failure 500 {object} schema.HTTPError
// @Param isSync query bool true "Execute type: if true, then in a multithreading, otherwise in a single thread"
// @Param execute body []dto.ExecBash true "List of execute bash script models"
// @Router /bash/execute/list [post]
func (h *BashHandler) ExecBashList(c *gin.Context) {
	isSync := c.GetBool("isSync")

	execBashDTOList := make([]dto.ExecBash, 0)

	if err := c.ShouldBindJSON(&execBashDTOList); err != nil {
		api.RaiseError(c, h.httpErrors.BashExecuteDTOList)
		return
	}

	if err := h.useCase.ExecBashList(isSync, execBashDTOList); err != nil {
		api.RaiseError(c, err)
		return
	}

	c.JSON(http.StatusOK, schema.Message{Message: msg.OK})
}

// RemoveBashById
// @Summary Remove by id
// @Tags Bash
// @Description Remove bash script by id
// @Produce json
// @Success 200 {object} model.Bash
// @Failure 500 {object} schema.HTTPError
// @Param id path string true "ID of bash script"
// @Router /bash/{id} [delete]
func (h *BashHandler) RemoveBashById(c *gin.Context) {
	bashId, err := uuid.FromString(c.Param("id"))
	if err != nil {
		api.RaiseError(c, h.httpErrors.BashId)
		return
	}

	bash, err := h.useCase.RemoveBashById(bashId)
	if err != nil {
		api.RaiseError(c, err)
		return
	}

	c.JSON(http.StatusOK, bash)
}

func GetBashHandler() api.IHandler {
	return &BashHandler{
		useCase:    usecase.GeBashUseCase(),
		httpErrors: config.GetHTTPErrors(),
	}
}
