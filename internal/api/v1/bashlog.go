package v1

import (
	"net/http"
	"pg-sh-scripts/internal/api"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/usecase"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	groupBashLogPath       = "/bash/log"
	getBashLogByBashIdPath = "/:bashId/list"
)

type (
	IBashLogHandler interface {
		GetBashLogUseCase(ctx *gin.Context)
	}

	BashLogHandler struct {
		useCase    usecase.IBashLogUseCase
		httpErrors *config.HTTPErrors
	}
)

// GetBashLogListByBashId
// @Summary Get list by bash id
// @Tags Bash Log
// @Description Get list of bash logs by bash id
// @Produce json
// @Success 200 {object} schema.SwagBashLogPaginationLimitOffsetPage
// @Failure 500 {object} schema.HTTPError
// @Param bashId path string true "ID of bash script"
// @Param limit query int true "Limit param of pagination"
// @Param offset query int true "Offset param of pagination"
// @Router /bash/log/{bashId}/list [get]
func (h *BashLogHandler) GetBashLogListByBashId(ctx *gin.Context) {
	bashId, err := uuid.FromString(ctx.Param("bashId"))
	if err != nil {
		api.RaiseError(ctx, h.httpErrors.Validate)
		return
	}

	bashLogList, err := h.useCase.GetBashLogListByBashId(bashId)
	if err != nil {
		api.RaiseError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, bashLogList)
}

func (h *BashLogHandler) Register(rg *gin.RouterGroup) {
	group := rg.Group(groupBashLogPath)
	{
		group.GET(getBashLogByBashIdPath, h.GetBashLogListByBashId)
	}
}

func GetBashLogHandler() api.IHandler {
	return &BashLogHandler{
		useCase:    usecase.GetBashLogUseCase(),
		httpErrors: config.GetHTTPErrors(),
	}
}
