package usecase

import (
	"bytes"
	"context"
	"mime/multipart"
	"os"
	"pg-sh-scripts/internal/common"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/service"
	"pg-sh-scripts/internal/types/alias"
	"pg-sh-scripts/internal/util"
	"pg-sh-scripts/pkg/gosha"
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	IBashUseCase interface {
		GetBashById(bashId uuid.UUID) (*model.Bash, error)
		GetBashFileBufferById(bashId uuid.UUID) (*bytes.Buffer, alias.BashTitle, error)
		GetBashList() ([]*model.Bash, error)
		CreateBash(file *multipart.FileHeader) (*model.Bash, error)
		ExecBash(isSync bool, dto dto.ExecBashDTO) error
		ExecBashList(isSync bool, dto []dto.ExecBashDTO) error
		RemoveBashById(bashId uuid.UUID) (*model.Bash, error)
	}

	BashUseCase struct {
		service    service.IBashService
		httpErrors *config.HTTPErrors
	}
)

func (u *BashUseCase) GetBashById(bashId uuid.UUID) (*model.Bash, error) {
	bash, err := u.service.GetOneById(context.Background(), bashId)
	if err != nil {
		return nil, u.httpErrors.BashGet
	}
	return bash, nil
}

func (u *BashUseCase) GetBashFileBufferById(bashId uuid.UUID) (*bytes.Buffer, alias.BashTitle, error) {
	bash, err := u.service.GetOneById(context.Background(), bashId)
	if err != nil {
		return nil, "", u.httpErrors.BashGet
	}

	bashFileBuffer, err := util.GetBashFileBuffer(bash.Title, bash.Body)
	if err != nil {
		return nil, "", u.httpErrors.BashGetFile
	}

	return bashFileBuffer, bash.Title, nil
}

func (u *BashUseCase) GetBashList() ([]*model.Bash, error) {
	bashList, err := u.service.GetAll(context.Background())
	if err != nil {
		return nil, u.httpErrors.BashGetList
	}
	return bashList, nil
}

func (u *BashUseCase) CreateBash(file *multipart.FileHeader) (*model.Bash, error) {
	fileName := file.Filename
	fileExtension := util.GetBashFileExtension(fileName)

	if err := util.ValidateBashFileExtension(fileExtension); err != nil {
		return nil, u.httpErrors.BashFileExtension
	}

	fileTitle := util.GetBashFileTitle(fileName)
	fileBody, err := util.GetBashFileBody(file)
	if err != nil {
		return nil, u.httpErrors.BashFileBody
	}

	createBashDTO := dto.CreateBashDTO{Title: fileTitle, Body: fileBody}
	bash, err := u.service.Create(context.Background(), createBashDTO)
	if err != nil {
		return nil, u.httpErrors.BashCreate
	}

	return bash, nil
}

func (u *BashUseCase) ExecBash(isSync bool, dto dto.ExecBashDTO) error {
	bash, err := u.service.GetOneById(context.Background(), dto.Id)
	if err != nil {
		return u.httpErrors.BashGet
	}

	tmpFile, err := gosha.GetTmpFile(bash.Body)
	if err != nil {
		return u.httpErrors.BashExecute
	}
	defer func() {
		_ = gosha.RemoveTmpFile(tmpFile)
	}()

	commands := []gosha.Cmd{
		{
			Title:   bash.Id.String(),
			Path:    tmpFile.Name(),
			Timeout: dto.TimeoutSeconds * time.Second,
		},
	}

	customGoshaExec := common.GetCustomGoshaExec(isSync, commands)
	customGoshaExec.Run()

	return nil
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
func (u *BashUseCase) ExecBashList(isSync bool, dto []dto.ExecBashDTO) error {
	execBashCount := len(dto)
	bashList := make([]*model.Bash, 0, execBashCount)
	for _, execBashDTO := range dto {
		bash, err := u.service.GetOneById(context.Background(), execBashDTO.Id)
		if err != nil {
			return u.httpErrors.BashGet
		}
		bashList = append(bashList, bash)
	}

	tmpFiles := make([]*os.File, 0, execBashCount)
	commands := make([]gosha.Cmd, 0, execBashCount)
	for i := 0; i < execBashCount; i++ {
		bash := bashList[i]
		execBashDTO := dto[i]

		tmpFile, err := gosha.GetTmpFile(bash.Body)
		if err != nil {
			return u.httpErrors.BashExecute
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

	customGoshaExec := common.GetCustomGoshaExec(isSync, commands)
	customGoshaExec.Run()

	return nil
}

func (u *BashUseCase) RemoveBashById(bashId uuid.UUID) (*model.Bash, error) {
	_, err := u.service.GetOneById(context.Background(), bashId)
	if err != nil {
		return nil, u.httpErrors.BashGet
	}

	bash, err := u.service.RemoveById(context.Background(), bashId)
	if err != nil {
		return nil, u.httpErrors.BashRemove
	}

	return bash, nil
}

func GeBashUseCase() IBashUseCase {
	return &BashUseCase{
		service:    service.GetBashService(),
		httpErrors: config.GetHTTPErrors(),
	}
}
