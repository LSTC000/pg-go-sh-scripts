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
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestBashLogHandler_GetBashLogListByBashId(t *testing.T) {
	type (
		inStruct struct {
			bashId           string
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
		mockBehavior func(*mock_usecase.MockIBashLogUseCase, *mock_api.MockIHelper, uuid.UUID, pagination.LimitOffsetParams, error)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				bashId:           uuid.NewV4().String(),
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          nil,
				limitExists:      true,
				offsetExists:     true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashLogUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams, err error) {
				mu.EXPECT().GetBashLogPaginationPageByBashId(bashId, paginationParams).Return(alias.BashLogLimitOffsetPage{}, nil)
			},
			expected: expectedStruct{
				body: `{"items":null,"limit":0,"offset":0,"total":0}`,
				code: http.StatusOK,
			},
		},
		{
			name: "Bash id must be uuid error",
			in: inStruct{
				bashId:           "uuid",
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          httpErrors.BashId,
				limitExists:      true,
				offsetExists:     true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashLogUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				mh.EXPECT().ParseError(err).Return(httpErr)
			},
			expected: expectedStruct{
				body: `{"httpCode":422,"serviceCode":200,"detail":"The bash id must be of type uuid4 like 151a583c-0ea0-46b8-b8a6-6bdcdd51655a"}`,
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Limit param must be int error",
			in: inStruct{
				bashId:           uuid.NewV4().String(),
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          httpErrors.PaginationLimitParamMustBeInt,
				limitExists:      false,
				offsetExists:     true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashLogUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams, err error) {
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
				bashId: uuid.NewV4().String(),
				paginationParams: pagination.LimitOffsetParams{
					Limit: -1,
				},
				httpErr:      httpErrors.PaginationLimitParamGTEZero,
				limitExists:  true,
				offsetExists: true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashLogUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams, err error) {
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
				bashId:           uuid.NewV4().String(),
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          httpErrors.PaginationOffsetParamMustBeInt,
				limitExists:      true,
				offsetExists:     false,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashLogUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams, err error) {
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
				bashId: uuid.NewV4().String(),
				paginationParams: pagination.LimitOffsetParams{
					Offset: -1,
				},
				httpErr:      httpErrors.PaginationOffsetParamGTEZero,
				limitExists:  true,
				offsetExists: true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashLogUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams, err error) {
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
			name: "Getting bash log pagination page error",
			in: inStruct{
				bashId:           uuid.NewV4().String(),
				paginationParams: pagination.LimitOffsetParams{},
				httpErr:          httpErrors.BashLogGetPaginationPageByBashId,
				limitExists:      true,
				offsetExists:     true,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashLogUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, paginationParams pagination.LimitOffsetParams, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				gomock.InOrder(
					mu.EXPECT().GetBashLogPaginationPageByBashId(bashId, paginationParams).Return(alias.BashLogLimitOffsetPage{}, err),
					mh.EXPECT().ParseError(err).Return(httpErr),
				)
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

			mockBashLogUseCase := mock_usecase.NewMockIBashLogUseCase(ctrl)
			mockApiHelper := mock_api.NewMockIHelper(ctrl)
			uuidBashId, _ := uuid.FromString(testCase.in.bashId)
			testCase.mockBehavior(mockBashLogUseCase, mockApiHelper, uuidBashId, testCase.in.paginationParams, testCase.in.httpErr)

			bashLogHandler := BashLogHandler{
				useCase:    mockBashLogUseCase,
				helper:     mockApiHelper,
				httpErrors: httpErrors,
			}

			path := groupBashLogPath + getBashLogListByBashIdPath
			casePath := strings.Replace(path, ":bashId", testCase.in.bashId, 1)

			r := gin.New()
			r.GET(path, bashLogHandler.GetBashLogListByBashId)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, casePath, nil)

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
