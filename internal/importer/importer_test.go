package importer

import (
	"path/filepath"
	"testing"

	"weddingguide/internal/clock"
	"weddingguide/internal/ids"
	"weddingguide/internal/store"
)

func TestImportSmallBatch(t *testing.T) {
	repository, err := store.NewWithReaderLimit(filepath.Join(t.TempDir(), "small.db"), 4)
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	service := New(repository, clock.NewFixed("t1"), ids.New("import"))
	rows := []store.ImportRow{{VisitorKey: "a", Action: "view", SeenAt: "t1"}, {VisitorKey: "b", Action: "view", SeenAt: "t1"}, {VisitorKey: "c", Action: "view", SeenAt: "t1"}}
	report, err := service.ImportRows("g", "planner", rows)
	if err != nil || report.Processed != 3 || report.Created != 3 {
		t.Fatalf("small import failed: %#v %v", report, err)
	}
}

func TestBatchImportProcessesAllRecords(t *testing.T) {
	repository, err := store.NewWithReaderLimit(filepath.Join(t.TempDir(), "batch.db"), 4)
	if err != nil {
		t.Fatal(err)
	}
	defer repository.Close()
	service := New(repository, clock.NewFixed("t1"), ids.New("import"))
	rows := make([]store.ImportRow, 0, 8)
	for index := 0; index < 8; index++ {
		rows = append(rows, store.ImportRow{VisitorKey: string(rune('a' + index)), Action: "view", SeenAt: "t1"})
	}
	report, err := service.ImportRows("g", "planner", rows)
	if err != nil {
		t.Fatalf("batch import returned error: %v", err)
	}
	if report.Processed != len(rows) {
		t.Fatalf("processed %d of %d records", report.Processed, len(rows))
	}
}
