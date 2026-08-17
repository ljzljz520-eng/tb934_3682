package store

import (
	"sync"

	bolt "go.etcd.io/bbolt"
)

type Reader struct {
	store *Store
	once  sync.Once
	open  bool
}

func (s *Store) OpenReader() (*Reader, error) {
	select {
	case s.readerSlots <- struct{}{}:
	default:
		return nil, ErrReaderLimit
	}
	return &Reader{store: s, open: true}, nil
}

func (r *Reader) Get(bucket, key string) ([]byte, error) {
	if r == nil || !r.open {
		return nil, ErrNotFound
	}
	var copyValue []byte
	err := r.store.db.View(func(tx *bolt.Tx) error {
		value := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if value == nil {
			return ErrNotFound
		}
		copyValue = append([]byte(nil), value...)
		return nil
	})
	return copyValue, err
}

func (r *Reader) Has(bucket, key string) bool {
	_, err := r.Get(bucket, key)
	return err == nil
}

func (r *Reader) Close() error {
	if r == nil {
		return nil
	}
	r.once.Do(func() {
		<-r.store.readerSlots
		r.open = false
	})
	return nil
}
