package web

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sshilko/go-signer/service/homework"
	"net/http"
	"time"
)

type HomeworkHandlerService interface {
	GetByID(homeworkID uuid.UUID) (*homework.Homework, error)
	Save(questions, answers string, userID uuid.UUID) (*uuid.UUID, error)
}

type HomeworkHandler struct {
	svc HomeworkHandlerService
}

func NewHomeworkHandler(svc HomeworkHandlerService) *HomeworkHandler {
	return &HomeworkHandler{svc: svc}
}

// GetRoutes returns routes defined on this handler.
func (m *HomeworkHandler) GetRoutes() map[string]map[string]echo.HandlerFunc {
	return map[string]map[string]echo.HandlerFunc{
		http.MethodPost: {
			"/users/:userID/homework": m.CreateHomeworkAssignment,
		},
		http.MethodGet: {
			"/users/:userID/homework/:homeworkID": m.RetrieveHomeworkAssignment,
		},
	}
}

type HomeworkFetchResponse struct {
	Questions   string    `json:"questions"`
	Answers     string    `json:"answers"`
	TimeCreated time.Time `json:"created_at"`
}

func (m *HomeworkHandler) RetrieveHomeworkAssignment(c echo.Context) error {
	userID := c.Param("userID")
	user, err := uuid.Parse(userID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	testID := c.Param("homeworkID")
	test, err := uuid.Parse(testID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	assignment, err := m.svc.GetByID(test)
	if err != nil {
		if errors.Is(err, homework.NotFoundError) {
			return c.NoContent(http.StatusNotFound)
		}
		return errors.Wrap(err, "failed to fetch test")
	}

	if assignment.UserID != user {
		return c.NoContent(http.StatusForbidden)
	}
	return c.JSON(http.StatusOK, &HomeworkFetchResponse{
		Questions:   assignment.Questions,
		Answers:     assignment.Answers,
		TimeCreated: *assignment.TimeCreated,
	})
}

type HomeworkSubmissionResponse struct {
	ID uuid.UUID `json:"id"`
}

type HomeworkRequest struct {
	Questions string `json:"questions" validate:"required"`
	Answers   string `json:"answers" validate:"required"`
}

func (m *HomeworkHandler) CreateHomeworkAssignment(c echo.Context) error {
	userID := c.Param("userID")
	user := uuid.MustParse(userID)

	input := new(HomeworkRequest)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	validate := validator.New()
	if err := validate.Struct(input); err != nil {
		return errors.Wrap(err, "request payload invalid")
	}

	id, err := m.svc.Save(input.Questions, input.Answers, user)
	if err != nil {
		return errors.Wrap(err, "failed to submit test")
	}

	return c.JSON(http.StatusCreated, &HomeworkSubmissionResponse{ID: *id})
}
