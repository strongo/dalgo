package dal

import (
	"bytes"
	"context"
	"fmt"
	"github.com/strongo/dalgo/query"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// CollectionRef points to a collection (e.g. table) in a database
type CollectionRef struct {
	Name   string
	Parent *Key
}

func (v CollectionRef) Path() string {
	if v.Parent == nil {
		return v.Name
	}
	return v.Parent.String() + "/" + v.Name
}

// Select holds definition of a query
type Select struct {

	// From defines target table/collection
	From *CollectionRef

	// Where defines filter condition
	Where query.Condition

	// GroupBy defines expressions to group by
	GroupBy []query.Expression

	// OrderBy defines expressions to order by
	OrderBy []query.Expression

	// Columns defines what columns to return
	Columns []query.Column

	Into func() interface{}

	// Limit specifies maximum number of records to be returned
	Limit int
}

func (q Select) String() string {
	writer := bytes.NewBuffer(make([]byte, 0, 1024))
	writer.WriteString("SELECT")
	if q.Limit > 0 {
		writer.WriteString(" TOP " + strconv.Itoa(q.Limit))
	}
	switch len(q.Columns) {
	case 0:
		writer.WriteString(" *")
	case 1:
		_, _ = fmt.Fprint(writer, " ", q.Columns[0].String())
	default:
		for _, col := range q.Columns {
			_, _ = fmt.Fprint(writer, "\n\t", col.String())
		}
	}
	is1liner := len(q.Columns) <= 1 &&
		(q.Where == nil || reflect.TypeOf(q.Where) == reflect.TypeOf(query.Comparison{}))

	if q.From != nil {
		if is1liner {
			writer.WriteString(" ")
		} else {
			writer.WriteString("\n")
		}
		fmt.Fprintf(writer, "FROM [%v]", q.From.Path())
	}
	if q.Where != nil {
		if is1liner {
			writer.WriteString(" ")
		} else {
			writer.WriteString("\n")
		}
		writer.WriteString("WHERE " + q.Where.String())
	}
	if len(q.GroupBy) > 0 {
		writer.WriteString("\nGROUP BY ")
		for _, expr := range q.GroupBy {
			writer.WriteString("\n\t")
			writer.WriteString(expr.String())
		}
	}
	return writer.String()
}

var _ fmt.Stringer = (*Select)(nil)

// And creates a new query by adding a condition to a predefined query
func (q Select) groupWithConditions(operator query.Operator, conditions ...query.Condition) Select {
	qry := Select{From: q.From}
	and := groupCondition{operator: operator, Conditions: make([]query.Condition, len(conditions)+1)}
	and.Conditions[0] = q.Where
	for i, condition := range conditions {
		and.Conditions[i+1] = condition
	}
	qry.Where = and
	return qry
}

// And creates an inherited query by adding AND conditions
func (q Select) And(conditions ...query.Condition) Select {
	return q.groupWithConditions(query.And, conditions...)
}

// Or creates an inherited query by adding OR conditions
func (q Select) Or(conditions ...query.Condition) Select {
	return q.groupWithConditions(query.Or, conditions...)
}

type groupCondition struct {
	operator   query.Operator
	Conditions []query.Condition
}

func (v groupCondition) Operator() query.Operator {
	return v.operator
}

func (v groupCondition) String() string {
	s := make([]string, len(v.Conditions))
	for i, condition := range v.Conditions {
		s[i] = condition.String()
	}
	return fmt.Sprintf("(%v)", strings.Join(s, string(v.operator)))
}

// ReadAll reads all records from a reader
func ReadAll(_ context.Context, reader Reader, limit int) (records []Record, err error) {
	var record Record
	if limit <= 0 {
		limit = math.MaxInt64
	}
	for i := 0; i < limit; i++ {
		if i >= limit {
			break
		}
		if record, err = reader.Next(); err != nil {
			if err == ErrNoMoreRecords {
				break
			}
			records = append(records, record)
		}
	}
	return records, err
}
