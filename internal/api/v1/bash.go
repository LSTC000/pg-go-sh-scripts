package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"pg-sh-scripts/internal/api"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/internal/usecase"
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
		GetBashById(*gin.Context)
		GetBashFileById(*gin.Context)
		GetBashList(*gin.Context)
		CreateBash(*gin.Context)
		ExecBash(*gin.Context)
		RemoveBashById(*gin.Context)
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
func (h *BashHandler) GetBashById(ctx *gin.Context) {
	bashId, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		api.RaiseError(ctx, h.httpErrors.Validate)
		return
	}

	bash, err := h.useCase.GetBashById(bashId)
	if err != nil {
		api.RaiseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, bash)
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
func (h *BashHandler) GetBashFileById(ctx *gin.Context) {
	bashId, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		api.RaiseError(ctx, h.httpErrors.Validate)
		return
	}

	bashFileBuffer, bashTitle, err := h.useCase.GetBashFileBufferById(bashId)
	if err != nil {
		api.RaiseError(ctx, err)
		return
	}

	extraHeaders := map[string]string{"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s.sh\"", bashTitle)}
	ctx.DataFromReader(http.StatusOK, int64(bashFileBuffer.Len()), "application/x-www-form-urlencoded", bashFileBuffer, extraHeaders)
}

// GetBashList
// @Summary Get list
// @Tags Bash
// @Description Get list of bash scripts
// @Produce json
// @Success 200 {object} schema.SwagBashPaginationLimitOffsetPage
// @Failure 500 {object} schema.HTTPError
// @Param limit query int true "Limit param of pagination"
// @Param offset query int true "Offset param of pagination"
// @Router /bash/list [get]
func (h *BashHandler) GetBashList(ctx *gin.Context) {
	bashList, err := h.useCase.GetBashList()
	if err != nil {
		api.RaiseError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, bashList)
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
func (h *BashHandler) CreateBash(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		api.RaiseError(ctx, h.httpErrors.Validate)
		return
	}

	bash, err := h.useCase.CreateBash(file)
	if err != nil {
		api.RaiseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, bash)
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
// @Param execute body []dto.ExecBashDTO true "List of execute bash script models"
// @Router /bash/execute/list [post]
func (h *BashHandler) ExecBashList(ctx *gin.Context) {
	execBashDTOList := make([]dto.ExecBashDTO, 0)

	isSync := ctx.GetBool("isSync")
	if err := ctx.ShouldBindJSON(&execBashDTOList); err != nil {
		api.RaiseError(ctx, h.httpErrors.Validate)
		return
	}

	if err := h.useCase.ExecBashList(isSync, execBashDTOList); err != nil {
		api.RaiseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, schema.Message{Message: "ok"})
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
func (h *BashHandler) RemoveBashById(ctx *gin.Context) {
	bashId, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		api.RaiseError(ctx, h.httpErrors.Validate)
		return
	}

	bash, err := h.useCase.RemoveBashById(bashId)
	if err != nil {
		api.RaiseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, bash)
}

func GetBashHandler() api.IHandler {
	return &BashHandler{
		useCase:    usecase.GeBashUseCase(),
		httpErrors: config.GetHTTPErrors(),
	}
}
