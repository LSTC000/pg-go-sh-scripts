package common

import (
	"bufio"
	"context"
	"errors"
	"io"
	"pg-sh-scripts/internal/dto"
	"pg-sh-scripts/internal/log"
	"pg-sh-scripts/internal/service"
	"pg-sh-scripts/pkg/gosha"
	"pg-sh-scripts/pkg/logging"

	uuid "github.com/satori/go.uuid"
)

type (
	ICustomGoshaExec interface {
		Run()
	}

	CustomGoshaExec struct {
		isSync    bool
		goshaExec gosha.IExec
		logger    *logging.Logger
	}

	CustomScanner struct{}
)

func (c *CustomGoshaExec) saveExecError(err error) {
	var execErr *gosha.ExecErr

	bashLogService := service.GetBashLogService()

	if errors.As(err, &execErr) {
		bashId, err := uuid.FromString(execErr.Title)
		if err == nil {
			createBashLogDTO := dto.CreateBashLog{
				BashId:  bashId,
				Body:    execErr.Detail,
				IsError: true,
			}
			_, _ = bashLogService.Create(context.Background(), createBashLogDTO)
		}
	} else {
		c.logger.Error("Unknown execute error: &v", err)
	}
}

func (c *CustomGoshaExec) Run() {
	if c.isSync {
		if errs := c.goshaExec.SyncRun(); errs != nil {
			for _, err := range errs {
				c.saveExecError(err)
			}
		}
	} else {
		if err := c.goshaExec.Run(); err != nil {
			c.saveExecError(err)
		}
	}
}

func (s *CustomScanner) Scan(stdout io.ReadCloser, cmd *gosha.Cmd) error {
	scanner := bufio.NewScanner(stdout)
	bashLogService := service.GetBashLogService()

	bashId, err := uuid.FromString(cmd.Title)
	if err != nil {
		return err
	}

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		msg := scanner.Text()
		createBashLogDTO := dto.CreateBashLog{
			BashId:  bashId,
			Body:    msg,
			IsError: false,
		}
		if _, err := bashLogService.Create(context.Background(), createBashLogDTO); err != nil {
			return err
		}
	}
	return nil
}

func GetCustomGoshaExec(isSync bool, commands []gosha.ICmd) ICustomGoshaExec {
	return &CustomGoshaExec{
		isSync:    isSync,
		goshaExec: gosha.GetExec(&CustomScanner{}, commands),
		logger:    log.GetLogger(),
	}
}
