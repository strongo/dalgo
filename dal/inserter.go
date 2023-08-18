package dal

import (
	"context"
	"fmt"
	"github.com/strongo/random"
	"reflect"
)

// IDGenerator defines a contract for ID generator function
type IDGenerator = func(ctx context.Context, record Record) error

// InsertOptions defines interface for insert options
type InsertOptions interface {
	IDGenerator() IDGenerator
}

type insertOptions struct {
	idGenerator IDGenerator
}

func (v insertOptions) IDGenerator() IDGenerator {
	return v.idGenerator
}

var _ InsertOptions = (*insertOptions)(nil)

// NewInsertOptions creates insert options
func NewInsertOptions(opts ...InsertOption) InsertOptions {
	var options insertOptions
	for _, o := range opts {
		o(&options)
	}
	return options
}

// InsertOption defines a contract for an insert option
type InsertOption func(options *insertOptions)

type randomStringOptions struct {
	prefix string
}

// RandomStringOptions defines settings for random string
type RandomStringOptions interface {
	Prefix() string
}

// WithPrefix sets prefix for a random string
func WithPrefix(prefix string) func(options *randomStringOptions) {
	return func(options *randomStringOptions) {
		options.prefix = prefix
	}
}

// Prefix returns a predefined prefix for a random string
func (v randomStringOptions) Prefix() string {
	return v.prefix
}

type randomStringOption func(opts *randomStringOptions)

// WithIDGenerator sets ID generator for a random string (usually random)
func WithIDGenerator(c context.Context, g IDGenerator) KeyOption {
	return func(key *Key) error {
		if key.ID != nil {
			panic("an attempt to set ID generator for a child that already have an ID value")
		}
		return g(c, &record{key: key})
	}
}

// WithRandomStringID sets ID generator to random string
func WithRandomStringID(ctx context.Context, length int, options ...randomStringOption) KeyOption {
	var rso randomStringOptions
	for _, setOption := range options {
		setOption(&rso)
	}
	return func(key *Key) error {
		key.IDKind = reflect.String
		return WithIDGenerator(ctx, func(ctx context.Context, record Record) error {
			key.ID = rso.prefix + random.ID(length)
			return nil
		})(key)
	}
}

//// WithParent sets Parent
//func WithParent[T comparable](collection string, id T, options ...KeyOption) KeyOption {
//	return func(key *Key) (err error) {
//		options = append(options, WithID(id))
//		key.parent, err = NewKeyWithOptions(collection, options...)
//		return err
//	}
//}

//// WithParentKey sets Parent key
//func WithParentKey(parent *Key) KeyOption {
//	return func(key *Key) error {
//		key.parent = parent
//		return nil
//	}
//}

// WithStringID sets ID as a predefined string
func WithStringID(id string) KeyOption {
	return func(key *Key) error {
		key.ID = id
		return nil
	}
}

// InsertWithRandomID inserts a record with a random ID
func InsertWithRandomID(
	c context.Context,
	r Record,
	generateID IDGenerator,
	attempts int,
	exists func(*Key) error,
	insert func(Record) error,
) error {
	key := r.Key()
	// We need a temp record to make sure we do not overwrite data during exists() check
	tmp := &record{key: key}
	for i := 1; i <= attempts; i++ {
		if err := generateID(c, tmp); err != nil {
			return fmt.Errorf("failed to generate random value: %w", err)
		}
		if err := exists(key); err == nil {
			continue
		} else if IsNotFound(err) {
			return insert(r) // r shares child with tmp
		} else {
			return fmt.Errorf("failed to check if record exists: %w", err)
		}
	}
	return fmt.Errorf("not able to generate unique Value in %v attempts", attempts)
}
