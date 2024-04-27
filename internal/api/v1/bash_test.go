package v1

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	mock_api "pg-sh-scripts/internal/api/mock"
	"pg-sh-scripts/internal/config"
	"pg-sh-scripts/internal/model"
	"pg-sh-scripts/internal/schema"
	"pg-sh-scripts/internal/type/alias"
	mock_usecase "pg-sh-scripts/internal/usecase/mock"
	"pg-sh-scripts/pkg/sql/pagination"
	"strconv"
	"strings"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBashHandler_GetBashById(t *testing.T) {
	type (
		inStruct struct {
			bashId  string
			httpErr error
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
		mockBehavior func(*mock_usecase.MockIBashUseCase, *mock_api.MockIHelper, uuid.UUID, error)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				bashId:  uuid.NewV4().String(),
				httpErr: nil,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
				mu.EXPECT().GetBashById(bashId).Return(&model.Bash{}, nil)
			},
			expected: expectedStruct{
				body: `{"id":"00000000-0000-0000-0000-000000000000","title":"","body":"","createdAt":"0001-01-01T00:00:00Z"}`,
				code: http.StatusOK,
			},
		},
		{
			name: "Bash id must be uuid error",
			in: inStruct{
				bashId:  "uuid",
				httpErr: httpErrors.BashId,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
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
			name: "Bash does not exists error",
			in: inStruct{
				bashId:  uuid.NewV4().String(),
				httpErr: httpErrors.BashDoesNotExists,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)

				gomock.InOrder(
					mu.EXPECT().GetBashById(bashId).Return(nil, err),
					mh.EXPECT().ParseError(err).Return(httpErr),
				)
			},
			expected: expectedStruct{
				body: `{"httpCode":404,"serviceCode":207,"detail":"The specified bash script does not exists"}`,
				code: http.StatusNotFound,
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
			uuidBashId, _ := uuid.FromString(testCase.in.bashId)
			testCase.mockBehavior(mockBashUseCase, mockApiHelper, uuidBashId, testCase.in.httpErr)

			bashHandler := BashHandler{
				useCase:    mockBashUseCase,
				helper:     mockApiHelper,
				httpErrors: httpErrors,
			}

			path := groupBashPath + getBashByIdPath
			casePath := strings.Replace(path, ":id", testCase.in.bashId, 1)

			r := gin.New()
			r.GET(path, bashHandler.GetBashById)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, casePath, nil)

			r.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, testCase.expected.body, recorder.Body.String())
		})
	}
}

func TestBashHandler_GetBashFileById(t *testing.T) {
	type (
		inStruct struct {
			bashId  string
			httpErr error
		}

		expectedStruct struct {
			header http.Header
			body   string
			code   int
		}
	)

	httpErrors := config.GetHTTPErrors()

	testCases := []struct {
		name         string
		in           inStruct
		mockBehavior func(*mock_usecase.MockIBashUseCase, *mock_api.MockIHelper, uuid.UUID, error)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				bashId:  uuid.NewV4().String(),
				httpErr: nil,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
				mu.EXPECT().GetBashFileBufferById(bashId).Return(&bytes.Buffer{}, "", nil)
			},
			expected: expectedStruct{
				header: http.Header{
					"Content-Disposition": []string{"attachment; filename=\".sh\""},
					"Content-Length":      []string{"0"},
					"Content-Type":        []string{"application/x-www-form-urlencoded"},
				},
				body: "",
				code: http.StatusOK,
			},
		},
		{
			name: "Bash id must be uuid error",
			in: inStruct{
				bashId:  "uuid",
				httpErr: httpErrors.BashId,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				mh.EXPECT().ParseError(err).Return(httpErr)
			},
			expected: expectedStruct{
				header: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
				body:   `{"httpCode":422,"serviceCode":200,"detail":"The bash id must be of type uuid4 like 151a583c-0ea0-46b8-b8a6-6bdcdd51655a"}`,
				code:   http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Bash does not exists error",
			in: inStruct{
				bashId:  uuid.NewV4().String(),
				httpErr: httpErrors.BashDoesNotExists,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)

				gomock.InOrder(
					mu.EXPECT().GetBashFileBufferById(bashId).Return(nil, "", err),
					mh.EXPECT().ParseError(err).Return(httpErr),
				)
			},
			expected: expectedStruct{
				header: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
				body:   `{"httpCode":404,"serviceCode":207,"detail":"The specified bash script does not exists"}`,
				code:   http.StatusNotFound,
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
			uuidBashId, _ := uuid.FromString(testCase.in.bashId)
			testCase.mockBehavior(mockBashUseCase, mockApiHelper, uuidBashId, testCase.in.httpErr)

			bashHandler := BashHandler{
				useCase:    mockBashUseCase,
				helper:     mockApiHelper,
				httpErrors: httpErrors,
			}

			path := groupBashPath + getBashFileByIdPath
			casePath := strings.Replace(path, ":id", testCase.in.bashId, 1)

			r := gin.New()
			r.GET(path, bashHandler.GetBashFileById)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, casePath, nil)

			r.ServeHTTP(recorder, request)
			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, testCase.expected.header, recorder.Result().Header)
			assert.Equal(t, testCase.expected.body, recorder.Body.String())
		})
	}
}

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

func TestBashHandler_RemoveBashById(t *testing.T) {
	type (
		inStruct struct {
			bashId  string
			httpErr error
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
		mockBehavior func(*mock_usecase.MockIBashUseCase, *mock_api.MockIHelper, uuid.UUID, error)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				bashId:  uuid.NewV4().String(),
				httpErr: nil,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
				mu.EXPECT().RemoveBashById(bashId).Return(&model.Bash{}, nil)
			},
			expected: expectedStruct{
				body: `{"id":"00000000-0000-0000-0000-000000000000","title":"","body":"","createdAt":"0001-01-01T00:00:00Z"}`,
				code: http.StatusOK,
			},
		},
		{
			name: "Bash id must be uuid error",
			in: inStruct{
				bashId:  "uuid",
				httpErr: httpErrors.BashId,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
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
			name: "Removing bash error",
			in: inStruct{
				bashId:  uuid.NewV4().String(),
				httpErr: httpErrors.BashRemove,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, bashId uuid.UUID, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)

				gomock.InOrder(
					mu.EXPECT().RemoveBashById(bashId).Return(nil, err),
					mh.EXPECT().ParseError(err).Return(httpErr),
				)
			},
			expected: expectedStruct{
				body: `{"httpCode":400,"serviceCode":211,"detail":"An error occurred while deleting the bash script"}`,
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
			uuidBashId, _ := uuid.FromString(testCase.in.bashId)
			testCase.mockBehavior(mockBashUseCase, mockApiHelper, uuidBashId, testCase.in.httpErr)

			bashHandler := BashHandler{
				useCase:    mockBashUseCase,
				helper:     mockApiHelper,
				httpErrors: httpErrors,
			}

			path := groupBashPath + removeBashPath
			casePath := strings.Replace(path, ":id", testCase.in.bashId, 1)

			r := gin.New()
			r.DELETE(path, bashHandler.RemoveBashById)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodDelete, casePath, nil)

			r.ServeHTTP(recorder, request)

			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, testCase.expected.body, recorder.Body.String())
		})
	}
}
