package v1

import (
	"net/http"
	"pg-sh-scripts/internal/api"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/usecase"
	"pg-sh-scripts/pkg/sql/pagination"
	"strconv"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	groupBashLogPath       = "/bash/log"
	getBashLogByBashIdPath = "/:bashId/list"
)

type (
	IBashLogHandler interface {
		GetBashLogUseCase(c *gin.Context)
	}

	BashLogHandler struct {
		useCase    usecase.IBashLogUseCase
		httpErrors *config.HTTPErrors
	}
)

func (h *BashLogHandler) Register(rg *gin.RouterGroup) {
	group := rg.Group(groupBashLogPath)
	{
		group.GET(getBashLogByBashIdPath, h.GetBashLogListByBashId)
	}
}

// GetBashLogListByBashId
// @Summary Get list by bash id
// @Tags Bash Log
// @Description Get list of bash logs by bash id
// @Produce json
// @Success 200 {object} schema.SwagBashLogPaginationPage
// @Failure 500 {object} schema.HTTPError
// @Param bashId path string true "ID of bash script"
// @Param limit query int true "Limit param of pagination" default(20)
// @Param offset query int true "Offset param of pagination" default(0)
// @Router /bash/log/{bashId}/list [get]
func (h *BashLogHandler) GetBashLogListByBashId(c *gin.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		api.RaiseError(c, h.httpErrors.Validate)
		return
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		api.RaiseError(c, h.httpErrors.Validate)
		return
	}

	bashId, err := uuid.FromString(c.Param("bashId"))
	if err != nil {
		api.RaiseError(c, h.httpErrors.Validate)
		return
	}

	paginationParams := pagination.LimitOffsetParams{
		Limit:  limit,
		Offset: offset,
	}

	bashLogList, err := h.useCase.GetBashLogPaginationPageByBashId(bashId, paginationParams)
	if err != nil {
		api.RaiseError(c, err)
		return
	}

	c.JSON(http.StatusOK, bashLogList)
}

func GetBashLogHandler() api.IHandler {
	return &BashLogHandler{
		useCase:    usecase.GetBashLogUseCase(),
		httpErrors: config.GetHTTPErrors(),
	}
}
