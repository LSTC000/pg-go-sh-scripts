package bashlog

import (
	"context"
	"net/http"
	"pg-sh-scripts/internal/api/v1/bash"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/pkg/logging"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type (
	IUseCase interface {
		GetBashLogListByBashID(ctx *gin.Context)
	}

	UseCase struct {
		service IService
		logger  *logging.Logger
	}
)

// GetBashLogListByBashID
// @Summary Get list by bash id
// @Tags Bash Log
// @Description Get list of bash logs by bash id
// @Produce json
// @Success 200 {array} BashLog
// @Failure 500 {object} model.HTTPError
// @Router /bash/log/{bashId}/list [get]
func (u *UseCase) GetBashLogListByBashID(ctx *gin.Context) {
	httpErrors := config.GetHTTPErrors()

	bashID, err := uuid.FromString(ctx.Param("bashId"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, httpErrors.Validate)
		return
	}

	bashService := bash.GetService()
	_, err = bashService.GetOneByID(context.Background(), bashID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashGet)
		return
	}

	bashLogList, err := u.service.GetAllByBashID(context.Background(), bashID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashLogGetListByBashID)
		return
	}

	ctx.JSON(http.StatusOK, bashLogList)
}

func GetUseCase() IUseCase {
	return &UseCase{
		service: GetService(),
		logger:  common.GetLogger(),
	}
}
