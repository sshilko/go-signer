package homework

import (
	"encoding/json"
	"github.com/google/uuid"
)

type RepositoryStorage interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
}

type Repository struct {
	db RepositoryStorage
}

func NewRepository(db RepositoryStorage) *Repository {
	return &Repository{
		db: db,
	}
}

func (d *Repository) FetchHomeworkByID(homeworkID uuid.UUID) (*Homework, error) {
	id, err := homeworkID.MarshalBinary()
	if err != nil {
		return nil, err
	}

	binary, err := d.db.Get(id)
	if err != nil {
		return nil, err
	}

	var homework Homework
	if err := json.Unmarshal(binary, &homework); err != nil {
		return nil, err
	}

	return &homework, nil
}

func (d *Repository) SaveHomework(homework Homework) error {
	homeworkSerialized, err := json.Marshal(homework)
	if err != nil {
		return err
	}

	homeworkID, err := homework.ID.MarshalBinary()
	if err != nil {
		return err
	}

	err = d.db.Set(homeworkID, homeworkSerialized)
	if err != nil {
		return err
	}

	return nil
}
