package homework

import (
	"github.com/google/uuid"
	"time"
)

type Homework struct {
	Questions   string     `json:"questions"`
	Answers     string     `json:"answers"`
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"user_id"`
	TimeCreated *time.Time `json:"created_at,omitempty"`
}

func NewHomework(questions, answers string, userID uuid.UUID) Homework {
	now := time.Now().UTC()
	return Homework{
		ID:          uuid.New(),
		Questions:   questions,
		Answers:     answers,
		UserID:      userID,
		TimeCreated: &now,
	}
}
