package importer

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"weddingguide/internal/store"
)

var ErrInvalidCSV = errors.New("invalid visitor csv")

type RowParser struct {
	separator rune
}

func NewRowParser(separator rune) RowParser {
	if separator == 0 {
		separator = ','
	}
	return RowParser{separator: separator}
}

func (p RowParser) Parse(reader io.Reader) ([]store.ImportRow, error) {
	csvReader := csv.NewReader(reader)
	csvReader.Comma = p.separator
	csvReader.FieldsPerRecord = -1
	rows := make([]store.ImportRow, 0)
	line := 0
	for {
		record, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		line++
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		if line == 1 && len(record) > 0 && strings.EqualFold(strings.TrimSpace(record[0]), "visitor") {
			continue
		}
		row, err := p.record(record)
		if err != nil {
			return nil, fmt.Errorf("line %d: %w", line, err)
		}
		rows = append(rows, row)
	}
	return NormalizeRows(rows), nil
}

func (p RowParser) record(record []string) (store.ImportRow, error) {
	if len(record) < 2 {
		return store.ImportRow{}, ErrInvalidCSV
	}
	row := store.ImportRow{VisitorKey: strings.TrimSpace(record[0]), Action: strings.TrimSpace(record[1])}
	if len(record) > 2 {
		row.SeenAt = strings.TrimSpace(record[2])
	}
	if row.VisitorKey == "" || row.Action == "" {
		return store.ImportRow{}, ErrInvalidCSV
	}
	return row, nil
}

func NormalizeRows(rows []store.ImportRow) []store.ImportRow {
	normalized := make([]store.ImportRow, 0, len(rows))
	seen := make(map[string]bool)
	for _, row := range rows {
		row.VisitorKey = strings.ToLower(strings.TrimSpace(row.VisitorKey))
		row.Action = strings.ToLower(strings.TrimSpace(row.Action))
		row.SeenAt = strings.TrimSpace(row.SeenAt)
		if row.SeenAt == "" {
			row.SeenAt = "fixture-time"
		}
		if row.VisitorKey == "" || row.Action == "" || seen[row.VisitorKey] {
			continue
		}
		seen[row.VisitorKey] = true
		normalized = append(normalized, row)
	}
	return normalized
}

func ValidateRows(rows []store.ImportRow) error {
	if len(rows) == 0 {
		return ErrInvalidCSV
	}
	for index, row := range rows {
		if row.VisitorKey == "" {
			return fmt.Errorf("row %d missing visitor", index)
		}
		if row.Action == "" {
			return fmt.Errorf("row %d missing action", index)
		}
	}
	return nil
}

func SortRows(rows []store.ImportRow) []store.ImportRow {
	ordered := append([]store.ImportRow(nil), rows...)
	sort.SliceStable(ordered, func(i, j int) bool {
		if ordered[i].VisitorKey == ordered[j].VisitorKey {
			return ordered[i].Action < ordered[j].Action
		}
		return ordered[i].VisitorKey < ordered[j].VisitorKey
	})
	return ordered
}

func EncodeRows(rows []store.ImportRow) string {
	ordered := SortRows(rows)
	var builder strings.Builder
	builder.WriteString("visitor,action,seen_at\n")
	for _, row := range ordered {
		builder.WriteString(row.VisitorKey)
		builder.WriteByte(',')
		builder.WriteString(row.Action)
		builder.WriteByte(',')
		builder.WriteString(row.SeenAt)
		builder.WriteByte('\n')
	}
	return builder.String()
}
