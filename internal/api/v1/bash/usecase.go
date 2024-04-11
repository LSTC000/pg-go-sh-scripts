package bash

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/pkg/gosha"
	"pg-sh-scripts/pkg/logging"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type (
	IUseCase interface {
		GetBashById(ctx *gin.Context)
		GetBashFileById(ctx *gin.Context)
		GetBashList(ctx *gin.Context)
		CreateBash(ctx *gin.Context)
		ExecBash(ctx *gin.Context)
		ExecBashList(ctx *gin.Context)
	}

	UseCase struct {
		service IService
		logger  *logging.Logger
	}
)

// GetBashById
// @Summary Get by id
// @Tags Bash
// @Description Get bash script by id
// @Produce json
// @Success 200 {object} Bash
// @Failure 500 {object} model.HTTPError
// @Param id path string true "ID of bash script"
// @Router /bash/{id} [get]
func (u *UseCase) GetBashById(ctx *gin.Context) {
	httpErrors := config.GetHTTPErrors()

	bashId, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, httpErrors.Validate)
		return
	}

	bash, err := u.service.GetOneById(context.Background(), bashId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashGet)
		return
	}

	ctx.JSON(http.StatusOK, bash)
}

// GetBashFileById
// @Summary Get file by id
// @Tags Bash
// @Description Get bash script file by id
// @Produce x-www-form-urlencoded
// @Success 200 {file} binary
// @Failure 500 {object} model.HTTPError
// @Param id path string true "ID of bash script"
// @Router /bash/{id}/file [get]
func (u *UseCase) GetBashFileById(ctx *gin.Context) {
	httpErrors := config.GetHTTPErrors()

	bashId, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, httpErrors.Validate)
		return
	}

	bash, err := u.service.GetOneById(context.Background(), bashId)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashGet)
		return
	}

	fileBuffer, err := GetBashFileBuffer(bash.Title, bash.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashGetFile)
		return
	}

	extraHeaders := map[string]string{"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s.sh\"", bash.Title)}
	ctx.DataFromReader(http.StatusOK, int64(fileBuffer.Len()), "application/x-www-form-urlencoded", fileBuffer, extraHeaders)
}

// GetBashList
// @Summary Get list
// @Tags Bash
// @Description Get list of bash scripts
// @Produce json
// @Success 200 {array} Bash
// @Failure 500 {object} model.HTTPError
// @Router /bash/list [get]
func (u *UseCase) GetBashList(ctx *gin.Context) {
	httpErrors := config.GetHTTPErrors()

	bashList, err := u.service.GetAll(context.Background())
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashGetList)
		return
	}

	ctx.JSON(http.StatusOK, bashList)
}

// CreateBash
// @Summary Create
// @Tags Bash
// @Description Create bash script
// @Accept mpfd
// @Produce json
// @Success 200 {object} Bash
// @Failure 500 {object} model.HTTPError
// @Param file formData file true "Bash script file"
// @Router /bash [post]
func (u *UseCase) CreateBash(ctx *gin.Context) {
	httpErrors := config.GetHTTPErrors()

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, httpErrors.Validate)
		return
	}

	fileName := file.Filename
	fileExtension := GetBashFileExtension(fileName)

	if err := ValidateBashFileExtension(fileExtension); err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashFileExtension)
		return
	}

	fileTitle := GetBashFileTitle(fileName)
	fileBody, err := GetBashFileBody(file)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashFileBody)
		return
	}

	createBash := CreateBashDTO{Title: fileTitle, Body: fileBody}
	bash, err := u.service.CreateBash(context.Background(), createBash)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashCreate)
		return
	}

	ctx.JSON(http.StatusOK, bash)
}

