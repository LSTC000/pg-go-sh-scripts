package usecase

import (
	"context"
	"pg-sh-scripts/internal/config"
	mock_service "pg-sh-scripts/internal/service/mock"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/sql/pagination"
	"testing"

	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"

	_ "github.com/stretchr/testify/assert"
)

func TestBashUseCase_GetBashPaginationPage(t *testing.T) {
	type expectedStruct struct {
		paginationPage alias.BashLimitOffsetPage
		err            error
	}

	httpErrors := config.GetHTTPErrors()

	testCases := []struct {
		name             string
		ctx              context.Context
		bashId           uuid.UUID
		paginationParams pagination.LimitOffsetParams
		mockBehavior     func(*mock_service.MockIBashService, context.Context, pagination.LimitOffsetParams)
		expected         expectedStruct
	}{
		{
			name:             "Success",
			ctx:              context.Background(),
			bashId:           uuid.NewV4(),
			paginationParams: pagination.LimitOffsetParams{},
			mockBehavior: func(m *mock_service.MockIBashService, ctx context.Context, paginationParams pagination.LimitOffsetParams) {
				m.EXPECT().GetPaginationPage(ctx, paginationParams).Return(alias.BashLimitOffsetPage{}, nil)
			},
			expected: expectedStruct{
				paginationPage: alias.BashLimitOffsetPage{},
				err:            nil,
			},
		},
		{
			name:             "Getting bash pagination page error",
			ctx:              context.Background(),
			bashId:           uuid.NewV4(),
			paginationParams: pagination.LimitOffsetParams{},
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
			testCase.mockBehavior(mockBashService, testCase.ctx, testCase.paginationParams)

			bashUseCase := BashUseCase{
				service:    mockBashService,
				httpErrors: httpErrors,
			}

			bashLogPaginationPage, err := bashUseCase.GetBashPaginationPage(testCase.paginationParams)

			assert.Equal(t, testCase.expected.paginationPage, bashLogPaginationPage)
			assert.Equal(t, testCase.expected.err, err)
		})
	}
}
