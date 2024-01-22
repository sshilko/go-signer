package web_test

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sshilko/go-signer/handler/web"
	"github.com/sshilko/go-signer/service/homework"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockedHomeworkService struct {
	mock.Mock
}

func (m *MockedHomeworkService) Save(questions, answers string, userID uuid.UUID) (*uuid.UUID, error) {
	args := m.Called()
	return args.Get(0).(*uuid.UUID), args.Error(1)
}

func (m *MockedHomeworkService) GetByID(homeworkID uuid.UUID) (*homework.Homework, error) {
	args := m.Called()
	return args.Get(0).(*homework.Homework), args.Error(1)
}

func TestRetrieveHomeworkAssignment(t *testing.T) {
	myEcho := echo.New()
	myRec := httptest.NewRecorder()

	payload := map[string]string{
		"questions": "Q1,Q2,C0,X1",
		"answers":   "A1,A2,C1,B2",
	}
	payloadBody, _ := json.Marshal(payload)

	requestPath := "/users/:userID/homework"

	testRequest := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(payloadBody))
	testRequest.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	testEchoContext := myEcho.NewContext(testRequest, myRec)

	userID1 := uuid.New()

	testEchoContext.SetPath(requestPath)
	testEchoContext.SetParamNames("userID")
	testEchoContext.SetParamValues(userID1.String())

	t.Run("Success", func(t *testing.T) {
		persistance := new(MockedHomeworkService)
		persistance.Mock.On("Save").Return(&userID1, nil)

		h := web.NewHomeworkHandler(persistance)
		assert.NoError(t, h.CreateHomeworkAssignment(testEchoContext))

		persistance.AssertCalled(t, "Save")
		require.Equal(t, http.StatusCreated, myRec.Code)

		type expectedResponse struct {
			ID uuid.UUID `json:"id"`
		}
		var expectResp expectedResponse
		require.NoError(t, json.Unmarshal(myRec.Body.Bytes(), &expectResp))
		assert.Equal(t, expectResp.ID.String(), userID1.String())
	})
}

func TestCreateHomeworkAssignmentNotFound(t *testing.T) {
	myEcho := echo.New()
	myRec := httptest.NewRecorder()

	requestPath := "/users/:userID/homework/:homeworkID"

	testRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	testEchoContext := myEcho.NewContext(testRequest, myRec)

	userID1 := uuid.New()
	testID1 := uuid.New()

	testEchoContext.SetPath(requestPath)
	testEchoContext.SetParamNames("userID", "homeworkID")
	testEchoContext.SetParamValues(userID1.String(), testID1.String())

	t.Run("Not found", func(t *testing.T) {
		persistence := new(MockedHomeworkService)
		persistence.Mock.On("GetByID").Return(&homework.Homework{}, homework.NotFoundError)

		h := web.NewHomeworkHandler(persistence)
		assert.NoError(t, h.RetrieveHomeworkAssignment(testEchoContext))

		persistence.AssertCalled(t, "GetByID")
		require.Equal(t, http.StatusNotFound, myRec.Code)
	})
}

func TestCreateHomeworkAssignmentSuccess(t *testing.T) {
	myEcho := echo.New()
	myRec := httptest.NewRecorder()

	requestPath := "/users/:userID/homework/:homeworkID"

	testRequest := httptest.NewRequest(http.MethodGet, "/", nil)
	testEchoContext := myEcho.NewContext(testRequest, myRec)

	userID1 := uuid.New()
	testID1 := uuid.New()

	testEchoContext.SetPath(requestPath)
	testEchoContext.SetParamNames("userID", "homeworkID")
	testEchoContext.SetParamValues(userID1.String(), testID1.String())

	t.Run("Success", func(t *testing.T) {
		nowTime := time.Now()

		persistence := new(MockedHomeworkService)
		dummyHomework := homework.Homework{
			Questions:   "Questions " + uuid.NewString(),
			Answers:     "Answers " + uuid.NewString(),
			ID:          testID1,
			UserID:      userID1,
			TimeCreated: &nowTime,
		}
		persistence.Mock.On("GetByID").Return(&dummyHomework, nil)

		h := web.NewHomeworkHandler(persistence)
		assert.NoError(t, h.RetrieveHomeworkAssignment(testEchoContext))

		persistence.AssertCalled(t, "GetByID")
		require.Equal(t, http.StatusOK, myRec.Code)
	})
}
