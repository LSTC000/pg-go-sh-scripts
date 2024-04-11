package v1

import (
	"pg-sh-scripts/internal/api"
	"pg-sh-scripts/internal/usecase"

	"github.com/gin-gonic/gin"
)

const (
	groupBashLogPath       = "/bash/log"
	getBashLogByBashIdPath = "/:bashId/list"
)

type BashLogHandler struct{}

func (h *BashLogHandler) Register(rg *gin.RouterGroup) {
	useCase := usecase.GetBashLogUseCase()
	group := rg.Group(groupBashLogPath)
	{
		group.GET(getBashLogByBashIdPath, useCase.GetBashLogListByBashId)
	}
}

func GetBashLogHandler() api.IHandler {
	return &BashLogHandler{}
}
