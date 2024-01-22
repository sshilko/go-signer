package homework

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sshilko/go-signer/pkg/db"
)

type Persistence interface {
	SaveHomework(homework Homework) error
	FetchHomeworkByID(homeworkID uuid.UUID) (*Homework, error)
}

var NotFoundError = errors.New("Homework not found")

type Service struct {
	Repo Persistence
}

func NewService(repo Persistence) *Service {
	return &Service{
		Repo: repo,
	}
}

func (d *Service) GetByID(homeworkID uuid.UUID) (*Homework, error) {
	ok, err := d.Repo.FetchHomeworkByID(homeworkID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, NotFoundError
		}
	}
	return ok, err
}

func (d *Service) Save(questions, answers string, userID uuid.UUID) (*uuid.UUID, error) {
	homeworkRecord := NewHomework(questions, answers, userID)
	if err := d.Repo.SaveHomework(homeworkRecord); err != nil {
		return nil, err
	}
	return &homeworkRecord.ID, nil
}
