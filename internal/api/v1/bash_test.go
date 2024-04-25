package v1

import (
	"errors"
	"net/http"
	"net/http/httptest"
	mock_api "pg-sh-scripts/internal/api/mock"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/internal/type/alias"
	mock_usecase "pg-sh-scripts/internal/usecase/mock"
	"pg-sh-scripts/pkg/sql/pagination"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBashHandler_GetBashList(t *testing.T) {
	type (
		inStruct struct {
			paginationParams pagination.LimitOffsetParams
			httpErr          error
			limitExists      bool
			offsetExists     bool
		}

		expectedStruct struct {
			body string
			code int
		}
	)

	httpErrors := config.GetHTTPErrors()

	testCases := []struct {
		name         string
		in           inStruct
		mockBehavior func(*mock_usecase.MockIBashUseCase, *mock_api.MockIHelper, pagination.LimitOffsetParams, error)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          nil,
				limitExists:      true,
				offsetExists:     true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, paginationParams pagination.LimitOffsetParams, err error) {
				mu.EXPECT().GetBashPaginationPage(paginationParams).Return(alias.BashLimitOffsetPage{}, nil)
			},
			expected: expectedStruct{
				body: `{"items":null,"limit":0,"offset":0,"total":0}`,
				code: http.StatusOK,
			},
		},
		{
			name: "Limit param must be int error",
			in: inStruct{
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          httpErrors.PaginationLimitParamMustBeInt,
				limitExists:      false,
				offsetExists:     true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, paginationParams pagination.LimitOffsetParams, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				mh.EXPECT().ParseError(err).Return(httpErr)
			},
			expected: expectedStruct{
				body: `{"httpCode":422,"serviceCode":100,"detail":"The limit pagination parameter must be integer"}`,
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Limit param gte to zero error",
			in: inStruct{
				paginationParams: pagination.LimitOffsetParams{
					Limit: -1,
				},
				httpErr:      httpErrors.PaginationLimitParamGTEZero,
				limitExists:  true,
				offsetExists: true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, paginationParams pagination.LimitOffsetParams, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				mh.EXPECT().ParseError(err).Return(httpErr)
			},
			expected: expectedStruct{
				body: `{"httpCode":422,"serviceCode":101,"detail":"The limit pagination parameter must be greater than or equal to zero"}`,
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Offset param must be int error",
			in: inStruct{
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          httpErrors.PaginationOffsetParamMustBeInt,
				limitExists:      true,
				offsetExists:     false,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, paginationParams pagination.LimitOffsetParams, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				mh.EXPECT().ParseError(err).Return(httpErr)
			},
			expected: expectedStruct{
				body: `{"httpCode":422,"serviceCode":102,"detail":"The offset pagination parameter must be integer"}`,
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Limit param gte to zero error",
			in: inStruct{
				paginationParams: pagination.LimitOffsetParams{
					Offset: -1,
				},
				httpErr:      httpErrors.PaginationOffsetParamGTEZero,
				limitExists:  true,
				offsetExists: true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, paginationParams pagination.LimitOffsetParams, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				mh.EXPECT().ParseError(err).Return(httpErr)
			},
			expected: expectedStruct{
				body: `{"httpCode":422,"serviceCode":103,"detail":"The offset pagination parameter must be greater than or equal to zero"}`,
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Getting bash pagination page error",
			in: inStruct{
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          httpErrors.BashLogGetPaginationPageByBashId,
				limitExists:      true,
				offsetExists:     true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, paginationParams pagination.LimitOffsetParams, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				mu.EXPECT().GetBashPaginationPage(paginationParams).Return(alias.BashLimitOffsetPage{}, err)
				mh.EXPECT().ParseError(err).Return(httpErr)
			},
			expected: expectedStruct{
				body: `{"httpCode":400,"serviceCode":300,"detail":"An error occurred while receiving the pagination page of bash log scripts"}`,
				code: http.StatusBadRequest,
			},
		},
	}

	gin.SetMode(gin.TestMode)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockBashUseCase := mock_usecase.NewMockIBashUseCase(ctrl)
			mockApiHelper := mock_api.NewMockIHelper(ctrl)
			testCase.mockBehavior(mockBashUseCase, mockApiHelper, testCase.in.paginationParams, testCase.in.httpErr)

			bashHandler := BashHandler{
				useCase:    mockBashUseCase,
				helper:     mockApiHelper,
				httpErrors: httpErrors,
			}

			path := groupBashPath + getBashListPath

			r := gin.New()
			r.GET(path, bashHandler.GetBashList)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, path, nil)

			requestQueryParams := request.URL.Query()
			if testCase.in.limitExists {
				requestQueryParams.Add("limit", strconv.Itoa(testCase.in.paginationParams.Limit))
			}
			if testCase.in.offsetExists {
				requestQueryParams.Add("offset", strconv.Itoa(testCase.in.paginationParams.Offset))
			}
			request.URL.RawQuery = requestQueryParams.Encode()

			r.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, testCase.expected.body, recorder.Body.String())
		})
	}
}
