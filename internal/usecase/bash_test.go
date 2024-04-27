package usecase

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/model"
	mock_service "pg-sh-scripts/internal/service/mock"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/sql/pagination"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	_ "github.com/stretchr/testify/assert"
)

func TestBashUseCase_GetBashById(t *testing.T) {
	type (
		inStruct struct {
			ctx    context.Context
			bashId uuid.UUID
		}

		expectedStruct struct {
			bash *model.Bash
			err  error
		}
	)

	httpErrors := config.GetHTTPErrors()

	testCases := []struct {
		name         string
		in           inStruct
		mockBehavior func(*mock_service.MockIBashService, context.Context, uuid.UUID)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				ctx:    context.Background(),
				bashId: uuid.NewV4(),
			},
			mockBehavior: func(m *mock_service.MockIBashService, ctx context.Context, bashId uuid.UUID) {
				m.EXPECT().GetOneById(ctx, bashId).Return(&model.Bash{}, nil)
			},
			expected: expectedStruct{
				bash: &model.Bash{},
				err:  nil,
			},
		},
		{
			name: "Getting bash does not exists error",
			in: inStruct{
				ctx:    context.Background(),
				bashId: uuid.NewV4(),
			},
			mockBehavior: func(m *mock_service.MockIBashService, ctx context.Context, bashId uuid.UUID) {
				m.EXPECT().GetOneById(ctx, bashId).Return(nil, httpErrors.BashDoesNotExists)
			},
			expected: expectedStruct{
				bash: nil,
				err:  httpErrors.BashDoesNotExists,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBashService := mock_service.NewMockIBashService(ctrl)
			testCase.mockBehavior(mockBashService, testCase.in.ctx, testCase.in.bashId)

			bashUseCase := BashUseCase{
				service:    mockBashService,
				httpErrors: httpErrors,
			}

			bash, err := bashUseCase.GetBashById(testCase.in.bashId)

			assert.Equal(t, testCase.expected.bash, bash)
			assert.Equal(t, testCase.expected.err, err)
		})
	}
}

func TestBashUseCase_GetBashPaginationPage(t *testing.T) {
	type (
		inStruct struct {
			ctx              context.Context
			paginationParams pagination.LimitOffsetParams
		}

		expectedStruct struct {
			paginationPage alias.BashLimitOffsetPage
			err            error
		}
	)

	httpErrors := config.GetHTTPErrors()

	testCases := []struct {
		name         string
		in           inStruct
		mockBehavior func(*mock_service.MockIBashService, context.Context, pagination.LimitOffsetParams)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				ctx:              context.Background(),
				paginationParams: pagination.LimitOffsetParams{},
			},
			mockBehavior: func(m *mock_service.MockIBashService, ctx context.Context, paginationParams pagination.LimitOffsetParams) {
				m.EXPECT().GetPaginationPage(ctx, paginationParams).Return(alias.BashLimitOffsetPage{}, nil)
			},
			expected: expectedStruct{
				paginationPage: alias.BashLimitOffsetPage{},
				err:            nil,
			},
		},
		{
			name: "Getting bash pagination page error",
			in: inStruct{
				ctx:              context.Background(),
				paginationParams: pagination.LimitOffsetParams{},
			},
			mockBehavior: func(m *mock_service.MockIBashService, ctx context.Context, paginationParams pagination.LimitOffsetParams) {
				m.EXPECT().GetPaginationPage(ctx, paginationParams).Return(alias.BashLimitOffsetPage{}, httpErrors.BashGetPaginationPage)
			},
			expected: expectedStruct{
				paginationPage: alias.BashLimitOffsetPage{},
				err:            httpErrors.BashGetPaginationPage,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBashService := mock_service.NewMockIBashService(ctrl)
			testCase.mockBehavior(mockBashService, testCase.in.ctx, testCase.in.paginationParams)

			bashUseCase := BashUseCase{
				service:    mockBashService,
				httpErrors: httpErrors,
			}

			bashLogPaginationPage, err := bashUseCase.GetBashPaginationPage(testCase.in.paginationParams)

			assert.Equal(t, testCase.expected.paginationPage, bashLogPaginationPage)
			assert.Equal(t, testCase.expected.err, err)
		})
	}
}
