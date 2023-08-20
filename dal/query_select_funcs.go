package dal

import (
	"errors"
	"math"
)

// SelectAllIDs is a helper method that for a given reader returns all IDs as a strongly typed slice.
func SelectAllIDs[T comparable](reader Reader, limit int) (ids []T, err error) {
	if reader == nil {
		panic("reader is a required parameter, got nil")
	}
	if limit >= 0 {
		ids = make([]T, 0, limit)
	} else {
		ids = make([]T, 0)
		limit = math.MaxInt
	}
	for ; limit > 0; limit-- {
		var record Record
		if record, err = reader.Next(); err != nil {
			if errors.Is(err, ErrNoMoreRecords) {
				err = nil
			}
			return
		}
		id := record.Key().ID.(T) // on separate line for debug purposes
		ids = append(ids, id)
	}
	return ids, reader.Close()
}

// SelectAllRecords	is a helper method that for a given reader returns all records as a slice.
func SelectAllRecords(reader Reader, limit int) (records []Record, err error) {
	if reader == nil {
		panic("reader is a required parameter, got nil")
	}
	if limit >= 0 {
		records = make([]Record, 0, limit)
	} else {
		records = make([]Record, 0)
		limit = math.MaxInt
	}
	for i := 0; limit <= 0 || i < limit; i++ {
		var record Record
		if record, err = reader.Next(); err != nil {
			if errors.Is(err, ErrNoMoreRecords) {
				err = nil
			}
			return
		}
		records = append(records, record)
	}
	return records, reader.Close()
}