// ExecBash
// @Summary Execute
// @Tags Bash
// @Description Execute bash script
// @Accept json
// @Produce json
// @Success 200 {object} model.Message
// @Failure 500 {object} model.HTTPError
// @Param isSync query bool true "Execute type: if true, then in a multithreading, otherwise in a single thread"
// @Param execute body ExecBashDTO true "Execute bash script model"
// @Router /bash/execute [post]
func (u *UseCase) ExecBash(ctx *gin.Context) {
	execBash := ExecBashDTO{}
	httpErrors := config.GetHTTPErrors()

	isSync := ctx.GetBool("isSync")
	if err := ctx.ShouldBindJSON(&execBash); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, httpErrors.Validate)
		return
	}

	bash, err := u.service.GetOneById(context.Background(), execBash.Id)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashGet)
		return
	}

	tmpFile, err := gosha.GetTmpFile(bash.Body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpErrors.BashExecute)
		return
	}
	defer func() {
		if err := gosha.RemoveTmpFile(tmpFile); err != nil {
			ctx.JSON(http.StatusBadRequest, httpErrors.BashExecute)
		}
	}()

	commands := []gosha.Cmd{
		{
			Title:   bash.Id.String(),
			Path:    tmpFile.Name(),
			Timeout: execBash.TimeoutSeconds * time.Second,
		},
	}

	goshaExec := common.GetGoshaExec(commands)
	message := model.Message{Message: "ok"}

	if isSync {
		if errs := goshaExec.SyncRun(); errs != nil {
			message.Message = fmt.Sprintf("Execute Error: %v", errs)
		}
	} else {
		if err := goshaExec.Run(); err != nil {
			message.Message = fmt.Sprintf("Execute Error: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, message)
}

// ExecBashList
// @Summary Execute list
// @Tags Bash
// @Description Execute list of bash scripts
// @Accept json
// @Produce json
// @Success 200 {object} model.Message
// @Failure 500 {object} model.HTTPError
// @Param isSync query bool true "Execute type: if true, then in a multithreading, otherwise in a single thread"
// @Param execute body []ExecBashDTO true "List of execute bash script models"
// @Router /bash/execute/list [post]
func (u *UseCase) ExecBashList(ctx *gin.Context) {
	execBashList := make([]ExecBashDTO, 0)
	httpErrors := config.GetHTTPErrors()

	isSync := ctx.GetBool("isSync")
	if err := ctx.ShouldBindJSON(&execBashList); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, httpErrors.Validate)
		return
	}

	execBashCount := len(execBashList)
	bashList := make([]*Bash, 0, execBashCount)
	for _, execBash := range execBashList {
		bash, err := u.service.GetOneById(context.Background(), execBash.Id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, httpErrors.BashGet)
			return
		}
		bashList = append(bashList, bash)
	}

	tmpFiles := make([]*os.File, 0, execBashCount)
	commands := make([]gosha.Cmd, 0, execBashCount)
	for i := 0; i < execBashCount; i++ {
		bash := bashList[i]
		execBash := execBashList[i]

		tmpFile, err := gosha.GetTmpFile(bash.Body)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, httpErrors.BashExecute)
			return
		}
		tmpFiles = append(tmpFiles, tmpFile)

		cmd := gosha.Cmd{
			Title:   bash.Id.String(),
			Path:    tmpFile.Name(),
			Timeout: execBash.TimeoutSeconds * time.Second,
		}
		commands = append(commands, cmd)
	}
	defer func() {
		for _, tmpFile := range tmpFiles {
			if err := gosha.RemoveTmpFile(tmpFile); err != nil {
				ctx.JSON(http.StatusBadRequest, httpErrors.BashExecute)
			}
		}
	}()

	goshaExec := common.GetGoshaExec(commands)
	message := model.Message{Message: "ok"}

	if isSync {
		if errs := goshaExec.SyncRun(); errs != nil {
			message.Message = fmt.Sprintf("Execute Error: %v", errs)
		}
	} else {
		if err := goshaExec.Run(); err != nil {
			message.Message = fmt.Sprintf("Execute Error: %v", err)
		}
	}

	ctx.JSON(http.StatusOK, message)
}

func GetUseCase() IUseCase {
	return &UseCase{
		service: GetService(),
		logger:  common.GetLogger(),
	}
}
