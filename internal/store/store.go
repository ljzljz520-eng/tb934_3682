package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	bolt "go.etcd.io/bbolt"
	"weddingguide/internal/domain"
)

var (
	ErrReaderLimit = errors.New("activity reader limit reached")
	ErrNotFound    = errors.New("entity not found")
)

const (
	guideBucket    = "WeddingGuide"
	scheduleBucket = "ScheduleItem"
	venueBucket    = "VenueAddress"
	attireBucket   = "AttireTip"
	linkBucket     = "ActionLink"
	visitorBucket  = "VisitorRecord"
	blessingBucket = "Blessing"
	auditBucket    = "AuditEntry"
)

type Store struct {
	db          *bolt.DB
	readerSlots chan struct{}
	mu          sync.RWMutex
}

func New(path string) (*Store, error) {
	return NewWithReaderLimit(path, 4)
}

func NewWithReaderLimit(path string, limit int) (*Store, error) {
	if limit < 1 {
		return nil, fmt.Errorf("reader limit must be positive")
	}
	if err := ensureParent(path); err != nil {
		return nil, err
	}
	db, err := bolt.Open(path, 0600, &bolt.Options{NoSync: true})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db, readerSlots: make(chan struct{}, limit)}
	if err := s.initialize(); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func ensureParent(path string) error {
	parent := filepath.Dir(path)
	if parent == "." || parent == "" {
		return nil
	}
	return os.MkdirAll(parent, 0755)
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range []string{guideBucket, scheduleBucket, venueBucket, attireBucket, linkBucket, visitorBucket, blessingBucket, auditBucket} {
			if _, err := tx.CreateBucketIfNotExists([]byte(name)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *Store) SaveGuide(guide domain.WeddingGuide) error {
	if err := guide.Validate(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := putJSON(tx, guideBucket, guide.ID, guide); err != nil {
			return err
		}
		if err := putJSON(tx, venueBucket, guide.ID, guide.Venue); err != nil {
			return err
		}
		if err := putJSON(tx, attireBucket, guide.ID, guide.Attire); err != nil {
			return err
		}
		for _, item := range guide.Schedule {
			if err := putJSON(tx, scheduleBucket, guide.ID+":"+item.ID, item); err != nil {
				return err
			}
		}
		for _, link := range guide.Links {
			if err := putJSON(tx, linkBucket, guide.ID+":"+link.ID, link); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) ReadGuide(id string) (domain.WeddingGuide, error) {
	var guide domain.WeddingGuide
	err := s.db.View(func(tx *bolt.Tx) error {
		return getJSON(tx, guideBucket, id, &guide)
	})
	if err != nil {
		return domain.WeddingGuide{}, err
	}
	return guide, nil
}

func (s *Store) SaveVisitor(visitor domain.VisitorRecord) error {
	if err := visitor.Validate(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return putJSON(tx, visitorBucket, visitor.GuideID+":"+visitor.VisitorKey, visitor)
	})
}

func (s *Store) ReadVisitor(guideID, visitorKey string) (domain.VisitorRecord, error) {
	var visitor domain.VisitorRecord
	err := s.db.View(func(tx *bolt.Tx) error {
		return getJSON(tx, visitorBucket, guideID+":"+visitorKey, &visitor)
	})
	if err != nil {
		return domain.VisitorRecord{}, err
	}
	return visitor, nil
}

func (s *Store) SaveBlessing(blessing domain.Blessing) error {
	if err := blessing.Validate(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return putJSON(tx, blessingBucket, blessing.ID, blessing)
	})
}

func (s *Store) ReadBlessing(id string) (domain.Blessing, error) {
	var blessing domain.Blessing
	err := s.db.View(func(tx *bolt.Tx) error {
		return getJSON(tx, blessingBucket, id, &blessing)
	})
	if err != nil {
		return domain.Blessing{}, err
	}
	return blessing, nil
}

func (s *Store) AppendAudit(entry domain.AuditEntry) error {
	if err := entry.Validate(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		return putJSON(tx, auditBucket, entry.ID, entry)
	})
}

func (s *Store) ListAudits(guideID string) ([]domain.AuditEntry, error) {
	entries := make([]domain.AuditEntry, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(auditBucket))
		return bucket.ForEach(func(_, value []byte) error {
			var entry domain.AuditEntry
			if err := json.Unmarshal(value, &entry); err != nil {
				return err
			}
			if entry.GuideID == guideID {
				entries = append(entries, entry)
			}
			return nil
		})
	})
	return entries, err
}

func putJSON(tx *bolt.Tx, bucket, key string, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return tx.Bucket([]byte(bucket)).Put([]byte(key), encoded)
}

func getJSON(tx *bolt.Tx, bucket, key string, destination any) error {
	value := tx.Bucket([]byte(bucket)).Get([]byte(key))
	if value == nil {
		return ErrNotFound
	}
	return json.Unmarshal(value, destination)
}
