package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	bolt "go.etcd.io/bbolt"
)

const checkpointBucket = "ImportCheckpoint"

type Checkpoint struct {
	ID          string
	GuideID     string
	BatchID     string
	Processed   int
	Accepted    int
	Rejected    int
	Digest      string
	CompletedAt string
}

func (s *Store) SaveCheckpoint(checkpoint Checkpoint) error {
	if strings.TrimSpace(checkpoint.ID) == "" || strings.TrimSpace(checkpoint.GuideID) == "" {
		return ErrNotFound
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(checkpointBucket)); err != nil {
			return err
		}
		encoded, err := json.Marshal(checkpoint)
		if err != nil {
			return err
		}
		return tx.Bucket([]byte(checkpointBucket)).Put([]byte(checkpoint.ID), encoded)
	})
}

func (s *Store) ReadCheckpoint(id string) (Checkpoint, error) {
	var checkpoint Checkpoint
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(checkpointBucket))
		if bucket == nil {
			return ErrNotFound
		}
		value := bucket.Get([]byte(id))
		if value == nil {
			return ErrNotFound
		}
		return json.Unmarshal(value, &checkpoint)
	})
	return checkpoint, err
}

func (s *Store) ListCheckpoints(guideID string) ([]Checkpoint, error) {
	checkpoints := make([]Checkpoint, 0)
	err := s.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(checkpointBucket))
		if bucket == nil {
			return nil
		}
		return bucket.ForEach(func(_, value []byte) error {
			var checkpoint Checkpoint
			if err := json.Unmarshal(value, &checkpoint); err != nil {
				return err
			}
			if checkpoint.GuideID == guideID {
				checkpoints = append(checkpoints, checkpoint)
			}
			return nil
		})
	})
	sort.SliceStable(checkpoints, func(i, j int) bool { return checkpoints[i].CompletedAt < checkpoints[j].CompletedAt })
	return checkpoints, err
}

func DigestRows(values []string) string {
	ordered := append([]string(nil), values...)
	sort.Strings(ordered)
	hash := sha256.New()
	for _, value := range ordered {
		_, _ = hash.Write([]byte(value))
		_, _ = hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func CheckpointID(batchID, guideID string) string {
	if strings.TrimSpace(batchID) == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", guideID, batchID)
}
