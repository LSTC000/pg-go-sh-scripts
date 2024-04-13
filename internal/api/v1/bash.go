package v1

import (
	"fmt"
	"net/http"
	"pg-sh-scripts/internal/api"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/internal/usecase"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	groupBashPath       = "/bash"
	getBashByIdPath     = "/:id"
	getBashFileByIdPath = "/:id/file"
	getBashListPath     = "/list"
	createBashPath      = ""
	execBashPath        = "/execute"
	execBashListPath    = "/execute/list"
)

type (
	IBashHandler interface {
		GetBashById(*gin.Context)
		GetBashFileById(*gin.Context)
		GetBashList(*gin.Context)
		CreateBash(*gin.Context)
		ExecBash(*gin.Context)
		ExecBashList(*gin.Context)
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
		group.POST(execBashPath, h.ExecBash)
		group.POST(execBashListPath, h.ExecBashList)
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
// @Success 200 {array} model.Bash
// @Failure 500 {object} schema.HTTPError
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

// ExecBash
// @Summary Execute
// @Tags Bash
// @Description Execute bash script
// @Accept json
// @Produce json
// @Success 200 {object} schema.Message
// @Failure 500 {object} schema.HTTPError
// @Param isSync query bool true "Execute type: if true, then in a multithreading, otherwise in a single thread"
// @Param execute body dto.ExecBashDTO true "Execute bash script schema"
// @Router /bash/execute [post]
func (h *BashHandler) ExecBash(ctx *gin.Context) {
	execBashDTO := dto.ExecBashDTO{}

	isSync := ctx.GetBool("isSync")
	if err := ctx.ShouldBindJSON(&execBashDTO); err != nil {
		api.RaiseError(ctx, h.httpErrors.Validate)
		return
	}

	if err := h.useCase.ExecBash(isSync, execBashDTO); err != nil {
		api.RaiseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, schema.Message{Message: "ok"})
}

// ExecBashList
// @Summary Execute list
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

func GetBashHandler() api.IHandler {
	return &BashHandler{
		useCase:    usecase.GeBashUseCase(),
		httpErrors: config.GetHTTPErrors(),
	}
}
