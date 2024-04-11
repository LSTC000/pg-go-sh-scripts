package usecase

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/internal/service"
	"pg-sh-scripts/internal/util"
	"pg-sh-scripts/pkg/gosha"
	"pg-sh-scripts/pkg/logging"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
)

type (
	IBashUseCase interface {
		GetBashById(ctx *gin.Context)
		GetBashFileById(ctx *gin.Context)
		GetBashList(ctx *gin.Context)
		CreateBash(ctx *gin.Context)
		ExecBash(ctx *gin.Context)
		ExecBashList(ctx *gin.Context)
	}

	BashUseCase struct {
		service    service.IBashService
		logger     *logging.Logger
		httpErrors *config.HTTPErrors
	}

	PgScanner struct{}
)

func (s *PgScanner) Scan(stdout io.ReadCloser, cmd gosha.Cmd) error {
	scanner := bufio.NewScanner(stdout)
	bashLogService := service.GetBashLogService()

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		msg := scanner.Text()
		bashId, err := uuid.FromString(cmd.Title)
		if err != nil {
			return err
		}

		createBashLogDTO := dto.CreateBashLogDTO{
			BashId: bashId,
			Body:   msg,
		}
		if _, err := bashLogService.Create(context.Background(), createBashLogDTO); err != nil {
			return err
		}
	}
	return nil
}

