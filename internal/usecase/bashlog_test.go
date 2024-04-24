package usecase

import (
	"context"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/model"
	mock_service "pg-sh-scripts/internal/service/mock"
	"pg-sh-scripts/internal/type/alias"
	"pg-sh-scripts/pkg/sql/pagination"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
)

func TestBashLogUseCase_GetBashLogPaginationPageByBashId(t *testing.T) {
	type (
		inStruct struct {
			ctx              context.Context
			bashId           uuid.UUID
			paginationParams pagination.LimitOffsetParams
		}

		expectedStruct struct {
			paginationPage alias.BashLogLimitOffsetPage
			err            error
		}
	)

	httpErrors := config.GetHTTPErrors()

	testCases := []struct {
		name         string
		in           inStruct
		mockBehavior func(*mock_service.MockIBashLogService, *mock_service.MockIBashService, context.Context, uuid.UUID, pagination.LimitOffsetParams)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				ctx:              context.Background(),
				bashId:           uuid.NewV4(),
				paginationParams: pagination.LimitOffsetParams{},
			},
			mockBehavior: func(mbl *mock_service.MockIBashLogService, mb *mock_service.MockIBashService, ctx context.Context, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) {
				gomock.InOrder(
					mb.EXPECT().GetOneById(ctx, bashId).Return(&model.Bash{}, nil),
					mbl.EXPECT().GetPaginationPageByBashId(ctx, bashId, paginationParams).Return(alias.BashLogLimitOffsetPage{}, nil),
				)
			},
			expected: expectedStruct{
				paginationPage: alias.BashLogLimitOffsetPage{},
				err:            nil,
			},
		},
		{
			name: "Bash does not exists",
			in: inStruct{
				ctx:              context.Background(),
				bashId:           uuid.NewV4(),
				paginationParams: pagination.LimitOffsetParams{},
			},
			mockBehavior: func(mbl *mock_service.MockIBashLogService, mb *mock_service.MockIBashService, ctx context.Context, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) {
				mb.EXPECT().GetOneById(ctx, bashId).Return(nil, httpErrors.BashDoesNotExists)
			},
			expected: expectedStruct{
				paginationPage: alias.BashLogLimitOffsetPage{},
				err:            httpErrors.BashDoesNotExists,
			},
		},
		{
			name: "Getting bash log pagination page error",
			in: inStruct{
				ctx:              context.Background(),
				bashId:           uuid.NewV4(),
				paginationParams: pagination.LimitOffsetParams{},
			},
			mockBehavior: func(mbl *mock_service.MockIBashLogService, mb *mock_service.MockIBashService, ctx context.Context, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams) {
				gomock.InOrder(
					mb.EXPECT().GetOneById(ctx, bashId).Return(&model.Bash{}, nil),
					mbl.EXPECT().GetPaginationPageByBashId(ctx, bashId, paginationParams).Return(alias.BashLogLimitOffsetPage{}, httpErrors.BashLogGetPaginationPageByBashId),
				)
			},
			expected: expectedStruct{
				paginationPage: alias.BashLogLimitOffsetPage{},
				err:            httpErrors.BashLogGetPaginationPageByBashId,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBashService := mock_service.NewMockIBashService(ctrl)
			mockBashLogService := mock_service.NewMockIBashLogService(ctrl)
			testCase.mockBehavior(mockBashLogService, mockBashService, testCase.in.ctx, testCase.in.bashId, testCase.in.paginationParams)

			bashLogUseCase := BashLogUseCase{
				service:     mockBashLogService,
				bashService: mockBashService,
				httpErrors:  httpErrors,
			}

			bashLogPaginationPage, err := bashLogUseCase.GetBashLogPaginationPageByBashId(testCase.in.bashId, testCase.in.paginationParams)

			assert.Equal(t, testCase.expected.paginationPage, bashLogPaginationPage)
			assert.Equal(t, testCase.expected.err, err)
		})
	}
}
