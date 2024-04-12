package usecase

import (
	"context"
	"net/http"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/service"
	"pg-sh-scripts/pkg/logging"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type (
	IBashLogUseCase interface {
		GetBashLogListByBashId(ctx *gin.Context)
	}

	BashLogUseCase struct {
		service    service.IBashLogService
		logger     *logging.Logger
		httpErrors *config.HTTPErrors
	}
)

// GetBashLogListByBashId
// @Summary Get list by bash id
// @Tags Bash Log
// @Description Get list of bash logs by bash id
// @Produce json
// @Success 200 {array} model.BashLog
// @Failure 500 {object} schema.HTTPError
// @Param bashId path string true "ID of bash script"
// @Router /bash/log/{bashId}/list [get]
func (u *BashLogUseCase) GetBashLogListByBashId(ctx *gin.Context) {
	bashId, err := uuid.FromString(ctx.Param("bashId"))
	if err != nil {
		ctx.JSON(u.httpErrors.Validate.HTTPCode, u.httpErrors.Validate)
		return
	}

	bashService := service.GetBashService()
	_, err = bashService.GetOneById(context.Background(), bashId)
	if err != nil {
		ctx.JSON(u.httpErrors.BashGet.HTTPCode, u.httpErrors.BashGet)
		return
	}

	bashLogList, err := u.service.GetAllByBashId(context.Background(), bashId)
	if err != nil {
		ctx.JSON(u.httpErrors.BashLogGetListByBashId.HTTPCode, u.httpErrors.BashLogGetListByBashId)
		return
	}

	ctx.JSON(http.StatusOK, bashLogList)
}

func GetBashLogUseCase() IBashLogUseCase {
	return &BashLogUseCase{
		service:    service.GetBashLogService(),
		logger:     log.GetLogger(),
		httpErrors: config.GetHTTPErrors(),
	}
}
