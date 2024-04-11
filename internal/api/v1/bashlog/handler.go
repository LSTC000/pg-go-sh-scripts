package bashlog

import (
	"pg-sh-scripts/internal/api"

	"github.com/gin-gonic/gin"
)

const (
	groupPath              = "/bash/log"
	getBashLogByBashIdPath = "/:bashId/list"
)

type Handler struct{}

func (h *Handler) Register(rg *gin.RouterGroup) {
	useCase := GetUseCase()
	group := rg.Group(groupPath)
	{
		group.GET(getBashLogByBashIdPath, useCase.GetBashLogListByBashId)
	}
}

func GetHandler() api.IHandler {
	return &Handler{}
}
