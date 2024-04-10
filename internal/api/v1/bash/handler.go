package bash

import (
	"pg-sh-scripts/internal/api"

	"github.com/gin-gonic/gin"
)

const (
	groupPath           = "/bash"
	getBashByIDPath     = "/:id"
	getBashFileByIDPath = "/:id/file"
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
		group.GET(getBashByIDPath, useCase.GetBashByID)
		group.GET(getBashFileByIDPath, useCase.GetBashFileByID)
		group.GET(getBashListPath, useCase.GetBashList)
		group.POST(createBashPath, useCase.CreateBash)
		group.POST(execBashPath, useCase.ExecBash)
		group.POST(execBashListPath, useCase.ExecBashList)
	}
}

func GetHandler() api.IHandler {
	return &Handler{}
}
