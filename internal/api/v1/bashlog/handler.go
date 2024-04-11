package bashlog

import (
	"pg-sh-scripts/internal/api"

	"github.com/gin-gonic/gin"
)

const (
	groupPath              = "/bash/log"
	getBashLogByBashIDPath = "/:bashId/list"
)

type Handler struct{}

func (h *Handler) Register(rg *gin.RouterGroup) {
	useCase := GetUseCase()
	group := rg.Group(groupPath)
	{
		group.GET(getBashLogByBashIDPath, useCase.GetBashLogListByBashID)
	}
}

func GetHandler() api.IHandler {
	return &Handler{}
}