// GetBashById
// @Summary Get by id
// @Tags Bash
// @Description Get bash script by id
// @Produce json
// @Success 200 {object} model.Bash
// @Failure 500 {object} schema.HTTPError
// @Param id path string true "ID of bash script"
// @Router /bash/{id} [get]
func (u *BashUseCase) GetBashById(ctx *gin.Context) {
	bashId, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		ctx.JSON(u.httpErrors.Validate.HTTPCode, u.httpErrors.Validate)
		return
	}

	bash, err := u.service.GetOneById(context.Background(), bashId)
	if err != nil {
		ctx.JSON(u.httpErrors.BashGet.HTTPCode, u.httpErrors.BashGet)
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
// @Failure 500 {object} schema.HTTPError
// @Param id path string true "ID of bash script"
// @Router /bash/{id}/file [get]
func (u *BashUseCase) GetBashFileById(ctx *gin.Context) {
	bashId, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		ctx.JSON(u.httpErrors.Validate.HTTPCode, u.httpErrors.Validate)
		return
	}

	bash, err := u.service.GetOneById(context.Background(), bashId)
	if err != nil {
		ctx.JSON(u.httpErrors.BashGet.HTTPCode, u.httpErrors.BashGet)
		return
	}

	fileBuffer, err := util.GetBashFileBuffer(bash.Title, bash.Body)
	if err != nil {
		ctx.JSON(u.httpErrors.BashGetFile.HTTPCode, u.httpErrors.BashGetFile)
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
// @Success 200 {array} model.Bash
// @Failure 500 {object} schema.HTTPError
// @Router /bash/list [get]
func (u *BashUseCase) GetBashList(ctx *gin.Context) {
	bashList, err := u.service.GetAll(context.Background())
	if err != nil {
		ctx.JSON(u.httpErrors.BashGetList.HTTPCode, u.httpErrors.BashGetList)
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
// @Success 200 {object} model.Bash
// @Failure 500 {object} schema.HTTPError
// @Param file formData file true "Bash script file"
// @Router /bash [post]
func (u *BashUseCase) CreateBash(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(u.httpErrors.Validate.HTTPCode, u.httpErrors.Validate)
		return
	}

	fileName := file.Filename
	fileExtension := util.GetBashFileExtension(fileName)

	if err := util.ValidateBashFileExtension(fileExtension); err != nil {
		ctx.JSON(u.httpErrors.BashFileExtension.HTTPCode, u.httpErrors.BashFileExtension)
		return
	}

	fileTitle := util.GetBashFileTitle(fileName)
	fileBody, err := util.GetBashFileBody(file)
	if err != nil {
		ctx.JSON(u.httpErrors.BashFileBody.HTTPCode, u.httpErrors.BashFileBody)
		return
	}

	createBashDTO := dto.CreateBashDTO{Title: fileTitle, Body: fileBody}
	bash, err := u.service.CreateBash(context.Background(), createBashDTO)
	if err != nil {
		ctx.JSON(u.httpErrors.BashCreate.HTTPCode, u.httpErrors.BashCreate)
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
// @Success 200 {object} schema.Message
// @Failure 500 {object} schema.HTTPError
// @Param isSync query bool true "Execute type: if true, then in a multithreading, otherwise in a single thread"
// @Param execute body dto.ExecBashDTO true "Execute bash script schema"
// @Router /bash/execute [post]
func (u *BashUseCase) ExecBash(ctx *gin.Context) {
	execBashDTO := dto.ExecBashDTO{}

	isSync := ctx.GetBool("isSync")
	if err := ctx.ShouldBindJSON(&execBashDTO); err != nil {
		ctx.JSON(u.httpErrors.Validate.HTTPCode, u.httpErrors.Validate)
		return
	}

	bash, err := u.service.GetOneById(context.Background(), execBashDTO.Id)
	if err != nil {
		ctx.JSON(u.httpErrors.BashGet.HTTPCode, u.httpErrors.BashGet)
		return
	}

	tmpFile, err := gosha.GetTmpFile(bash.Body)
	if err != nil {
		ctx.JSON(u.httpErrors.BashExecute.HTTPCode, u.httpErrors.BashExecute)
		return
	}
	defer func() {
		_ = gosha.RemoveTmpFile(tmpFile)
	}()

	commands := []gosha.Cmd{
		{
			Title:   bash.Id.String(),
			Path:    tmpFile.Name(),
			Timeout: execBashDTO.TimeoutSeconds * time.Second,
		},
	}

	goshaExec := common.GetGoshaExec(&PgScanner{}, commands)
	message := schema.Message{Message: "ok"}

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
// @Success 200 {object} schema.Message
// @Failure 500 {object} schema.HTTPError
// @Param isSync query bool true "Execute type: if true, then in a multithreading, otherwise in a single thread"
// @Param execute body []dto.ExecBashDTO true "List of execute bash script models"
// @Router /bash/execute/list [post]
func (u *BashUseCase) ExecBashList(ctx *gin.Context) {
	execBashDTOList := make([]dto.ExecBashDTO, 0)

	isSync := ctx.GetBool("isSync")
	if err := ctx.ShouldBindJSON(&execBashDTOList); err != nil {
		ctx.JSON(u.httpErrors.Validate.HTTPCode, u.httpErrors.Validate)
		return
	}

	execBashCount := len(execBashDTOList)
	bashList := make([]*model.Bash, 0, execBashCount)
	for _, execBashDTO := range execBashDTOList {
		bash, err := u.service.GetOneById(context.Background(), execBashDTO.Id)
		if err != nil {
			ctx.JSON(u.httpErrors.BashGet.HTTPCode, u.httpErrors.BashGet)
			return
		}
		bashList = append(bashList, bash)
	}

	tmpFiles := make([]*os.File, 0, execBashCount)
	commands := make([]gosha.Cmd, 0, execBashCount)
	for i := 0; i < execBashCount; i++ {
		bash := bashList[i]
		execBashDTO := execBashDTOList[i]

		tmpFile, err := gosha.GetTmpFile(bash.Body)
		if err != nil {
			ctx.JSON(u.httpErrors.BashExecute.HTTPCode, u.httpErrors.BashExecute)
			return
		}
		tmpFiles = append(tmpFiles, tmpFile)

		cmd := gosha.Cmd{
			Title:   bash.Id.String(),
			Path:    tmpFile.Name(),
			Timeout: execBashDTO.TimeoutSeconds * time.Second,
		}
		commands = append(commands, cmd)
	}
	defer func() {
		for _, tmpFile := range tmpFiles {
			_ = gosha.RemoveTmpFile(tmpFile)
		}
	}()

	goshaExec := common.GetGoshaExec(&PgScanner{}, commands)
	message := schema.Message{Message: "ok"}

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

func GeBashUseCase() IBashUseCase {
	return &BashUseCase{
		service:    service.GetBashService(),
		logger:     common.GetLogger(),
		httpErrors: config.GetHTTPErrors(),
	}
}
