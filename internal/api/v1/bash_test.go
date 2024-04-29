package v1

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
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

const (
	bashTestFile    = "helloworld.sh"
	bashTestDataDir = "bash_testdata"
)

func TestBashHandler_GetBashById(t *testing.T) {
	type (
		inStruct struct {
			bashId  string
			httpErr error
		}

		expectedStruct struct {
			golden string
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
				mu.EXPECT().GetBashById(bashId).Return(&model.Bash{}, nil)
			},
			expected: expectedStruct{
				golden: "default_bash",
				code:   http.StatusOK,
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
				golden: "bash_id_error",
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
					mu.EXPECT().GetBashById(bashId).Return(nil, err),
					mh.EXPECT().ParseError(err).Return(httpErr),
				)
			},
			expected: expectedStruct{
				golden: "bash_does_not_exists_error",
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

			handlerPath := groupBashPath + getBashByIdPath
			handlerCasePath := strings.Replace(handlerPath, ":id", testCase.in.bashId, 1)

			r := gin.New()
			r.GET(handlerPath, bashHandler.GetBashById)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, handlerCasePath, nil)

			r.ServeHTTP(recorder, request)

			content, err := os.ReadFile(path.Join(bashTestDataDir, testCase.expected.golden+".golden"))
			if err != nil {
				t.Fatalf("%s Error: %s", t.Name(), err)
			}
			expectedBody := string(content)

			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, expectedBody, recorder.Body.String())
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
			golden string
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
				golden: "default_bash_file",
				code:   http.StatusOK,
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
				golden: "bash_id_error",
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
				golden: "bash_does_not_exists_error",
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

			handlerPath := groupBashPath + getBashFileByIdPath
			handlerCasePath := strings.Replace(handlerPath, ":id", testCase.in.bashId, 1)

			r := gin.New()
			r.GET(handlerPath, bashHandler.GetBashFileById)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, handlerCasePath, nil)

			r.ServeHTTP(recorder, request)

			content, err := os.ReadFile(path.Join(bashTestDataDir, testCase.expected.golden+".golden"))
			if err != nil {
				t.Fatalf("%s Error: %s", t.Name(), err)
			}
			expectedBody := string(content)

			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, testCase.expected.header, recorder.Result().Header)
			assert.Equal(t, expectedBody, recorder.Body.String())
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
			golden string
			code   int
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
				golden: "default_pagination_page",
				code:   http.StatusOK,
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
				golden: "pagination_limit_param_int_error",
				code:   http.StatusUnprocessableEntity,
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
				golden: "pagination_limit_param_gte_zero_error",
				code:   http.StatusUnprocessableEntity,
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
				golden: "pagination_offset_param_int_error",
				code:   http.StatusUnprocessableEntity,
			},
		},
		{
			name: "Offset param gte to zero error",
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
				golden: "pagination_offset_param_gte_zero_error",
				code:   http.StatusUnprocessableEntity,
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
				golden: "get_pagination_page_error",
				code:   http.StatusBadRequest,
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

			handlerPath := groupBashPath + getBashListPath

			r := gin.New()
			r.GET(handlerPath, bashHandler.GetBashList)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, handlerPath, nil)

			requestQueryParams := request.URL.Query()
			if testCase.in.limitExists {
				requestQueryParams.Add("limit", strconv.Itoa(testCase.in.paginationParams.Limit))
			}
			if testCase.in.offsetExists {
				requestQueryParams.Add("offset", strconv.Itoa(testCase.in.paginationParams.Offset))
			}
			request.URL.RawQuery = requestQueryParams.Encode()

			r.ServeHTTP(recorder, request)

			content, err := os.ReadFile(path.Join(bashTestDataDir, testCase.expected.golden+".golden"))
			if err != nil {
				t.Fatalf("%s Error: %s", t.Name(), err)
			}
			expectedBody := string(content)

			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, expectedBody, recorder.Body.String())
		})
	}
}

