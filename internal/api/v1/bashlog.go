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
	groupBashLogPath           = "/bash/log"
	getBashLogListByBashIdPath = "/:bashId/list"
)

type (
	IBashLogHandler interface {
		GetBashLogUseCase(c *gin.Context)
	}

	BashLogHandler struct {
		useCase    usecase.IBashLogUseCase
		helper     api.IHelper
		httpErrors *config.HTTPErrors
	}
)

func (h *BashLogHandler) Register(rg *gin.RouterGroup) {
	group := rg.Group(groupBashLogPath)
	{
		group.GET(getBashLogListByBashIdPath, h.GetBashLogListByBashId)
	}
}

// GetBashLogListByBashId
// @Summary Get list by bash id
// @Tags Bash Log
// @Description Get list of bash logs by bash id
// @Produce json
// @Success 200 {object} schema.BashLogPaginationPage
// @Failure 500 {object} schema.HTTPError
// @Param bashId path string true "ID of bash script"
// @Param limit query int true "Limit param of pagination" default(20)
// @Param offset query int true "Offset param of pagination" default(0)
// @Router /bash/log/{bashId}/list [get]
func (h *BashLogHandler) GetBashLogListByBashId(c *gin.Context) {
	bashId, err := uuid.FromString(c.Param("bashId"))
	if err != nil {
		httpError := h.helper.ParseError(c, h.httpErrors.BashId)
		c.JSON(httpError.HTTPCode, httpError)
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		httpError := h.helper.ParseError(c, h.httpErrors.PaginationLimitParamMustBeInt)
		c.JSON(httpError.HTTPCode, httpError)
		return
	}
	if limit < 0 {
		httpError := h.helper.ParseError(c, h.httpErrors.PaginationLimitParamGTEZero)
		c.JSON(httpError.HTTPCode, httpError)
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		httpError := h.helper.ParseError(c, h.httpErrors.PaginationOffsetParamMustBeInt)
		c.JSON(httpError.HTTPCode, httpError)
		return
	}
	if offset < 0 {
		httpError := h.helper.ParseError(c, h.httpErrors.PaginationOffsetParamGTEZero)
		c.JSON(httpError.HTTPCode, httpError)
		return
	}

	paginationParams := pagination.LimitOffsetParams{
		Limit:  limit,
		Offset: offset,
	}

	bashLogList, err := h.useCase.GetBashLogPaginationPageByBashId(bashId, paginationParams)
	if err != nil {
		httpError := h.helper.ParseError(c, err)
		c.JSON(httpError.HTTPCode, httpError)
		return
	}

	c.JSON(http.StatusOK, bashLogList)
}

func GetBashLogHandler() api.IHandler {
	return &BashLogHandler{
		useCase:    usecase.GetBashLogUseCase(),
		helper:     api.GetHelper(),
		httpErrors: config.GetHTTPErrors(),
	}
}
