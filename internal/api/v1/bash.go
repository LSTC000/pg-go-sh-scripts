package v1

import (
	"pg-sh-scripts/internal/api"
	"pg-sh-scripts/internal/usecase"

	"github.com/gin-gonic/gin"
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

type BashHandler struct{}

func (h *BashHandler) Register(rg *gin.RouterGroup) {
	useCase := usecase.GeBashUseCase()
	group := rg.Group(groupBashPath)
	{
		group.GET(getBashByIdPath, useCase.GetBashById)
		group.GET(getBashFileByIdPath, useCase.GetBashFileById)
		group.GET(getBashListPath, useCase.GetBashList)
		group.POST(createBashPath, useCase.CreateBash)
		group.POST(execBashPath, useCase.ExecBash)
		group.POST(execBashListPath, useCase.ExecBashList)
	}
}

func GetBashHandler() api.IHandler {
	return &BashHandler{}
}
