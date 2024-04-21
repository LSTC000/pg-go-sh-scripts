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
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/internal/util"
	"pg-sh-scripts/pkg/gosha"
	"pg-sh-scripts/pkg/sql/pagination"
	"time"

	uuid "github.com/satori/go.uuid"
)

type (
	IBashUseCase interface {
		GetBashById(bashId uuid.UUID) (*model.Bash, error)
		GetBashFileBufferById(bashId uuid.UUID) (*bytes.Buffer, alias.BashTitle, error)
		GetBashPaginationPage(paginationParams pagination.LimitOffsetParams) (alias.BashLimitOffsetPage, error)
		CreateBash(file *multipart.FileHeader) (*model.Bash, error)
		ExecBashList(isSync bool, dto []dto.ExecBash) error
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
		return nil, u.httpErrors.BashDoesNotExists
	}
	return bash, nil
}

func (u *BashUseCase) GetBashFileBufferById(bashId uuid.UUID) (*bytes.Buffer, alias.BashTitle, error) {
	bash, err := u.service.GetOneById(context.Background(), bashId)
	if err != nil {
		return nil, "", u.httpErrors.BashDoesNotExists
	}

	bashFileBuffer, err := util.GetBashFileBuffer(bash.Title, bash.Body)
	if err != nil {
		return nil, "", u.httpErrors.BashGetFileBuffer
	}

	return bashFileBuffer, bash.Title, nil
}

func (u *BashUseCase) GetBashPaginationPage(paginationParams pagination.LimitOffsetParams) (alias.BashLimitOffsetPage, error) {
	bashPaginationPage, err := u.service.GetPaginationPage(context.Background(), paginationParams)
	if err != nil {
		return bashPaginationPage, u.httpErrors.BashGetPaginationPage
	}
	return bashPaginationPage, nil
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

	createBashDTO := dto.CreateBash{Title: fileTitle, Body: fileBody}
	bash, err := u.service.Create(context.Background(), createBashDTO)
	if err != nil {
		return nil, u.httpErrors.BashCreate
	}

	return bash, nil
}

func (u *BashUseCase) ExecBashList(isSync bool, dto []dto.ExecBash) error {
	execBashCount := len(dto)
	bashList := make([]*model.Bash, 0, execBashCount)

	for _, execBashDTO := range dto {
		bash, err := u.service.GetOneById(context.Background(), execBashDTO.Id)
		if err != nil {
			return u.httpErrors.BashDoesNotExists
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
		return nil, u.httpErrors.BashDoesNotExists
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
