package dalmock

import (
	"context"
	"github.com/strongo/dalgo/dal"
)

type readonlySession struct {
	onSelectFrom map[string]SelectResult
}

var _ dal.ReadSession = (*readonlySession)(nil)

func (d readonlySession) Get(ctx context.Context, record dal.Record) error {
	//TODO implement me
	panic("implement me")
}

func (d readonlySession) GetMulti(ctx context.Context, records []dal.Record) error {
	//TODO implement me
	panic("implement me")
}

func (d readonlySession) Select(ctx context.Context, query dal.Select) (dal.Reader, error) {
	collectionPath := query.From.Path()
	result := d.onSelectFrom[collectionPath]
	return result.reader(query.Into), result.err
}
