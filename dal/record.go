package dal

import (
	"errors"
	"fmt"
)

// ErrDoesNotExist indicates a record does not exist
var ErrDoesNotExist = errors.New("does not exist")

// Record is a gateway to a database record.
type Record interface {
	// Key keeps a `table` name of an entity and an ID within that table or a chain of nested keys
	Key() *Key

	// Error keeps an error for the last operation on the record. Not found is not treated as an error
	Error() error

	// Exists indicates if record was found in database. Throws panic if called before a `Get` or `Set`.
	Exists() bool

	// SetError sets error relevant to specific record. Intended to be used only by DALgo DB drivers.
	SetError(err error)

	// Data returns record data (without ID/key).
	// Requires either record to be created by NewRecordWithData()
	// or DataTo() to be called first, otherwise panics.
	Data() any

	// HasChanged & MarkAsChanged are methods of convenience
	HasChanged() bool

	// MarkAsChanged & HasChanged are methods of convenience
	MarkAsChanged()

	//// SetDataTo sets DataTo handler
	//SetDataTo(dataTo func(target any) error)

	//// DataTo deserializes record data into a struct. Throws panic if called before `Get`.
	//// Uses a handler set by SetDataTo.
	//DataTo(target any) error
}

type record struct {
	key     *Key
	err     error
	changed bool
	data    any
	//dataTo func(target any) error
}

// Exists returns if records exists.
func (v *record) Exists() bool {
	if v.err != nil {
		if IsNotFound(v.err) {
			return false
		}
		if v.err == NoError {
			return true
		}
	}
	panic("an attempt to check if record exists before it was retrieved from database and SetError(error) called")
}

// Key returns key of a record
func (v *record) Key() *Key {
	return v.key
}

// HasChanged indicates if the record has changed since loading
func (v *record) HasChanged() bool {
	return v.changed
}

// MarkAsChanged marks the record as changed since loading
func (v *record) MarkAsChanged() {
	v.changed = true
}

func (v *record) Data() any {
	if v.err != nil {
		if v.err == NoError {
			return v.data
		}
		if !IsNotFound(v.err) {
			panic("an attempt to retrieve data from a record with an error")
		}
	}
	//panic("an attempt to access record data before it was retrieved from database and SetError(error) called")
	return v.data
}

//// SetDataTo sets DataTo handler
//func (v *record) SetDataTo(dataTo func(target any) error) {
//	v.dataTo = dataTo
//}
//
//// DataTo marshals record data into target
//func (v record) DataTo(target any) error {
//	if target == nil {
//		panic("not possible to marshall data into a nil value")
//	}
//	if v.dataTo == nil {
//		panic(fmt.Sprintf("method DataTo(%T) is called before data retrieval", target))
//	}
//	if err := v.dataTo(target); err != nil {
//		return fmt.Errorf("failed to marshal record data into %T: %w", target, err)
//	}
//	v.data = target
//	return nil
//}

// Error returns error associated with a record
func (v *record) Error() error {
	if v.err == NoError || IsNotFound(v.err) {
		return nil
	}
	return v.err
}

// SetError sets error associated with a record
func (v *record) SetError(err error) {
	if err == nil {
		v.err = NoError
	} else {
		v.err = err
	}
}

// NewRecord creates a new record
func NewRecord(key *Key) Record {
	return newRecord(key)
}

func newRecord(key *Key) *record {
	if key == nil {
		panic("parameter 'key' is required for dalgo.NewRecord()")
	}
	if err := key.Validate(); err != nil {
		panic(fmt.Errorf("invalid key: %w", err))
	}
	return &record{key: key}
}

// NewRecordWithData creates a new record with a data target struct
func NewRecordWithData(key *Key, data any) Record {
	record := newRecord(key)
	record.data = data
	return record
}

func NewRecordWithoutKey(data any) Record {
	return &record{data: data}
}
