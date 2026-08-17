package store

import (
	"encoding/json"
	"sort"
	"strings"

	bolt "go.etcd.io/bbolt"
	"weddingguide/internal/domain"
)

type Stats struct {
	Guides    int
	Schedules int
	Venues    int
	Attire    int
	Links     int
	Visitors  int
	Blessings int
	Audits    int
}

type EntityCounts map[string]int

func (s *Store) Stats() (Stats, error) {
	counts, err := s.EntityCounts()
	if err != nil {
		return Stats{}, err
	}
	return Stats{Guides: counts[guideBucket], Schedules: counts[scheduleBucket], Venues: counts[venueBucket], Attire: counts[attireBucket], Links: counts[linkBucket], Visitors: counts[visitorBucket], Blessings: counts[blessingBucket], Audits: counts[auditBucket]}, nil
}

func (s *Store) EntityCounts() (EntityCounts, error) {
	counts := make(EntityCounts)
	err := s.db.View(func(tx *bolt.Tx) error {
		for _, bucket := range []string{guideBucket, scheduleBucket, venueBucket, attireBucket, linkBucket, visitorBucket, blessingBucket, auditBucket} {
			counts[bucket] = tx.Bucket([]byte(bucket)).Stats().KeyN
		}
		return nil
	})
	return counts, err
}

func (s *Store) ListVisitors(guideID string) ([]domain.VisitorRecord, error) {
	visitors := make([]domain.VisitorRecord, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		return forEachJSON(tx, visitorBucket, func(value []byte) error {
			var visitor domain.VisitorRecord
			if err := json.Unmarshal(value, &visitor); err != nil {
				return err
			}
			if visitor.GuideID == guideID {
				visitors = append(visitors, visitor)
			}
			return nil
		})
	})
	sort.SliceStable(visitors, func(i, j int) bool { return visitors[i].VisitorKey < visitors[j].VisitorKey })
	return visitors, err
}

func (s *Store) ListBlessings(guideID string) ([]domain.Blessing, error) {
	blessings := make([]domain.Blessing, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		return forEachJSON(tx, blessingBucket, func(value []byte) error {
			var blessing domain.Blessing
			if err := json.Unmarshal(value, &blessing); err != nil {
				return err
			}
			if blessing.GuideID == guideID {
				blessings = append(blessings, blessing)
			}
			return nil
		})
	})
	sort.SliceStable(blessings, func(i, j int) bool { return blessings[i].CreatedAt < blessings[j].CreatedAt })
	return blessings, err
}

func (s *Store) ExportGuide(guideID string) ([]byte, error) {
	guide, err := s.ReadGuide(guideID)
	if err != nil {
		return nil, err
	}
	visitors, err := s.ListVisitors(guideID)
	if err != nil {
		return nil, err
	}
	blessings, err := s.ListBlessings(guideID)
	if err != nil {
		return nil, err
	}
	audits, err := s.ListAudits(guideID)
	if err != nil {
		return nil, err
	}
	payload := map[string]any{"guide": guide, "visitors": visitors, "blessings": blessings, "audits": audits}
	return json.MarshalIndent(payload, "", "  ")
}

func (s *Store) VerifyGuide(guideID string) error {
	guide, err := s.ReadGuide(guideID)
	if err != nil {
		return err
	}
	if err := guide.Validate(); err != nil {
		return err
	}
	if !domain.HasRequiredActions(guide.Links) {
		return domain.ErrInvalidGuide
	}
	if !domain.ValidateScheduleOrder(guide.Schedule) {
		return domain.ErrInvalidGuide
	}
	return nil
}

func (s *Store) DeleteGuide(guideID string) error {
	if strings.TrimSpace(guideID) == "" {
		return ErrNotFound
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range []string{guideBucket, venueBucket, attireBucket} {
			if err := tx.Bucket([]byte(bucket)).Delete([]byte(guideID)); err != nil {
				return err
			}
		}
		for _, bucket := range []string{scheduleBucket, linkBucket} {
			b := tx.Bucket([]byte(bucket))
			keys := make([][]byte, 0)
			if err := b.ForEach(func(key, _ []byte) error {
				if strings.HasPrefix(string(key), guideID+":") {
					keys = append(keys, append([]byte(nil), key...))
				}
				return nil
			}); err != nil {
				return err
			}
			for _, key := range keys {
				if err := b.Delete(key); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func forEachJSON(tx *bolt.Tx, bucket string, visit func([]byte) error) error {
	return tx.Bucket([]byte(bucket)).ForEach(func(_, value []byte) error {
		return visit(value)
	})
}
