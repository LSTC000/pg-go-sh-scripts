package bash

import (
	"pg-sh-scripts/internal/api"

	"github.com/gin-gonic/gin"
)

const (
	groupPath           = "/bash"
	getBashByIdPath     = "/:id"
	getBashFileByIdPath = "/:id/file"
	getBashListPath     = "/list"
	createBashPath      = ""
	execBashPath        = "/execute"
	execBashListPath    = "/execute/list"
)

type Handler struct{}

func (h *Handler) Register(rg *gin.RouterGroup) {
	useCase := GetUseCase()
	group := rg.Group(groupPath)
	{
		group.GET(getBashByIdPath, useCase.GetBashById)
		group.GET(getBashFileByIdPath, useCase.GetBashFileById)
		group.GET(getBashListPath, useCase.GetBashList)
		group.POST(createBashPath, useCase.CreateBash)
		group.POST(execBashPath, useCase.ExecBash)
		group.POST(execBashListPath, useCase.ExecBashList)
	}
}

func GetHandler() api.IHandler {
	return &Handler{}
}
