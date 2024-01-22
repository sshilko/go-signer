package db

import (
	"github.com/nutsdb/nutsdb"
	"github.com/pkg/errors"
)

type Database struct {
	db     *nutsdb.DB
	bucket string
}

var ErrNotFound = errors.New("Not found")

func NewDB(path, bucket string) (*Database, error) {
	db, err := nutsdb.Open(
		nutsdb.DefaultOptions,
		nutsdb.WithDir(path),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to access db")
	}

	if err := db.Update(
		func(tx *nutsdb.Tx) error {
			return tx.NewBucket(nutsdb.DataStructureBTree, bucket)
		}); err != nil {
		if !errors.Is(err, nutsdb.ErrBucketAlreadyExist) {
			return nil, err
		}
	}

	return &Database{
		db:     db,
		bucket: bucket,
	}, nil
}

func (d *Database) Disconnect() error {
	return d.db.Close()
}

func (d *Database) Get(key []byte) ([]byte, error) {
	var result []byte
	err := d.db.View(
		func(tx *nutsdb.Tx) error {
			value, err := tx.Get(d.bucket, key)
			if err != nil {
				if errors.Is(err, nutsdb.ErrKeyNotFound) ||
					errors.Is(err, nutsdb.ErrBucketNotFound) ||
					errors.Is(err, nutsdb.ErrNotFoundBucket) ||
					errors.Is(err, nutsdb.ErrBucketNotExist) {
					return ErrNotFound
				}
				return err
			}
			result = value
			return nil
		})
	return result, err
}

func (d *Database) Set(key, value []byte) error {
	return d.db.Update(
		func(tx *nutsdb.Tx) error {
			return tx.Put(d.bucket, key, value, 0)
		})
}