func TestBashHandler_CreateBash(t *testing.T) {
	type (
		inStruct struct {
			isUploadFile bool
			httpErr      error
		}

		expectedStruct struct {
			golden string
			code   int
		}
	)

	httpErrors := config.GetHTTPErrors()

	testCases := []struct {
		name         string
		in           inStruct
		mockBehavior func(*mock_usecase.MockIBashUseCase, *mock_api.MockIHelper, error)
		expected     expectedStruct
	}{
		{
			name: "Success",
			in: inStruct{
				isUploadFile: true,
				httpErr:      nil,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, err error) {
				mu.EXPECT().CreateBash(gomock.Any()).Return(&model.Bash{}, nil)
			},
			expected: expectedStruct{
				golden: "default_bash",
				code:   http.StatusOK,
			},
		},
		{
			name: "Uploading bash file error",
			in: inStruct{
				isUploadFile: false,
				httpErr:      httpErrors.BashFileUpload,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				mh.EXPECT().ParseError(err).Return(httpErr)
			},
			expected: expectedStruct{
				golden: "upload_file_error",
				code:   http.StatusBadRequest,
			},
		},
		{
			name: "Creating bash error",
			in: inStruct{
				isUploadFile: true,
				httpErr:      httpErrors.BashCreate,
			},
			mockBehavior: func(mu *mock_usecase.MockIBashUseCase, mh *mock_api.MockIHelper, err error) {
				var httpErr *schema.HTTPError
				errors.As(err, &httpErr)
				gomock.InOrder(
					mu.EXPECT().CreateBash(gomock.Any()).Return(nil, err),
					mh.EXPECT().ParseError(err).Return(httpErr),
				)
			},
			expected: expectedStruct{
				golden: "create_error",
				code:   http.StatusBadRequest,
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
			testCase.mockBehavior(mockBashUseCase, mockApiHelper, testCase.in.httpErr)

			bashHandler := BashHandler{
				useCase:    mockBashUseCase,
				helper:     mockApiHelper,
				httpErrors: httpErrors,
			}

			handlerPath := groupBashPath + createBashPath

			r := gin.New()
			r.POST(handlerPath, bashHandler.CreateBash)

			var (
				recorder *httptest.ResponseRecorder
				request  *http.Request
			)

			if testCase.in.isUploadFile {
				var b bytes.Buffer
				w := multipart.NewWriter(&b)

				file, err := os.Open(path.Join(bashTestDataDir, bashTestFile))
				if err != nil {
					t.Fatalf("%s Error: %s", t.Name(), err)
				}
				defer func() {
					if err := file.Close(); err != nil {
						t.Fatalf("%s Error: %s", t.Name(), err)
					}
				}()

				fw, err := w.CreateFormFile("file", file.Name())
				if err != nil {
					t.Fatalf("%s Error: %s", t.Name(), err)
				}

				_, err = io.Copy(fw, file)
				if err != nil {
					t.Fatalf("%s Error: %s", t.Name(), err)
				}
				if err := w.Close(); err != nil {
					t.Fatalf("%s Error: %s", t.Name(), err)
				}

				recorder = httptest.NewRecorder()
				request = httptest.NewRequest(http.MethodPost, handlerPath, &b)
				request.Header.Add("Content-Type", w.FormDataContentType())
			} else {
				recorder = httptest.NewRecorder()
				request = httptest.NewRequest(http.MethodPost, handlerPath, nil)
			}

			r.ServeHTTP(recorder, request)

			content, err := os.ReadFile(path.Join(bashTestDataDir, testCase.expected.golden+".golden"))
			if err != nil {
				t.Fatalf("%s Error: %s", t.Name(), err)
			}
			expectedBody := string(content)

			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, expectedBody, recorder.Body.String())
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
			golden string
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
				mu.EXPECT().RemoveBashById(bashId).Return(&model.Bash{}, nil)
			},
			expected: expectedStruct{
				golden: "default_bash",
				code:   http.StatusOK,
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
				golden: "bash_id_error",
				code:   http.StatusUnprocessableEntity,
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
				golden: "removing_bas_error",
				code:   http.StatusBadRequest,
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

			handlerPath := groupBashPath + removeBashPath
			handlerCasePath := strings.Replace(handlerPath, ":id", testCase.in.bashId, 1)

			r := gin.New()
			r.DELETE(handlerPath, bashHandler.RemoveBashById)

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodDelete, handlerCasePath, nil)

			r.ServeHTTP(recorder, request)

			content, err := os.ReadFile(path.Join(bashTestDataDir, testCase.expected.golden+".golden"))
			if err != nil {
				t.Fatalf("%s Error: %s", t.Name(), err)
			}
			expectedBody := string(content)

			assert.Equal(t, testCase.expected.code, recorder.Code)
			assert.Equal(t, expectedBody, recorder.Body.String())
		})
	}
}
